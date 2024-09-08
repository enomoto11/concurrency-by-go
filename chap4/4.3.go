package chap4

import (
	"fmt"
	"time"
)

type chapter4_3 struct{}

func NewChapter4_3() *chapter4_3 {
	return &chapter4_3{}
}

// GoroutineLeak
// p.93
// nilチャネルを渡して、サブルーチンで読み込みをしてしまうためそのgoroutineは破棄されない
func (chapter4_3) GoroutineLeak() {
	dowork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})

		// メインスレッドが起動している限りこのgoroutineはメモリに残り続ける
		go func() {
			defer fmt.Println("dowrk exited")
			defer close(completed)

			for s := range strings { // 実際に渡されるのはnilであり、読み込みまちが発生するためfor文の中身に処理が移ることはない
				// do something
				fmt.Println(s)
			}
		}()

		return completed
	}

	dowork(nil)
	fmt.Println("done")
}

// AvoidGoroutineLeakBySignalingBetweenParentAndChildGoroutines
// 親（→）子goroutine間で(慣習的にdoneと呼ばれる)チャネルをシグナルとして共有する
// p.93
func (chapter4_3) AvoidGoroutineLeakBySignalingBetweenParentAndChildGoroutines() {
	dowork := func(
		done chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})

		go func() {
			defer fmt.Println("dowork exited")
			defer close(terminated)

			for {
				select {
				case s := <-strings:
					// do something
					_ = s
				case <-done: // メモリリークする可能性のあるループ内で、親からのシグナルに従うようにする
					return
				}
			}
		}()

		return terminated
	}

	done := make(chan interface{})
	// これはdoneと違い、シンプルに子が親に完了を通達するシグナルなのでメモリリークする場合でも書いてる
	terminated := dowork(done, nil)

	// これは別に別goroutineである必要はない
	go func() {
		// 1秒後に操作をキャンセルする
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelling dowork goroutine...")
		close(done) // 子goroutineに伝達される
	}()

	<-terminated
	fmt.Println("done")
}
