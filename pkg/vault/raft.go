package vault

import (
	"bytes"
	"context"

	vault2 "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

type Vault struct {
	client *vault2.Client
}

func NewVault(cl *vault2.Client) *Vault {
	return &Vault{
		client: cl,
	}
}

func (v *Vault) Backup(ctx context.Context) ([]byte, error) {
	var buf bytes.Buffer

	if err := v.client.Sys().RaftSnapshotWithContext(ctx, &buf); err != nil {
		return nil, errors.WithStack(err)
	}

	return buf.Bytes(), nil
}
