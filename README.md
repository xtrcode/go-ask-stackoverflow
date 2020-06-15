# go-ask-stackoverflow
go-ask-stackoverflow is a simple terminal appilcation which quickly responds to your question with the accepted answer from StackOverflow. As we programmers tend to love StackOverflow for inspiration or just quick & dirty answers, go-ask-stackoverflow intent is to speed this up and be easily extendable.

## Clone

```bash
$ git clone github.com/xtrcode/go-ask-stackoverflow
$ cd go-ask-stackoverflow/
$ go build -o ask main.go
```
Now you can move the executable to your preferred location (e.g. `/usr/bin`) and just

## Enjoy
`& ask whatever you want`

## Customize to your needs
### Add custom search engine

```go
type Engine interface {
	Request(str string) error
	Get() (string, error)
}
```

### Add custom website parser

```go
type Website interface {
	Get(url string) error
	Parse() (string, error)
}
```

### Add custom cache driver
```go
type Cache interface {
	Init() error
	Open() error
	Close() error
	Get(key string) (string, error)
	Set(key, value string) error
}
```

# Credits
go-ask-stackoverflow is inspired by https://github.com/juliusmh/ask

# LICENSE
MIT
