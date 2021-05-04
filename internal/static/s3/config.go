package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/static"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SSL             bool
	Location        string
}

func (c *Config) NewStorage() (static.Storage, error) {
	minioClient, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure: c.SSL,
		Region: c.Location,
	})
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "MINIO-4m90d", "Errors.Assets.Store.NotInitialized")
	}
	return &Minio{
		Client:   minioClient,
		Location: c.Location,
	}, nil
}
