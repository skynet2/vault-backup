package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	time2 "time"

	"github.com/hashicorp/go-multierror"
	vault2 "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"

	"github.com/skynet2/vault-backup/pkg/desitnation"
	"github.com/skynet2/vault-backup/pkg/vault"
)

var (
	successCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "vault_backup_success_total",
		Help: "The total number of success backups",
	}, []string{"vault_name"})

	failCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "vault_backup_errors_total",
		Help: "The total number of errors during backups",
	}, []string{"vault_name", "err_text"})

	isRegistered = false
)

func registerMetrics() {
	if isRegistered {
		return
	}

	prometheus.MustRegister(successCounter)
	prometheus.MustRegister(failCounter)
	isRegistered = true
}

func main() {
	// trigger build
	logger := log.Logger
	registerMetrics()

	if vaultSecretsErr := checkVault(os.Getenv("VAULT_SECRETS_PATH"), logger); vaultSecretsErr != nil {
		panic(vaultSecretsErr)
	}

	name := os.Getenv("VAULT_NAME")
	if name == "" {
		name = "default"
	}

	defer func() {
		if pushErr := pushMetrics(os.Getenv("PROMETHEUS_PUSH_GATEWAY_URL")); pushErr != nil {
			logger.Err(pushErr).Send()
		}
	}()

	s3Destination, err := desitnation.NewS3()
	if err != nil {
		failCounter.WithLabelValues(name, err.Error()).Inc()
		panic(err)
	}

	var client *vault2.Client
	var finalErr error
	for _, url := range strings.Split(os.Getenv("VAULT_URLS"), ",") {
		logger.Info().Msgf("connecting to %v", url)
		vaultClient, clientErr := vault2.NewClient(&vault2.Config{
			Address: url,
		})

		if clientErr != nil {
			clientErr = errors.Wrapf(clientErr, "can not connect to %v", url)
			finalErr = multierror.Append(finalErr, clientErr)
			continue
		}

		_, clientErr = vaultClient.Sys().Health()
		if clientErr != nil {
			clientErr = errors.Wrapf(clientErr, "can not get health of %v", url)
			finalErr = multierror.Append(finalErr, clientErr)
			continue
		}

		client = vaultClient
		logger.Info().Msgf("connected to %v", url)
		break
	}

	if client == nil {
		panic(finalErr)
	}

	ctx := logger.WithContext(context.Background())
	if err != nil {
		failCounter.WithLabelValues(name, err.Error()).Inc()
		panic(err)
	}

	token, err := getVaultToken(client, logger)
	if err != nil {
		failCounter.WithLabelValues(name, err.Error()).Inc()
		panic(err)
	}

	client.SetToken(token)

	backupData, err := vault.NewVault(client).Backup(ctx)
	if err != nil {
		failCounter.WithLabelValues(name, err.Error()).Inc()
		panic(err)
	}

	tt := time2.Now().UTC()
	finalPath := path.Join(
		fmt.Sprint(tt.Year()),
		fmt.Sprintf("%d", tt.Month()),
		fmt.Sprint(tt.Day()),
		fmt.Sprintf("%s.data", tt.Format("2006-01-02_15_04_05")),
	)

	if err = s3Destination.Upload(ctx, finalPath, backupData); err != nil {
		failCounter.WithLabelValues(name, err.Error()).Inc()
		panic(err)
	}

	successCounter.WithLabelValues(name).Inc()
}
