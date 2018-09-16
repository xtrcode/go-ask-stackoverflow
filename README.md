# go-ask-stackoverflow
go-ask-stackoverflow is a simple terminal appilcation which quickly response to your question with the accepted answer of StackOverflow. As we programmers tend to love StackOverflow for inspiration or just quick & dirty answers, go-ask-stackoverflow intent is to speed this up and be easily extendable.

# Install
## Required Packages
To cache every question & answer go-ask-stackoverflow uses boltdb as key-value store.

`$ go get github.com/boltdb/bolt/...`

To parse Google.com and StackOverflow.com goquery is used.

`$ go get github.com/PuerkitoBio/goquery`

## Clone
`$ git clone github.com/xtrcode/go-ask-stackoverflow`

`$ cd go-ask-stackoverflow`

`$ go build -o ask main.go`

Now you can move the executable to your preferred location (e.g. `/usr/bin`) and just

## Enjoy
`& ask what every you want`

# Credits
go-ask-stackoverflow is inspired by https://github.com/juliusmh/ask

# LICENCE
MIT
