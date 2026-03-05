package commands

type BaseCommand interface {
	Execute() error
}
