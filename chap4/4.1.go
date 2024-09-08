package chap4

import (
	"bytes"
	"fmt"
	"sync"
)

type chapter4_1 struct{}

func NewChapter4_1() *chapter4_1 {
	return &chapter4_1{}
}

// LexicalConstraintWithConcurrentlyUnsafeData
// 並行安全ではないデータ構造を使った拘束
func (chapter4_1) LexicalConstraintWithConcurrentlyUnsafeData() {
	// 1. チャネルによる同期はパフォーマンス上ネックになるためこれを回避できる
	// 2. レキシカルスコープのコンテキストの中では同期的なコードが書ける
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)

	data := []byte("golang")

	// スライスの部分集合を別のインスタンスとして渡しているので並行安全
	// 起動したgoroutineは直接元のインスタンス data にアクセスしないように拘束している
	go printData(&wg, data[:4])
	go printData(&wg, data[3:])

	wg.Wait()
}
