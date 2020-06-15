package cache

type Cache interface {
	Init() error
	Open() error
	Close() error
	Get(key string) (string, error)
	Set(key, value string) error
}
