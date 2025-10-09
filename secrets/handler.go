package secrets

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
)

type Handler struct {
	cacheDir string
}

func NewHandler(cacheDir string) *Handler {
	return &Handler{
		cacheDir: cacheDir,
	}
}

func (h *Handler) Get(key string) ([]byte, error) {
	path := filepath.Join(h.cacheDir, key) + ".gpg"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("could not find key %s", key)
	}

	cmd := exec.Command("gpg", "--decrypt", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (h *Handler) Template(tmpl string) ([]byte, error) {
	t := template.New("tpl")
	if t == nil {
		return nil, errors.New("Could not build template parser")
	}

	funcs := template.FuncMap{
		"pass": func(k string) (string, error) {
			val, err := h.Get(k)
			if err != nil {
				return "", err
			}

			return string(val), nil
		},
	}

	t = t.Funcs(funcs)

	parsed, err := t.Parse(tmpl)
	if err != nil {
		return nil, err
	}

	var data interface{}
	out := new(bytes.Buffer)
	err = parsed.Execute(out, data)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
