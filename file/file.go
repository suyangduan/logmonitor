package file

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func WriteFileHelper(fileName string, numOfLines int) {
	fo, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	for ; numOfLines > 0; numOfLines-- {
		newLine := fmt.Sprintf("%s\n", time.Now())
		if _, err := fo.Write([]byte(newLine)); err != nil {
			panic(err)
		}
	}
}

func GenerateTestFile(fileName string, numOfLines int) {
	fo, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	animals := []string{"mouse", "ox", "tiger", "rabbit", "dragon", "snake", "horse", "ram", "monkey", "rooster", "dog", "pig"}
	for i := 0; i < numOfLines; i++ {

		newLine := fmt.Sprintf("%s%d%s%s\n",
			"Line ", i+1, " this line contains a random animal: ", animals[rand.Intn(12)])
		if _, err := fo.Write([]byte(newLine)); err != nil {
			panic(err)
		}
	}
}

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

const READ_BUFFER_SIZE = 128

// ReadLastLines reads a buffer size of READ_BUFFER_SIZE at the end of a file
// returns full lines within that buffer in reverse order and an error if any
func ReadLastLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, READ_BUFFER_SIZE)
	stat, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	fileSize := stat.Size()
	curBufStart := fileSize - READ_BUFFER_SIZE
	if fileSize < READ_BUFFER_SIZE {
		curBufStart = 0
	}
	_, err = file.ReadAt(buf, curBufStart)
	if err != nil {
		panic(err)
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

	lastLines := []string{}
	for i := 0; i < len(indices)-1; i++ {
		// create a buffer the size between two line breaks (excluding the second line break)
		// and then read everything in between
		endLineBreakIndex := indices[len(indices)-1-i]
		startLineBreakIndex := indices[len(indices)-2-i]
		bufferSize := endLineBreakIndex - startLineBreakIndex - 1
		newBuf := make([]byte, bufferSize)
		_, err = file.ReadAt(newBuf, curBufStart+startLineBreakIndex+1)
		if err == nil {
			lastLines = append(lastLines, string(newBuf))
		} else {
			return lastLines, err
		}
	}

	return lastLines, nil
}

// ReadLastLinesWithOffset reads the last initBufSize bytes in front of the fileOffset bytes before EOF
// fileOffset needs to be at a line break. otherwise the incomplete line at the end of the buffer will be lost
// returns the complete lines in reverse order, a new offset for the next call and an error if any
// Note: the initBufSize needs to be longer than the maximum length of a log line
func ReadLastLinesWithOffset(fileName string, fileOffset int64, initBufSize int) ([]string, int64, error) {
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
		return []string{}, fileSize, nil
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

	lastLines := []string{}
	for i := 0; i < len(indices)-1; i++ {
		// create a buffer the size between two line breaks (excluding the second line break)
		// and then read everything in between
		endLineBreakIndex := indices[len(indices)-1-i]
		startLineBreakIndex := indices[len(indices)-2-i]
		bufferSize := endLineBreakIndex - startLineBreakIndex - 1

		newBuf := make([]byte, bufferSize)
		_, err = file.ReadAt(newBuf, curBufStart+startLineBreakIndex+1)
		if err == nil {
			lastLines = append(lastLines, string(newBuf))
		} else {
			return lastLines, 0, err
		}
	}

	// if this is the beginning of the file, append the first line
	// (which doesn't start with a line break :facepalm)
	if curBufStart == 0 {
		lastLines = append(lastLines, string(buf[:indices[0]]))
		// set fileOffset to be file size, indication we're scanned through
		return lastLines, fileSize, nil
	}

	return lastLines, int64(bufSize) - indices[0] - 1 + fileOffset, nil
}

func ReadLastNLinesWithKeyword(fileName string, n int, query string) ([]string, error) {
	// buffer size 32KB
	// Note: if log lines are longer than 32KB, the program won't work
	initBufSize := 1 << 15

	lines := []string{}
	var fileOffset int64 = 0
	newlines := []string{"dummy"}
	var err error
	for len(lines) < n && len(newlines) != 0 {
		newlines, fileOffset, err = ReadLastLinesWithOffset(fileName, fileOffset, initBufSize)
		if err != nil {
			panic(err)
		}

		if query == "" {
			lines = append(lines, newlines...)
		} else {
			for _, newline := range newlines {
				if strings.Contains(newline, query) {
					lines = append(lines, newline)
				}
			}
		}
	}

	if len(lines) >= n {
		return lines[:n], nil
	}

	return lines, nil
}
