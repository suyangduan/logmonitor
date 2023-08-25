package file

import (
	"bytes"
	"os"
	"strings"
)

type LineReturn struct {
	Line   string
	Offset int64
}

// ReadLastLinesWithOffsetPagination reads the last initBufSize bytes in front of the fileOffset bytes before EOF
// fileOffset needs to be at a line break. otherwise the incomplete line at the end of the buffer will be lost
// returns the complete lines in reverse order and their individual line's file offset and an error if any
// Note: the initBufSize needs to be longer than the maximum length of a log line
func ReadLastLinesWithOffsetPagination(fileName string, fileOffset int64, initBufSize int) ([]LineReturn, error) {
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
		return []LineReturn{}, nil
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

	// validate that the input fileOffset value is correct
	if buf[len(buf)-1] != '\n' {
		panic("the last byte in the buffer is not a line break")
	}

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

	lastLines := []LineReturn{}
	for i := 0; i < len(indices)-1; i++ {
		// between two adjacent line breaks is a complete line
		endLineBreakIndex := indices[len(indices)-1-i]
		startLineBreakIndex := indices[len(indices)-2-i]

		lastLines = append(lastLines, LineReturn{
			Line:   string(buf[startLineBreakIndex+1 : endLineBreakIndex]),
			Offset: int64(bufSize) - startLineBreakIndex - 1 + fileOffset,
		})
	}

	// if this is the beginning of the file, append the first line
	// which starts at the beginning of the buffer and ends at the first line break
	if curBufStart == 0 {
		lastLines = append(lastLines, LineReturn{
			Line:   string(buf[:indices[0]]),
			Offset: fileSize,
		})
		// set fileOffset to be file size, indicating we've scanned through
		return lastLines, nil
	}

	// the new offset indicates the first line break's location from the end of the file
	// which will be the end location of next round of buffer reading
	return lastLines, nil
}

// ReadLastNLinesWithKeywordPagination keeps calling ReadLastLinesWithOffsetPagination
// until we reach the target lines of log.
// if input query is not empty, log lines are filtered first before they are appended
func ReadLastNLinesWithKeywordPagination(fileName string, n int, query string, offset int64) ([]string, int64, error) {
	return ReadLastNLinesWithKeywordPaginationInternal(fileName, n, query, offset, READ_BUFFER_SIZE)
}

// ReadLastNLinesWithKeywordPaginationInternal exposes buffer size for testing only
// internal only!!!
func ReadLastNLinesWithKeywordPaginationInternal(
	fileName string, n int, query string, offset int64, initBufSize int) ([]string, int64, error) {
	lines := []LineReturn{}
	newlines := []LineReturn{{"", offset}}
	var err error
	for len(lines) < n && len(newlines) != 0 {
		newlines, err = ReadLastLinesWithOffsetPagination(fileName, newlines[len(newlines)-1].Offset, initBufSize)
		if err != nil {
			panic(err)
		}

		if query == "" {
			lines = append(lines, newlines...)
		} else {
			for _, newline := range newlines {
				if strings.Contains(newline.Line, query) {
					lines = append(lines, newline)
				}
			}
		}
	}

	returnSize := n
	if len(lines) < n {
		returnSize = len(lines)
	}

	returnVal := []string{}

	for i := 0; i < returnSize; i++ {
		returnVal = append(returnVal, lines[i].Line)
	}

	return returnVal, lines[returnSize-1].Offset, nil
}
