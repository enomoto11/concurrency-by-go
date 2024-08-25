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

// MyPool
// p.60
func (chapter3_2) MyPool() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()            // ここではインスタンス生成される
	instace := myPool.Get() // ここでもインスタンスを生成する
	myPool.Put(instace)
	myPool.Get() // ここではプール内に使用可能なインスタンスがあるため、生成処理Newは行われない
}

// AFewMemmoryAllocationsIsEnoughThanksToPool
// p. 61
func (chapter3_2) AFewMemmoryAllocationsIsEnoughThanksToPool() {
	var numCalcsCreated int

	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated++
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// 事前に4KBを確保する
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem) // 使用したインスタンスはpoolにputしなおす

			//メモリに対して任意の処理を行う
		}()
	}

	wg.Wait()
	fmt.Printf("%d calcurators ware created", numCalcsCreated) // 11や12など比較的小さな結果を得る
}
