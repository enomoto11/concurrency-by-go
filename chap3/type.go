package chap3

type Chapter3 struct {
	c3_2 *chapter3_2
	c3_3 *chapter3_3
	c3_4 *chapter3_4
}

// Number
// 章番号を返す
func (Chapter3) Number() int { return 3 }

// NewChapter3
// mainで呼び出されるファクトリ
func NewChapter3() *Chapter3 {
	return &Chapter3{
		c3_2: NewChapter3_2(),
		c3_3: NewChapter3_3(),
		c3_4: NewChapter3_4(),
	}
}

// Exec
// 当該章の処理を実行する
func (c Chapter3) Exec() error {
	// c.c3_2.ClickAndBroadCastToGoroutines()
	// c.c3_2.AddingAndRemovingQueue()
	// c.c3_2.MyPool()
	// c.c3_2.AFewMemmoryAllocationsIsEnoughThanksToPool()
	// c.c3_3.ReadAfterChannelWasClosed()
	// c.c3_3.RangeStatementWithChannel()
	// c.c3_3.UnblockGoroutinesSimultaneously()
	// c.c3_3.EncapsulationChannelInProducerGoroutine()
	// c.c3_4.BasicSelect()
	c.c3_4.ForSelectStatementWithDefault()

	return nil
}
