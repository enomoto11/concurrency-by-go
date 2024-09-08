package chap4

import (
	"fmt"
	"math/rand"
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

// BlockGoroutinesTryingToWriteOnChannel
// これはアンチパターン
func (chapter4_3) BlockGoroutinesTryingToWriteOnChannel() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)

		// このサブルーチンのfor文が終わらないため、deferの処理が実行されない
		// for文が終わらないのは、読み込みがすでにメインgoroutineで完了しきっているチャネルに送信しているため
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()
		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints: ")

	// このメインgoroutineにおいては、3回の読み込み処理が発生するのは約束されている
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}

// TryToWriteAfterCompletelyReadChannel
// これはアンチパターン
func (chapter4_3) TryToWriteAfterCompletelyReadChannel() {
	c := make(chan interface{})
	c <- "test"
	c <- "test"
	c <- "test"

	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-c)
	}

	fmt.Println(<-c) // ここまでの時系列的で並行的にチャネルに送信をしようとするサブルーチンがないため、これもブロックされる（永遠に待つことになる）
	fmt.Println("ここには到達しないはず")

	// この時点で、自分以外のgoroutineにてこれ以上読み込み処理がなされないと確定しているチャネルに対して書き込みを行おうとするからブロックされる
	// この場合、そもそも他にサブルーチンもないのでGoランタイムがブロックして読み込み待ちするのを諦めてdead lockとなる
	c <- "test"
	fmt.Println("ここには到達しないはず")

	fmt.Println(<-c) // 同じgoroutine上で書き込み回数と同等の数の読み込みが完了しているならば、問題の書き込み部分より後に読み込みを書いても同じくdead lockを起こす
	fmt.Println("ここにも到達しないはず")
}

// AvoidBlockBySignalForDoneReadingAllValuesInChannel
// p.95
func (chapter4_3) AvoidBlockBySignalForDoneReadingAllValuesInChannel() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				// どこで書き込みを止めるべきかをメインgoroutineから教えてもらうことで、
				// これから読み込まれる見込みのないチャネルに送信するのを防ぐことができる
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan interface{})

	randStream := newRandStream(done)
	fmt.Println("3 random ints: ")

	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done) // 読み込みが完了したため、これ以上書き込みをしないようにサブルーチンにシグナルを送信する

	// 処理が実行中であることをエミュレート
	// これがない場合、サブルーチンよりも先にメインルーチンが終了してしまうため、期待する出力が得られない
	time.Sleep(1 * time.Second)
}
