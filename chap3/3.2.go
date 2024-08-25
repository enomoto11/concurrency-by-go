package chap3

import (
	"fmt"
	"sync"
	"time"
)

type chapter3_2 struct{}

func NewChapter3_2() *chapter3_2 {
	return &chapter3_2{}
}

// AddingAndRemovingQueue
// p.54
func (chapter3_2) AddingAndRemovingQueue() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)

		c.L.Lock() // このロックはqueue変数に複数goroutineがアクセスする可能性があるため
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal() // Signal()をしないと、Cond型に子goroutineで目的を達成できたことを知らせることができず、ループ内のwait()がずっと待機したままになる
	}

	for i := 0; i < 10; i++ {
		c.L.Lock() // こちらのロックも、これから使おうとしているqueue変数に複数goroutineがアクセスする可能性があるため

		for len(queue) == 2 {
			// queueの中身が2個にまで達してしまったらmain goroutine上ではwaitする
			// その間子goroutineでremovingが行われるため、いずれこのループからは抜けられる
			c.Wait()
		}

		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{}) // Addingは常にmain goroutineで行われる

		go removeFromQueue(1 * time.Second) //removingは常に子goroutineで行われ、それぞれのremovingは別々のgoroutine上である

		c.L.Unlock()
	}
}

// ClickAndBroadCastToGoroutines
// p.56
func (chapter3_2) ClickAndBroadCastToGoroutines() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait() // シグナル（というよりはbroadcastされる情報）を受け取るまでは待機し、後続の処理を実行しない
			fn()
		}()

		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("Maxmizing window")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	fmt.Println("ブロードキャスト中...")
	clickRegistered.Wait()
	fmt.Println("ブロードキャスト完了")
}
