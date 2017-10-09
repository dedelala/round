# round
Go command line spinner.


## About
Package round is a command line spinner. Start one with Go.

The package will intelligently decide whether to write spinners on stdout,
stderr or neither, depending if a terminal is present.

Wrappers for Stdout and Stderr are provided so as not to interfere with the
spinner while running.


## Get It
```
go get github.com/dedelala/round
```


### go-round
```
go install github.com/dedelala/round/go-round
```

[go-round](go-round/) is a command line program that adds a spinner to almost
anything. It copies stdin to stdout. You can read `go-round help` for usage.


## Use It
```go
import "github.com/dedelala/round"
```


## Make It Go
```go
round.Go(round.Pipe)
```


## Write Out
```go
fmt.Fprintln(round.Stdout, "Like a record!")
fmt.Fprintln(round.Stderr, "Wahhhhh!")
```
Don't write directly to os.Stdout or os.Stderr while a spinner is running. There be dragons.


## Make It Stop
```go
round.Stop()
```


## Built-In Styles!

Style    | =  | Unicode Set
---------|----|--------------
`Block`  | █  | 2580—259F Block Elements
`Cylon`  | @  | 0020—007F Basic Latin
`Hearts` | 💖 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
`Moon`   | 🌓 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
`Pipe`   | -  | 0020—007F Basic Latin


## Make Your Own Scroller

```go
s := round.NewScroller(12, "[%v]", "Do Not Panic!")
round.Start(s)
```
