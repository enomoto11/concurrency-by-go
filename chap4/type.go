package chap4

type Chapter4 struct {
	c4_1 *chapter4_1
	c4_3 *chapter4_3
	c4_4 *chapter4_4
}

// Number
// 章番号を返す
func (Chapter4) Number() int { return 4 }

// NewChapter4
// mainで呼び出されるファクトリ
func NewChapter4() *Chapter4 {
	return &Chapter4{
		c4_1: NewChapter4_1(),
		c4_3: NewChapter4_3(),
		c4_4: NewChapter4_4(),
	}
}

// Exec
// 当該章の処理を実行する
func (c Chapter4) Exec() error {
	c.c4_3.AvoidBlockBySignalForDoneReadingAllValuesInChannel()
	// fmt.Println("end execScripts num goroutine: ", runtime.NumGoroutine())
	return nil
}
