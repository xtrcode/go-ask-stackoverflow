package cache

import (
	"os"

	"github.com/boltdb/bolt"
)

type Bolt struct {
	db         *bolt.DB
	cache_dir  string
	cache_file string
	cache_full string
}

func (b *Bolt) Init() error {
	var err error

	b.cache_dir, err = os.UserHomeDir()
	if err != nil {
		return err
	}

	b.cache_dir += "/.ask"
	b.cache_file = "cache.db"
	b.cache_full = b.cache_dir + "/" + b.cache_file

	// create cache dir and db if not exists
	if _, err := os.Stat(b.cache_dir); os.IsNotExist(err) {
		if err = os.MkdirAll(b.cache_dir, 0777); err != nil {
			return err
		}

		f, err := os.Create(b.cache_full)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	return nil
}

func (b *Bolt) Open() error {
	var err error

	b.db, err = bolt.Open(b.cache_full, 0700, nil)
	return err
}

func (b *Bolt) Close() error {
	return b.db.Close()
}

func (b *Bolt) Get(key string) (string, error) {
	var answer []byte

	if err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("42"))
		if b == nil {
			return nil
		}

		answer = b.Get([]byte(key))

		return nil
	}); err != nil {
		return "", err
	}

	return string(answer), nil
}

func (b *Bolt) Set(key, value string) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("42"))
		if err != nil {
			return err
		}

		return b.Put([]byte(key), []byte(value))
	}); err != nil {
		return err
	}

	return nil
}
