# round
A Go command line spinner library.

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

## Write to It
```go
fmt.Fprintln(round.Stdout, "Like a record!")
```
Don't write directly to os.Stdout or os.Stderr while a spinner is running. There be dragons.

## Make It Stop
```go
round.Stop()
```

## Built-In Styles!

Style    | Unicode Set
-------- | -----------
`Pipe`   | 0020—007F Basic Latin
`Moon`   | 1F300—1F5FF Miscellaneous Symbols and Pictographs
`Block`  | 2580—259F Block Elements
`Hearts` | 1F300—1F5FF Miscellaneous Symbols and Pictographs

## Make Your Own Scroller

```go
s := round.NewScroller(12, "[%v]", "Do Not Panic!")
round.Start(s)
```
