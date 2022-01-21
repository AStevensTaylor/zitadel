package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"time"
)

type KeyPair struct {
	es_models.ObjectRoot

	Usage      KeyUsage
	Algorithm  string
	PrivateKey *Key
	PublicKey  *Key
}

type KeyUsage int32

const (
	KeyUsageSigning KeyUsage = iota
	KeyUsageSAMLMetadataSigning
	KeyUsageSAMLResponseSinging
	KeyUsageSAMLCA
)

func (u KeyUsage) String() string {
	switch u {
	case KeyUsageSigning:
		return "sig"
	case KeyUsageSAMLCA:
		return "saml_ca"
	case KeyUsageSAMLResponseSinging:
		return "saml_response_sig"
	case KeyUsageSAMLMetadataSigning:
		return "saml_metadata_sig"
	}
	return ""
}

type Key struct {
	Key    *crypto.CryptoValue
	Expiry time.Time
}

func (k *KeyPair) IsValid() bool {
	return k.Algorithm != "" &&
		k.PrivateKey != nil && k.PrivateKey.IsValid() &&
		k.PublicKey != nil && k.PublicKey.IsValid()
}

func (k *Key) IsValid() bool {
	return k.Key != nil
}
