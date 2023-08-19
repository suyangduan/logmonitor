package file

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
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

// TODO: no lb within the first scan
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
