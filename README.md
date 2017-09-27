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
w := round.NewSpinMe(os.Stdout, round.Pipe)
```
os.Stdout could be any io.Writer with an Fd. If it's not a terminal the spinner does nothing.

## Write to It
```go
fmt.Fprintln(&w, "Like a record!")
```
Don't write directly to the underlying writer before the spinner is closed. There be dragons.

## Make It Stop
```go
w.Close()
```
It _does not_ close the underlying writer.

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
w := round.NewSpinMe(os.Stdout, s)
```
