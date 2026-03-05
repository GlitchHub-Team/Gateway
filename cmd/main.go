package main

import (
	"go.uber.org/fx"
)

func main() {
	fx.New(
	// ci vanno i provider, NewSpawnerService
	).Run()
}
