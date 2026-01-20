package kubeconfigSourceModel

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Expand(m KubeconfigSourceModel) (id string, bytes []byte, context string, err error) {
	var sources int

	if len(m.Kubeconfig.File) > 0 {
		sources++
		f := m.Kubeconfig.File[0]

		path := f.Path.ValueString()
		if path == "" {
			return "", nil, "", fmt.Errorf("kubeconfig.file.path is required")
		}

		path, err = expandPath(path)
		if err != nil {
			return "", nil, "", err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return "", nil, "", err
		}

		bytes = data
		context = f.Context.ValueString()
		id = hash("file", path, context)
	}

	if len(m.Kubeconfig.Env) > 0 {
		sources++
		e := m.Kubeconfig.Env[0]

		name := e.Name.ValueString()
		if name == "" {
			return "", nil, "", fmt.Errorf("kubeconfig.env.name is required")
		}

		val := os.Getenv(name)
		if val == "" {
			return "", nil, "", fmt.Errorf("env %q is empty", name)
		}

		switch strings.ToLower(e.Encoding.ValueString()) {
		case "", "plain":
			bytes = []byte(val)

		case "base64":
			decoded, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return "", nil, "", fmt.Errorf("failed to decode base64 env %q: %w", name, err)
			}
			bytes = decoded

		default:
			return "", nil, "", fmt.Errorf("unsupported kubeconfig.env.encoding %q", e.Encoding.ValueString())
		}

		context = e.Context.ValueString()
		id = hash("env", name, context)
	}

	if len(m.Kubeconfig.S3) > 0 {
		sources++
		return "", nil, "", fmt.Errorf("kubeconfig.s3 not implemented")
	}

	if sources == 0 {
		return "", nil, "", fmt.Errorf("one kubeconfig source must be specified")
	}

	if sources > 1 {
		return "", nil, "", fmt.Errorf("only one kubeconfig source is allowed")
	}

	return id, bytes, context, nil
}

func hash(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func expandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, strings.TrimPrefix(p, "~")), nil
	}
	return p, nil
}
