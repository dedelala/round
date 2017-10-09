# go-round

## install

```
go install github.com/dedelala/round/go-round
```

## about

```
go-round help
```

```
go-round - copies stdin to stdout and shows a spinner
Usage: go-round [style]

The default style is pipe.

Style   | =  | Unicode Set
--------|----|--------------
block   | █  | 2580—259F Block Elements
cylon   | @  | 0020—007F Basic Latin
hearts  | 💖 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
moon    | 🌓 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
pipe    | -  | 0020—007F Basic Latin

Scrollers
Usage: go-round [options] [scroll|bounce] [message...]

  -f string
    	format for a scroller or bouncer frame (default "[%v]")
  -w int
    	field width of a scroller or bouncer (default 8)
```

## examples

```
ping google.com |go-round
```

```
go get -u github.com/dedelala/round |go-round hearts
```

## just for fun

Note: go-round can be redirected and it still goes round

```
go-round > some-file
# (then type stuff)
# Ctrl-D
```
