# Go言語による並行処理　ノート

## Usage

- 第N章のノートをとり始めたい時は以下のコマンドからボイラープレートを作成する
```sh
go run cmd/main.go n
```

例えば第3章のノートをとりたい時は
```sh
go run cmd/main.go 3
```
をコマンドラインから実行する

- 第N章の実装を実行したい場合は[main.go](./main.go)の #newChapters にて第N章の関数が実装されている構造体を`Chapter`構造体に格納する

例えば第3章の実装を試したい場合は
```main.go
func newChapters() []Chapter {
	chapter3 := chap3.NewChapter3()
	return []Chapter{
		chapter3,
	}
}
```
とする