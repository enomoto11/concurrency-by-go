package main

import (
	"fmt"

	"github.com/enomoto11/concurrency-by-go/chap3"
)

type Chapter interface {
	Exec() error
	Number() int
}

func newChapters() []Chapter {
	chapter3 := chap3.NewChapter3()
	return []Chapter{
		chapter3,
	}
}

func main() {
	chapters := newChapters()
	for _, c := range chapters {
		if err := c.Exec(); err != nil {
			fmt.Printf("err in chapter %d: %v\n", c.Number(), err)
		}
	}
}
