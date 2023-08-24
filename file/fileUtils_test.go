package file_test

import (
	"cribl/logmonitor/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("write file helper", func() {
	It("work as expected", func() {
		fileName := "test.txt"
		Expect(true).To(BeTrue())
		file.WriteFileHelper(fileName, 100)

		// open input file
		fi, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		// close fi on exit and check for its returned error
		defer func() {
			if err := fi.Close(); err != nil {
				panic(err)
			}
			os.Remove(fileName)
		}()

		numOfLines, err := file.LineCounter(fi)

		Expect(err).To(BeNil())
		Expect(numOfLines).To(Equal(100))
	})
})

//var _ = Describe("This test generates a test file", func() {
//	It("", func() {
//		file.GenerateTestFile("1KB.txt", 100)
//	})
//})
