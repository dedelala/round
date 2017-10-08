# round
Go command line spinner.

### go-round
go-round is a command line program that can add a spinner to things. It copies
stdin to stdout. You can `go install github.com/dedelala/round/go-round`, and
read `go-round help` for usage.

## About
Package round is a command line spinner. Start one with Go.

If Stdout is a terminal, the spinner will be written there. If not, Stderr
will be checked. If neither, the spinner will quietly do nothing.

The exported Stdout and Stderr, if connected to a terminal, will block
against the spinner. Use these instead of the handles in package os.

## Get It
`go get github.com/dedelala/round`

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
