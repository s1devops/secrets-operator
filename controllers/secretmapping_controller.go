/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	secretsv1alpha1 "github.com/s1devops/secrets-operator/api/v1alpha1"
	"github.com/s1devops/secrets-operator/secrets"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	FILE_MODE = os.FileMode(int(0700))
)

// SecretMappingReconciler reconciles a SecretMapping object
type SecretMappingReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	GitCacheDir string
}

//+kubebuilder:rbac:groups=secrets.s1devops.com,resources=secretmappings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=secrets.s1devops.com,resources=secretmappings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=secrets.s1devops.com,resources=secretmappings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SecretMapping object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SecretMappingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// grab SecretMapping
	secretMapping := &secretsv1alpha1.SecretMapping{}
	err := r.Get(ctx, req.NamespacedName, secretMapping)
	if err != nil {

		if errors.IsNotFound(err) {
			// secretMapping must have been deleted, no further action
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	secretSource := &secretsv1alpha1.SecretSource{}
	err = r.Get(ctx, types.NamespacedName(secretMapping.Spec.Source), secretSource)
	if err != nil {
		return ctrl.Result{}, err
	}

	privateKeySecret := &apiv1.Secret{}
	err = r.Get(
		ctx,
		types.NamespacedName{Namespace: secretSource.Namespace, Name: secretSource.Spec.PrivateKey.Name},
		privateKeySecret,
	)
	if err != nil {
		logger.Error(err, "Failed to find private key secret", "SecretMapping.Namespace", secretMapping.Namespace, "SecretMapping.Name", secretMapping.Name)
		return ctrl.Result{}, err
	}

	privateKey, ok := privateKeySecret.Data[secretSource.Spec.PrivateKey.Key]
	if !ok {
		err = fmt.Errorf("missing secret key field %s", secretSource.Spec.PrivateKey.Key)
		logger.Error(err, "Failed to find private key secret field", "SecretMapping.Namespace", secretMapping.Namespace, "SecretMapping.Name", secretMapping.Name)
		return ctrl.Result{}, err
	}

	privateKeyHasher := sha256.New()
	privateKeyHasher.Write(privateKey)
	privateKeyHash := hex.EncodeToString(privateKeyHasher.Sum(nil))

	privateKeyPath := filepath.Join(r.GitCacheDir, privateKeyHash) + ".asc"
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		// Never imported
		if err = ioutil.WriteFile(privateKeyPath, privateKey, FILE_MODE); err != nil {
			logger.Error(err, "Failed to find save private key", "SecretMapping.Namespace", secretMapping.Namespace, "SecretMapping.Name", secretMapping.Name)
			return ctrl.Result{}, err
		}

		if err = exec.Command("gpg", "--import", privateKeyPath).Run(); err != nil {
			logger.Error(err, "Failed to import private key", "SecretMapping.Namespace", secretMapping.Namespace, "SecretMapping.Name", secretMapping.Name)
			return ctrl.Result{}, err
		}
	}

	gitRepo := &sourcev1.GitRepository{}
	err = r.Get(ctx, types.NamespacedName(secretSource.Spec.GitRepository), gitRepo)
	if err != nil {
		return ctrl.Result{}, err
	}

	revision := determineRevision(gitRepo)
	revisionCacheDir := filepath.Join(r.GitCacheDir, revision)
	if _, err := os.Stat(revisionCacheDir); os.IsNotExist(err) {
		logger.Info("Downloading git revision", "SecretMapping.Namespace", secretMapping.Namespace, "SecretMapping.Name", secretMapping.Name, "Revision", revision)

		revisionUrl := gitRepo.Status.Artifact.URL
		if err = downloadArtifact(revisionUrl, revisionCacheDir); err != nil {
			logger.Error(err, "Failed downloading git revision", "SecretMapping.Namespace", secretMapping.Namespace, "SecretMapping.Name", secretMapping.Name, "Revision", revision)
			return ctrl.Result{}, err
		}
	}

	handler := secrets.NewHandler(revisionCacheDir)

	secretName := secretMapping.Spec.Name
	if secretName == "" {
		secretName = secretMapping.Name
	}

	found := &apiv1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Namespace: secretMapping.Namespace, Name: secretName}, found)
	if err != nil {

		if errors.IsNotFound(err) {
			// Create

			newSecret := &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: secretMapping.Namespace,
					Name:      secretName,
				},
			}

			ctrl.SetControllerReference(secretMapping, newSecret, r.Scheme)
			err = r.Create(ctx, newSecret)
			if err != nil {
				logger.Error(err, "Failed to create new Secret", "Secret.Namespace", newSecret.Namespace, "Secret.Name", newSecret.Name)

				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, err
	}

	data := make(map[string][]byte)

	for _, mapping := range secretMapping.Spec.Mappings {
		if mapping.Type == secretsv1alpha1.TYPE_STRING {
			data[mapping.Name] = []byte(mapping.Value)
		} else if mapping.Type == secretsv1alpha1.TYPE_PASS {

			value, err := handler.Get(mapping.Value)
			if err != nil {
				logger.Error(err, "Failed to get secret", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name, "Secret.Key", mapping.Value)
			}
			data[mapping.Name] = value
		} else if mapping.Type == secretsv1alpha1.TYPE_TEMPLATE {
			value, err := handler.Template(mapping.Value)
			if err != nil {
				logger.Error(err, "Failed to get secret", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name, "Secret.Key", mapping.Value)
			}
			data[mapping.Name] = value
		}
	}

	if reflect.DeepEqual(found.Data, data) {
		return ctrl.Result{}, nil
	}

	found.Data = data

	err = r.Update(ctx, found)
	if err != nil {
		logger.Error(err, "Failed to update Secret", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func determineRevision(gitRepo *sourcev1.GitRepository) string {
	rawRevision := gitRepo.Status.Artifact.Revision
	portions := strings.Split(rawRevision, "/")
	if len(portions) != 2 {
		return ""
	}

	return portions[1]
}

func downloadArtifact(href string, cacheDir string) error {

	req, err := http.NewRequest(http.MethodGet, href, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	gr, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if(header.Typeflag != tar.TypeReg) {
			// Ignore any non regular files
			continue;
		}
		destinationPath := filepath.Join(cacheDir, header.Name)

		destinationDir := filepath.Dir(destinationPath)
		if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
			err = os.MkdirAll(destinationDir, FILE_MODE)
			if err != nil {
				return err
			}
		}

		w, err := os.Create(destinationPath)
		if err != nil {
			return err
		}
		defer w.Close()

		if _, err = io.Copy(w, tr); err != nil {
			return err
		}

	}

	return nil
}

func secretSourceToSecretMappings(mgr ctrl.Manager, secretSource *secretsv1alpha1.SecretSource) []reconcile.Request {
	// logger := log.Log.WithName("controller")

	secretSourceNamespace := secretSource.Namespace

	secretSourceName := secretSource.Name

	cli := mgr.GetClient()

	secretMappingList := &secretsv1alpha1.SecretMappingList{}
	err := cli.List(context.TODO(), secretMappingList)
	if err != nil {
		return []reconcile.Request{}
	}

	var reconcileRequests []reconcile.Request

	for _, secretMapping := range secretMappingList.Items {
		mappingSource := secretMapping.Spec.Source
		if mappingSource.Namespace == secretSourceNamespace && mappingSource.Name == secretSourceName {
			rec := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      secretMapping.Name,
					Namespace: secretMapping.Namespace,
				},
			}
			reconcileRequests = append(reconcileRequests, rec)
		}
	}

	return reconcileRequests
}

func gitRepositoryToSecretMappings(mgr ctrl.Manager, gitRepo *sourcev1.GitRepository) []reconcile.Request {
	// logger := log.Log.WithName("controller")

	namespace := gitRepo.Namespace

	name := gitRepo.Name

	cli := mgr.GetClient()

	secretSourceList := &secretsv1alpha1.SecretSourceList{}
	err := cli.List(context.TODO(), secretSourceList)
	if err != nil {
		return []reconcile.Request{}
	}

	var reconcileRequests []reconcile.Request

	for _, secretSource := range secretSourceList.Items {
		sourceRepo := secretSource.Spec.GitRepository
		if sourceRepo.Namespace == namespace && sourceRepo.Name == name {
			rec := secretSourceToSecretMappings(mgr, &secretSource)

			reconcileRequests = append(reconcileRequests, rec...)
		}
	}

	return reconcileRequests
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretMappingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretsv1alpha1.SecretMapping{}).
		Watches(
			&source.Kind{Type: &secretsv1alpha1.SecretSource{}},
			handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
				secretSource := obj.(*secretsv1alpha1.SecretSource)
				return secretSourceToSecretMappings(mgr, secretSource)
			}),
		).
		Watches(
			&source.Kind{Type: &sourcev1.GitRepository{}},
			handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
				gitRepo := obj.(*sourcev1.GitRepository)
				return gitRepositoryToSecretMappings(mgr, gitRepo)
			}),
		).
		Owns(&apiv1.Secret{}).
		Complete(r)
}
