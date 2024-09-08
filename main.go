package main

import (
	"fmt"

	"github.com/enomoto11/concurrency-by-go/chap4"
)

type Chapter interface {
	Exec() error
	Number() int
}

func newChapters() []Chapter {
	// chapter3 := chap3.NewChapter3()
	chapter4 := chap4.NewChapter4()
	return []Chapter{
		// chapter3,
		chapter4,
	}
}

func main() {
	chapters := newChapters()
	for _, c := range chapters {
		fmt.Println("CHAPTER", c.Number(), "------------------------")
		if err := c.Exec(); err != nil {
			fmt.Printf("err in chapter %d: %v\n", c.Number(), err)
		}
		fmt.Println("CHAPTER", c.Number(), "------------------------")
	}
}
