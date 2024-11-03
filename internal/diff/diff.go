package diff

import (
	"bufio"
	"io"
	"slices"
)

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

func fileDiff(old, new io.Reader, splitFunc bufio.SplitFunc) ([]Diff, error) {
	oldScanner := bufio.NewScanner(old)
	oldScanner.Split(splitFunc)
	oldBlocks := make([][]byte, 0, 10)
	for oldScanner.Scan() {
		oldBlocks = append(oldBlocks, oldScanner.Bytes())
	}

	newScanner := bufio.NewScanner(new)
	newScanner.Split(splitFunc)
	newBlocks := make([][]byte, 0, 10)
	for newScanner.Scan() {
		newBlocks = append(newBlocks, newScanner.Bytes())
	}

	var oldEnded, newEnded bool
	var oldContent, newContent []byte
	currOld, currNew := -1, -1
	totalOld := len(oldBlocks)
	totalNew := len(newBlocks)
	diffs := make([]Diff, 0, 10)

	for {
		oldEnded = currOld+1 >= totalOld
		newEnded = currNew+1 >= totalNew

		if oldEnded && newEnded {
			break
		}

		if !oldEnded {
			currOld++
			oldContent = oldBlocks[currOld]
		}

		if !newEnded {
			currNew++
			newContent = newBlocks[currNew]
		}

		if slices.Equal(oldContent, newContent) {
			continue
		}

		if oldEnded && !newEnded {
			diffs = append(diffs, Diff{
				Block:   currNew,
				Type:    ADD,
				Content: newContent,
			})
			continue
		} else if !oldEnded && newEnded {
			diffs = append(diffs, Diff{
				Block:   currOld,
				Type:    REMOVE,
				Content: oldContent,
			})
			continue
		}

		hit := -1
		if !oldEnded {
			for i := currOld + 1; i < totalOld; i++ {
				if slices.Equal(oldBlocks[i], newContent) {
					hit = i
					break
				}
			}
			if hit != -1 {
				for i := currOld; i < hit; i++ {
					diffs = append(diffs, Diff{
						Block:   i,
						Type:    REMOVE,
						Content: oldBlocks[i],
					})
				}
				currOld = hit - 1
				currNew--
				continue
			}
		}

		if !newEnded {
			for i := currNew + 1; i < totalNew; i++ {
				if slices.Equal(newBlocks[i], oldContent) {
					hit = i
					break
				}
			}
			if hit != -1 {
				for i := currNew; i < hit; i++ {
					diffs = append(diffs, Diff{
						Block:   i,
						Type:    ADD,
						Content: newBlocks[i],
					})
				}
				currNew = hit - 1
				currOld--
				continue
			}
		}

		diffs = append(diffs, Diff{
			Block:   currNew,
			Type:    REPLACE,
			Content: newContent,
		})
	}

	return diffs, nil
}

func TextFileDiff(old, new io.Reader) ([]Diff, error) {
	return fileDiff(old, new, bufio.ScanLines)
}

func BinaryFileDiff(old, new io.Reader) ([]Diff, error) {
	return fileDiff(
		old,
		new,
		func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			amount := min(len(data), 512)
			return amount, data[0:amount], nil
		},
	)
}
