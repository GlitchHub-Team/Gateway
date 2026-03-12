package credentialsgenerator

type Credentials struct {
	PublicIdentifier string
	SecretKey        string
	Token            string
}

type CredentialsGeneratorPort interface {
	GenerateCredentials() (*Credentials, error)
}
