package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

func TestS3Success(t *testing.T) {
	t.Setenv("VAULT_CLUSTER_URL", "http://localhost:8200")
	//t.Setenv("VAULT_CLUSTER_ROOT_TOKEN", "hvs.AVdPZp6PoY13KSbxR47YI5L9")
	logger := log.Logger

	vaultUrl := os.Getenv("VAULT_CLUSTER_URL")
	if err := waitCluster(vaultUrl, logger); err != nil {
		logger.Panic().Err(err).Send()
	}

	initResp, err := initCluster(vaultUrl)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	t.Setenv("VAULT_NAME", "e2e")
	t.Setenv("VAULT_URL", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", initResp.RootToken)

	t.Setenv("S3_ACCESS_KEY", "ROOTUSER")
	t.Setenv("S3_ENDPOINT", "CHANGEME123")
	t.Setenv("S3_REGION", "http://127.0.0.1:9000")
	t.Setenv("S3_DISABLE_SSL", "true")
	t.Setenv("S3_BUCKET", "backup")

	main()

	fmt.Println(initResp)
}

type initResponse struct {
	Keys       []string `json:"keys"`
	KeysBase64 []string `json:"keys_base64"`
	RootToken  string   `json:"root_token"`
}

func waitCluster(vaultUrl string, logger zerolog.Logger) error {
	for i := 0; i < 30; i++ {
		resp, err := req.Get(fmt.Sprintf("%v/v1/sys/health", vaultUrl))
		if err != nil {
			logger.Err(errors.Wrap(err, "waiting for healthy cluster"))
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.IsErrorState() {
			logger.Err(errors.Wrapf(err, "health status code %v, expected 200", resp.Status))
			time.Sleep(1 * time.Second)
			continue
		}

		return nil
	}

	return errors.New("can not receive healthy cluster after 30 retries")
}

func initCluster(vaultUrl string) (*initResponse, error) {
	req.DevMode()
	getBytes := req.MustGet(fmt.Sprintf("%v/v1/sys/seal-status", vaultUrl)).Bytes()
	var initResp initResponse

	if !gjson.GetBytes(getBytes, "initialized").Bool() {
		_ = req.NewRequest().SetBodyJsonString(`{"secret_shares":1,"secret_threshold":1}`).
			SetSuccessResult(&initResp).
			MustPut(fmt.Sprintf("%v/v1/sys/seal-status", vaultUrl))
	} else {
		// if its local tests, we already should have that data in envs
		initResp = initResponse{
			Keys:       []string{os.Getenv("VAULT_CLUSTER_KEY")},
			KeysBase64: nil,
			RootToken:  os.Getenv("VAULT_CLUSTER_ROOT_TOKEN"),
		}
	}

	if gjson.GetBytes(
		req.MustGet(fmt.Sprintf("%v/v1/sys/seal-status", vaultUrl)).Bytes(),
		"sealed",
	).Bool() { // lets unseal it
		req.NewRequest().SetBodyJsonMarshal(map[string]interface{}{
			"key": initResp.Keys[0],
		}).MustPut(fmt.Sprintf("%v/v1/sys/unseal", vaultUrl))
	}

	return &initResp, nil
}
