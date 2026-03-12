package modules

import (
	credentialsgenerator "Gateway/internal/credentialsGenerator"

	"go.uber.org/fx"
)

var CredGenerator = fx.Module("CredGenerator",
	fx.Provide(
		fx.Annotate(
			credentialsgenerator.NewNKeyCredentialsGenerator,
			fx.As(new(credentialsgenerator.CredentialsGeneratorPort)),
		),
	),
)
