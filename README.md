# round
A Go command line spinner library.

# Get It

```go
import "github.com/dedelala/round"
```

# Make It Go

```go
w := round.NewSpinMe(os.Stdout, round.Pipe)
defer w.close()
```

# Built-In Styles!

Style    | Unicode Set
-------- | -----------
`Pipe`   | 0020—007F Basic Latin
`Moon`   | 1F300—1F5FF Miscellaneous Symbols and Pictographs
`Block`  | 2580—259F Block Elements
`Hearts` | 1F300—1F5FF Miscellaneous Symbols and Pictographs

# Make Your Own Scroller

```go
s := round.NewScroller(12, "[%v]", "Do Not Panic!")
w := round.NewSpinMe(os.Stdout, s)
```
