package chap3

import (
	"fmt"
	"time"
)

type chapter3_4 struct{}

func NewChapter3_4() *chapter3_4 {
	return &chapter3_4{}
}

// BasicSelect
// p.80
func (chapter3_4) BasicSelect() {
	start := time.Now()
	c5 := make(chan interface{})
	c3 := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c5)
	}()

	go func() {
		time.Sleep(3 * time.Second)
		close(c3)
	}()

	println("Blocking on read...")

	// どれかだけが実行される
	// case 文と異なり、該当するcaseがなければ実行されない
	select {
	case <-c5:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	case <-c3:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

// BasicSelectWithChannels
// Goのランタイムはselectに対して疑似乱数による一様選択をしていることがわかる
// p.80
func (chapter3_4) BasicSelectWithChannels() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	// 大体1:1くらいの比になるようになっているが、正確ではない
	// c1Count: 545
	// c1Count: 456
	fmt.Printf("c1Count: %d\nc1Count: %d", c1Count, c2Count)
}

// ForSelectStatementWithDefault
// for文の中であるgoroutineの結果を待つ間に、別のgoroutineが仕事を進められる
// 通常、selectのdefault説はfor文のなかで使われる
// p.83
func (chapter3_4) ForSelectStatementWithDefault() {
	done := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
	// forループ自体にエイリアス？的なをつけられる？
loop:
	for {
		select {
		case <-done:
			break loop
		// default がないと必ずdoneのclose待ちをGoのランタイムが選択し、loopが回らない
		default:
		}

		// simulate
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achived %v cycles of work before signalled to stop.\n", workCounter)

}
