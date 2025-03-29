package internal

type Book interface {
	Id() (int, error)
	Name() (string, error)
	URL() (string, error)
	Language() (string, error)
}
