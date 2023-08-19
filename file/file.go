package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

func WriteFileHelper(fileName string, numOfLines int) {
	// open output file
	fo, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
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

func WriteFileHelper2(fileName string, numOfLines int) {
	fileByte := []byte{}
	for ; numOfLines > 0; numOfLines-- {
		newLine := fmt.Sprintf("%s\n", time.Now())
		fileByte = append(fileByte, []byte(newLine)...)
	}

	os.WriteFile(fileName, fileByte, os.ModePerm)
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
