package diff_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/vitorwdson/tig/internal/diff"
)

var (
	oldFile = `abc
def
ghi
jkl
mno
pqr
stu
vwx
yz
`
	newFile = `abc-
def
---
---
ghi
jkl
stu
vwx-
yz
`
)

func TestTextDiff(t *testing.T) {
	diffs, err := diff.TextFileDiff(strings.NewReader(oldFile), strings.NewReader(newFile))
	if err != nil {
		t.Fatal("generating text file diffs errored: ", err)
	}

	expected := []diff.Diff{
		{Block: 0, Type: diff.REPLACE, Content: []byte("abc-")},
		{Block: 2, Type: diff.ADD, Content: []byte("---")},
		{Block: 3, Type: diff.ADD, Content: []byte("---")},
		{Block: 4, Type: diff.REMOVE, Content: []byte("mno")},
		{Block: 5, Type: diff.REMOVE, Content: []byte("pqr")},
		{Block: 7, Type: diff.REPLACE, Content: []byte("vwx-")},
	}

	if len(expected) != len(diffs) {
		t.Fatalf("expected %v, but got %v", expected, diffs)
	}

	for i := range expected {
		e := expected[i]
		d := diffs[i]

		if e.Block != d.Block || e.Type != d.Type || !slices.Equal(e.Content, d.Content) {
			t.Fatalf("expected %v, but got %v", expected, diffs)
		}
	}
}

func TestEmptyTextDiff(t *testing.T) {
	diffs, err := diff.TextFileDiff(strings.NewReader(""), strings.NewReader(newFile))
	if err != nil {
		t.Fatal("generating text file diffs errored: ", err)
	}

	expected := []diff.Diff{
		{Block:0, Type: diff.ADD, Content: []byte("abc-")},
		{Block:1, Type: diff.ADD, Content: []byte("def")},
		{Block:2, Type: diff.ADD, Content: []byte("---")},
		{Block:3, Type: diff.ADD, Content: []byte("---")},
		{Block:4, Type: diff.ADD, Content: []byte("ghi")},
		{Block:5, Type: diff.ADD, Content: []byte("jkl")},
		{Block:6, Type: diff.ADD, Content: []byte("stu")},
		{Block:7, Type: diff.ADD, Content: []byte("vwx-")},
		{Block:8, Type: diff.ADD, Content: []byte("yz")},
	}

	if len(expected) != len(diffs) {
		t.Fatalf("expected %v, but got %v", expected, diffs)
	}

	for i := range expected {
		e := expected[i]
		d := diffs[i]

		if e.Block != d.Block || e.Type != d.Type || !slices.Equal(e.Content, d.Content) {
			t.Fatalf("expected:\n%vbut got:\n%v", e, d)
		}
	}
}
