# secrets-operator

This project is a Kubernetes operator for managing secrets in your cluster.  This relies on the GitRepository CRD and [Source Controller](https://fluxcd.io/docs/components/source/) from the Flux project.  The secrets are stored gpg encrypted in a git repository, such as those managed by [Pass](https://www.passwordstore.org/) or [Gopass](https://github.com/gopasspw/gopass).

A more general alternative is [external-secrets.io](https://external-secrets.io/)

# Installation

TBD

# Usage

Add a Flux GitRepository:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: GitRepository
metadata:
  name: credentials
  namespace: secrets-operator
spec:
  interval: 360m
  ref:
    branch: master
  secretRef:
    name: repo-secret
  timeout: 20s
  url: ssh://git@github.com/s1devops/<repo>.git

```

Create a secret which holds an ascii armoured gpg private key (must not be passphrase protected) for the SecretSource:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: credentials-private-key
  namespace: secrets-operator
stringData:
  privateKey: |
    -----BEGIN PGP PRIVATE KEY BLOCK-----

    .....
    -----END PGP PRIVATE KEY BLOCK-----

```


Create a SecretSource linking the GitRepository and the gpg key together:

```yaml
apiVersion: secrets.s1devops.com/v1alpha1
kind: SecretSource
metadata:
  name: credentials
  namespace: secrets-operator
spec:
  gitRepository:
    namespace: secrets-operator
    name: credentials
  privateKey:
    name: credentials-private-key
    key: privateKey
  
```

Create a SecretMapping for every Secret you wish to create.  Example:

```yaml
apiVersion: secrets.s1devops.com/v1alpha1
kind: SecretMapping
metadata:
  namespace: some-target-namespace
  name: some-secret
spec:
  source:
    namespace: secret-operator
    name: credentials
  name: name-of-secret-to-create # defaults to SecretMapping name
  mappings:
    - name: secretKey
      type: string
      value: this would be stored in the secret verbatim
    - name: anotherSecretKey
      type: pass
      value: path/to/secret/in/gitRepo
    - name: yetAnotherSecretKey
      type: template
      value: Uses golang template, can look up secrets like {{ pass "path/to/secret/in/gitRepo" }}
```
