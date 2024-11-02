package diff

type DiffType int

const (
	ADD DiffType = iota
	REMOVE
	REPLACE
)

type Diff struct {
	Block   int
	Type    DiffType
	Content []byte
}
