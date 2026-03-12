package credentialsgenerator

import (
	"github.com/nats-io/nkeys"
)

type BaseJWT string

type NKeyCredentialsGenerator struct{}

func (g *NKeyCredentialsGenerator) GenerateCredentials() (*Credentials, error) {
	kp, err := nkeys.CreateUser()
	if err != nil {
		return nil, err
	}

	publicKey, err := kp.PublicKey()
	if err != nil {
		return nil, err
	}

	seed, err := kp.Seed()
	if err != nil {
		return nil, err
	}

	return &Credentials{
		PublicIdentifier: publicKey,
		SecretKey:        string(seed),
	}, nil
}

func NewNKeyCredentialsGenerator() *NKeyCredentialsGenerator {
	return &NKeyCredentialsGenerator{}
}
