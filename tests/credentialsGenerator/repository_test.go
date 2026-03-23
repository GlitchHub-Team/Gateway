package credentialsgeneratortests

import (
	"testing"

	credentialsgenerator "Gateway/internal/credentialsGenerator"

	"github.com/nats-io/nkeys"
)

func TestGenerateCredentialsReturnsConsistentKeyPair(t *testing.T) {
	creds, err := credentialsgenerator.GenerateCredentials()
	if err != nil {
		t.Fatalf("expected credentials generation to succeed, got %v", err)
	}
	if creds == nil {
		t.Fatal("expected non nil credentials")
	}
	if creds.PublicIdentifier == "" {
		t.Fatal("expected public identifier to be populated")
	}
	if creds.SecretKey == "" {
		t.Fatal("expected secret key to be populated")
	}

	kp, err := nkeys.FromSeed([]byte(creds.SecretKey))
	if err != nil {
		t.Fatalf("expected valid seed, got %v", err)
	}

	publicKey, err := kp.PublicKey()
	if err != nil {
		t.Fatalf("expected public key derivation to succeed, got %v", err)
	}
	if publicKey != creds.PublicIdentifier {
		t.Fatalf("expected public key %s, got %s", publicKey, creds.PublicIdentifier)
	}
}
