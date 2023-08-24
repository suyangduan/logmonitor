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
