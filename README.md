This is a hodgepodge of Go-lang for working with
Mars 2020 Mission Perseverance Rover [raw images](https://mars.nasa.gov/mars2020/multimedia/raw-images/).
It's a mess!  And it's a learning exercise, to re-acquaint myself with [Go](https://golang.org).

This code is based loosely on my [mars_perseverance_images](https://mchapman87501@github.com/mchapman87501/mars_perseverance_images.git) repo.  In addition to helping with Go fluency, it has provided a chance to learn basic image processing tasks such as de-mosaicing images taken under a Bayer filter.

The repo contains a few separate Go modules.  `lib` holds code used by the command-line applications in `app`.  The various `go.mod` files contain `replace` directives that assume everything is compiled from a common git repository.
