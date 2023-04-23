package main

import (
	"context"
	"os"

	vault2 "github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"

	"github.com/skynet2/vault-backup/pkg/vault"
)

var (
	successCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "vault_backup_success_total",
		Help: "The total number of success backups",
	})

	failCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "vault_backup_errors_total",
		Help: "The total number of errors during backups",
	})
)

func main() {
	prometheus.MustRegister(successCounter)
	prometheus.MustRegister(failCounter)
	defer func() {
		_ = pushMetrics(os.Getenv("PROMETHEUS_PUSH_GATEWAY_URL"))
	}()
	logger := log.Logger

	client, err := vault2.NewClient(&vault2.Config{
		Address: os.Getenv("VAULT_URL"),
	})

	ctx := logger.WithContext(context.Background())
	if err != nil {
		logger.Panic().Err(err).Send()
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	backupData, err := vault.NewVault(client).Backup(ctx)
	if err != nil {
		logger.Panic().Err(err).Send()
	}

}
