package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func checkVault(vaultPath string, logger zerolog.Logger) error {
	if vaultPath == "" {
		return nil
	}

	logger.Info().Msgf("Checking vault files on %v", vaultPath)

	files, err := os.ReadDir(vaultPath)
	if err != nil {
		logger.Error().Err(err).Msg("can not read vault folder")
		return errors.WithStack(err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	tempFile := map[string]interface{}{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		targetFile := path.Join(vaultPath, f.Name())
		logger.Info().Msgf("[VAULT] Reading file [%v]", targetFile)

		data, fileErr := os.ReadFile(targetFile)
		if fileErr != nil {
			return errors.Wrapf(err, "can not read file [%v]", targetFile)
		}

		if err = json.Unmarshal(data, &tempFile); err != nil {
			return errors.Wrapf(err, "can not parse file [%v]", targetFile)
		}

		for k, v := range tempFile {
			logger.Info().Msgf("[VAULT] Applying key [%v] from [%v]", k, targetFile)

			if err = os.Setenv(k, fmt.Sprint(v)); err != nil {
				return errors.Wrapf(err, "can not set env [%v] from [%v]", k, targetFile)
			}
		}
	}

	return nil
}
