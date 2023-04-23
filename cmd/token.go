package main

import (
	"os"
	"strings"

	vault2 "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func getVaultToken(client *vault2.Client, logger zerolog.Logger) (string, error) {
	switch strings.ToLower(os.Getenv("AUTH_TYPE")) {
	case "kubernetes":
		kubernetesRole := os.Getenv("KUBERNETES_ROLE")
		logger.Info().Msgf("using kubernetes auth with role: %v", kubernetesRole)
		serviceAccountTokenBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			return "", errors.WithStack(err)
		}
		options := map[string]interface{}{
			"jwt":  string(serviceAccountTokenBytes),
			"role": kubernetesRole,
		}
		secret, err := client.Logical().Write("auth/kubernetes/login", options)
		if err != nil {
			return "", errors.WithStack(err)
		}
		logger.Info().Msg("kubernetes auth success")
		return secret.Auth.ClientToken, nil
	default:
		logger.Info().Msg("reading VAULT_TOKEN env")
		return os.Getenv("VAULT_TOKEN"), nil
	}
}
