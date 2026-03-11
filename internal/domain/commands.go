package domain

type BaseCommand interface {
	Execute() error
	String() string
}
