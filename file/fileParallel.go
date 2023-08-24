package file

import (
	"bytes"
	"os"
	"strings"
)

// ReadLastLinesWithOffsetP reads the last initBufSize bytes in front of the fileOffset bytes before EOF
// returns the lines (complete and segmented) in reverse order and an error if any
// Note: the initBufSize needs to be longer than the maximum length of a log line
func ReadLastLinesWithOffsetP(fileName string, fileOffset int64, initBufSize int) ([][]byte, error) {
	if initBufSize == 0 {
		return [][]byte{}, nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	fileSize := stat.Size()
	// the fileOffset has exceeded the size of the file,
	// i.e., we've scanned through the whole file
	if fileSize <= fileOffset {
		return [][]byte{}, nil
	}

	curBufStart := fileSize - fileOffset - int64(initBufSize)
	bufSize := initBufSize

	if fileSize-fileOffset < int64(initBufSize) {
		// the remainder of the file is not big enough for a full buffer size
		curBufStart = 0
		bufSize = int(fileSize - fileOffset)
	}

	buf := make([]byte, bufSize)
	// ReadAt can't start from the very beginning of a file
	if curBufStart == 0 {
		_, err = file.Read(buf)
	} else {
		_, err = file.ReadAt(buf, curBufStart)
	}
	if err != nil {
		panic(err)
	}

	return RevertBufferByLineBreak(buf), nil
}

// RevertBufferByLineBreak split the input buf by line break byte
// Return the segments in reverse order
func RevertBufferByLineBreak(buf []byte) [][]byte {
	// find all the line break's locations (their index within the buffer)
	offset := 0
	indices := []int64{}
	for {
		index := bytes.IndexByte(buf[offset:], '\n')
		// no more line breaks
		if index == -1 {
			break
		}

		indices = append(indices, int64(index+offset))
		offset += index + 1
	}

	// no line break in the buffer
	// append the whole thing and return
	if len(indices) == 0 {
		return [][]byte{buf}
	}

	reverseBuf := [][]byte{}
	// if the last line break is not at the end of the buffer
	// append the segment from last line break to the end of buffer first
	if indices[len(indices)-1] != int64(len(buf)-1) {
		reverseBuf = append(reverseBuf, buf[indices[len(indices)-1]+1:])
	}

	// append bytes between two adjacent line breaks (a complete log line)
	for i := 0; i < len(indices)-1; i++ {
		endLineBreakIndex := indices[len(indices)-1-i]
		startLineBreakIndex := indices[len(indices)-2-i]

		// Here we include the second line break as part of the output instead of stripping it
		// this information is useful when we try to combine outputs from two ReadLastLinesWithOffsetP
		// if the first neighboring segment ends with a line break, then we append
		// if it doesn't we concatenate instead. See implementation of CombineLines
		reverseBuf = append(reverseBuf, buf[startLineBreakIndex+1:endLineBreakIndex+1])
	}

	reverseBuf = append(reverseBuf, buf[:indices[0]+1])

	return reverseBuf
}

// CombineLines combine content of the two adjacent buffers into one
// the content of the buffer is expected to be already in reverse order
// first and second refers to where they will be in the output (first will come before second)
// in the original input log file, first should be immediately after the second buffer
func CombineLines(first, second [][]byte) [][]byte {
	if len(first) == 0 {
		return second
	}

	if len(second) == 0 {
		return first
	}

	if second[0][len(second[0])-1] == '\n' {
		return append(first, second...)
	}

	jointLine := append(second[0], first[len(first)-1]...)
	if len(second) == 1 {
		return append(first[:len(first)-1], jointLine)
	}

	return append(append(first[:len(first)-1], jointLine), second[1:]...)
}

const FILE_OFFSET_UNIT_SIZE = 1 << 15

// ReadLastNLinesWithKeywordP keeps calling ReadLastLinesWithOffsetP until we reach the target lines of log
// if input query is not empty, log lines are filtered first before they are appended
func ReadLastNLinesWithKeywordP(fileName string, n int, query string) ([][]byte, error) {
	initBufSize := FILE_OFFSET_UNIT_SIZE

	lines := [][]byte{}
	fileOffsetCounter := 0
	newlines := [][]byte{[]byte("dummy")}
	var err error

	// if we have a filter key word, things get a bit tricky because
	// we can't blindly combine all the results with CombineLines
	// instead we need to filter the lines that contain the keyword before appending to the result
	// given that the lines at the border of two buffers might both be segmented
	// we need to keep a record of the last line of the previous result
	// and join this last line with the new result (bar new last line) and search for keyword
	// the new last line will be stored for next time when we have new results
	rollingLastLine := []byte{}

	// try to read one more line in case the last line is segmented
	for len(lines) < n+1 {
		newlines, err = ReadLastLinesWithOffsetP(
			fileName, int64(fileOffsetCounter*initBufSize), initBufSize)
		if err != nil {
			panic(err)
		}

		if len(newlines) == 0 {
			// if there's content in the rolling last line and we don't have any newlines coming
			// then we've reached the beginning of the file
			// we need to append the content if it matches the query
			if len(rollingLastLine) > 0 {
				if query == "" || strings.Contains(string(rollingLastLine), query) {
					lines = CombineLines(lines, [][]byte{rollingLastLine})
				}
			}
			break
		}

		if query == "" {
			lines = CombineLines(lines, newlines)
		} else {
			combinedNewLines := CombineLines([][]byte{rollingLastLine}, newlines[:len(newlines)-1])
			rollingLastLine = newlines[len(newlines)-1]

			for _, newline := range combinedNewLines {
				if strings.Contains(string(newline), query) {
					lines = append(lines, newline)
				}
			}
		}

		fileOffsetCounter++
	}

	if len(lines) >= n {
		return lines[:n], nil
	}

	return lines, nil
}
