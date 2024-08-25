package chap3

import (
	"fmt"
	"sync"
)

type chapter3_3 struct{}

func NewChapter3_3() *chapter3_3 {
	return &chapter3_3{}
}

// ReadAfterChannelWasClosed
// close()というのは書き込みに対するロック
// closeされたチャネルに何も事前に送信されていなくても読み取ることができる
// p.69
func (chapter3_3) ReadAfterChannelWasClosed() {
	stream := make(chan int)
	close(stream)
	// closeがない場合、一切書き込みがないチャネルを読み取ろうとしてしまい、deadlockを起こす
	// これはなぜなら、Goのチャネルはブロックの挙動をとるから

	integer, ok := <-stream

	fmt.Printf("(%v): %v", ok, integer)
}

// RangeStatementWithChannel
// p.70
func (chapter3_3) RangeStatementWithChannel() {
	intStream := make(chan int)
	go func() { // チャネルに書き込む専用のgoroutine
		defer close(intStream)

		for i := 1; i < 5; i++ {
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Printf("%v", integer)
	}
}

// UnblockGoroutinesSimultaneously
// p.70
func (chapter3_3) UnblockGoroutinesSimultaneously() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			<-begin // 書き込みがされていないチャネルの読み取りを行うことによって、当該処理をブロックする
			fmt.Printf("%d has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin) // closeによってブロックされているチャネルの読み取り部分に対してシグナルを送信することができる
	// closeされたチャネルは無限に読み取り可能となるので、チャネルの読み取りによってブロックされている各goroutine上の残りの処理はきちんと実行されるようになる
	// これはつまり、同時に複数のgoroutineを解放していることと同義となる
	wg.Wait()
}
