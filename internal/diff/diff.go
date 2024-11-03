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
	oldLines := make([][]byte, 0, 10)
	for oldScanner.Scan() {
		oldLines = append(oldLines, oldScanner.Bytes())
	}

	newScanner := bufio.NewScanner(new)
	newScanner.Split(splitFunc)
	newLines := make([][]byte, 0, 10)
	for newScanner.Scan() {
		newLines = append(newLines, newScanner.Bytes())
	}

	var oldEnded, newEnded bool
	currOld, currNew := -1, -1
	totalOld := len(oldLines)
	totalNew := len(newLines)
	diffs := make([]Diff, 0, 10)

	for {

		oldEnded = currOld+1 >= totalOld
		newEnded = currNew+1 >= totalNew

		if oldEnded && newEnded {
			break
		}

		if !oldEnded {
			currOld++
		}

		if !newEnded {
			currNew++
		}

		if slices.Equal(oldLines[currOld], newLines[currNew]) {
			continue
		}

		hit := -1
		if !oldEnded {
			for i := currOld + 1; i < totalOld; i++ {
				if slices.Equal(oldLines[i], newLines[currNew]) {
					hit = i
					break
				}
			}
			if hit != -1 {
				for i := currOld; i < hit; i++ {
					diffs = append(diffs, Diff{
						Block:   i,
						Type:    REMOVE,
						Content: oldLines[i],
					})
				}
				currOld = hit - 1
				currNew--
				continue
			}
		}

		if !newEnded {
			for i := currNew + 1; i < totalNew; i++ {
				if slices.Equal(newLines[i], oldLines[currOld]) {
					hit = i
					break
				}
			}
			if hit != -1 {
				for i := currNew; i < hit; i++ {
					diffs = append(diffs, Diff{
						Block:   i,
						Type:    ADD,
						Content: newLines[i],
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
			Content: newLines[currNew],
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
