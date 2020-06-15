package engine

type Engine interface {
	Request(str string) error
	Get() (string, error)
}
