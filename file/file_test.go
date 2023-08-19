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

var _ = Describe("Read last line", func() {
	It("works as expected for 100 line file (5KB)", func() {
		expectedLines := []string{
			"Line 100 this line contains a random animal: dragon",
			"Line 99 this line contains a random animal: mouse",
		}
		Expect(file.ReadLastLines("5KB.txt")).To(Equal(expectedLines))
	})

	It("works as expected for 100k line file (5MB)", func() {
		expectedLines := []string{
			"Line 100000 this line contains a random animal: pig",
			"Line 99999 this line contains a random animal: mouse",
		}
		Expect(file.ReadLastLines("5MB.txt")).To(Equal(expectedLines))
	})

	It("works as expected for 20M line file (1GB)", func() {
		expectedLines := []string{
			"Line 20000000 this line contains a random animal: mouse",
			"Line 19999999 this line contains a random animal: horse",
		}
		Expect(file.ReadLastLines("1GB.txt")).To(Equal(expectedLines))
	})
})
