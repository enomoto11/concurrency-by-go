package chap3

type Chapter3 struct {
	chapter3_2 *chapter3_2
}

// Number
// 章番号を返す
func (Chapter3) Number() int { return 3 }

// NewChapter3
// mainで呼び出されるファクトリ
func NewChapter3() *Chapter3 {
	return &Chapter3{
		chapter3_2: NewChapter3_2(),
	}
}

// Exec
// 当該章の処理を実行する
func (c Chapter3) Exec() error {
	// c.chapter3_2.ClickAndBroadCastToGoroutines()
	// c.chapter3_2.AddingAndRemovingQueue()
	// c.chapter3_2.MyPool()
	c.chapter3_2.AFewMemmoryAllocationsIsEnoughThanksToPool()

	return nil
}
