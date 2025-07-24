package ports

type ResultScore interface {
	SaveResult(a, b, result int) error
}