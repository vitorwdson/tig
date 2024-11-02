package diff

import (
	"bufio"
	"io"
	"slices"
)

func TextFileDiff(old, new io.Reader) ([]Diff, error) {
	oldScanner := bufio.NewScanner(old)
	oldScanner.Split(bufio.ScanLines)
	oldLines := make([][]byte, 0, 10)
	for oldScanner.Scan() {
		oldLines = append(oldLines, oldScanner.Bytes())
	}

	newScanner := bufio.NewScanner(new)
	newScanner.Split(bufio.ScanLines)
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
