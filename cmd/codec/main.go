package main

import (
	"io"
	"os"
)

const KEY = "k?P9KnIN5YM~kWLJcthi,gOnG2RT-qY00Yq-TR2GnOg,ihtcJLWk~MY5NInK9P?k"

type Codec struct {
	r        io.Reader
	keyBytes []byte
	keyPos   int
}

func (c *Codec) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	for i := range n {
		p[i] ^= c.keyBytes[c.keyPos]
		c.keyPos = (c.keyPos + 1) % len(c.keyBytes)
	}
	return
}

func New(rdr io.Reader, key []byte) Codec {
	return Codec{
		r:        rdr,
		keyBytes: key,
		keyPos:   0,
	}
}

func main() {
	c := New(os.Stdin, []byte(KEY))
	io.Copy(os.Stdout, &c)
}
