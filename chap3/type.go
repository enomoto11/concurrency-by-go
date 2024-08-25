package chap3

type Chapter3 struct{}

func (Chapter3) Number() int { return 3 }

func NewChapter3() *Chapter3 {
	return &Chapter3{}
}
