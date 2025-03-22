package internal

type Book interface {
	Id() (int, error)
	Title() (string, error)
	URL() (string, error)
	Language() (string, error)
}
