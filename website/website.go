package website

type Website interface {
	Get(url string) error
	Parse() (string, error)
}
