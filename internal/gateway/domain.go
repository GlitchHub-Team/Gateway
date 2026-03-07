package gateway

type Response struct {
	Success bool
	Message string
}

type BaseCommand interface {
	Execute() error
}
