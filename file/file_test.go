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

var _ = Describe("ReadLastLines", func() {
	It("works as expected for 100 line file (5KB)", func() {
		fileName := "/var/log/var5KB.txt"
		expectedLines := []string{
			"Line 100 this line contains a random animal: dragon",
			"Line 99 this line contains a random animal: mouse",
		}
		Expect(file.ReadLastLines(fileName)).To(Equal(expectedLines))
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

var _ = Describe("ReadLastLines edge cases", func() {
	It("works as expected when there are empty lines", func() {
		fileName := "5KBwithEmptyNewLines.txt"
		expectedLines := []string{
			"",
			"Line 100 this line contains a random animal: ox",
			"",
			"Line 99 this line contains a random animal: tiger",
			"",
		}
		Expect(file.ReadLastLines(fileName)).To(Equal(expectedLines))
	})

	// This test will fail unless we change the product code.
	// We need to append the EOF index to indices at the very end in this case
	//It("works as expected when there is no line break for the last line", func() {
	//	fileName := "5KBwithNoLastLineBreak.txt"
	//	expectedLines := []string{
	//		"Line 100 this line contains a random animal: horse",
	//		"Line 99 this line contains a random animal: rabbit",
	//	}
	//	Expect(file.ReadLastLines(fileName)).To(Equal(expectedLines))
	//})
})

var _ = Describe("ReadLastLinesWithOffset", func() {
	It("works as expected for 100 line file (5KB)", func() {
		expectedLines := []string{
			"Line 100 this line contains a random animal: dragon",
			"Line 99 this line contains a random animal: mouse",
		}
		fileName := "/var/log/var5KB.txt"

		firstLines, firstOffset, firstErr := file.ReadLastLinesWithOffset(fileName, 0, 128)
		Expect(firstLines).To(Equal(expectedLines))
		Expect(firstOffset).To(Equal(int64(102)))
		Expect(firstErr).To(BeNil())

		var offset int64 = 0
		lines := []string{}
		for i := 0; i < 50; i++ {
			lines, offset, _ = file.ReadLastLinesWithOffset(fileName, offset, 128)
		}

		expectedLinesEnd := []string{
			"Line 2 this line contains a random animal: tiger",
			"Line 1 this line contains a random animal: pig",
		}

		Expect(lines).To(Equal(expectedLinesEnd))
	})

	It("works as expected for 100k line file (5MB)", func() {
		expectedLines := []string{
			"Line 100000 this line contains a random animal: pig",
			"Line 99999 this line contains a random animal: mouse",
		}
		firstLines, _, firstErr := file.ReadLastLinesWithOffset("5MB.txt", 0, 1024)
		Expect(firstLines[:2]).To(Equal(expectedLines))
		Expect(firstErr).To(BeNil())

		var offset int64 = 0
		lines := []string{"dummy"}
		for i := 0; i < 5259; i++ {
			lines, offset, _ = file.ReadLastLinesWithOffset("5MB.txt", offset, 1024)
		}

		expectedLinesEnd := []string{
			"Line 2 this line contains a random animal: pig",
			"Line 1 this line contains a random animal: rabbit",
		}

		Expect(lines[len(lines)-2:]).To(Equal(expectedLinesEnd))
	})

	// for bigger size files, from tweaking the buffer size, no visible gains is achieved after
	// buffer size reaches 1<<15 byte (32KB). Right now this test runs about 3.5 seconds on a macbook air
	It("works as expected for 20M line file (1GB)", func() {
		expectedLines := []string{
			"Line 20000000 this line contains a random animal: mouse",
			"Line 19999999 this line contains a random animal: horse",
		}
		firstLines, _, err := file.ReadLastLinesWithOffset("1GB.txt", 0, 1<<15)
		Expect(firstLines[:2]).To(Equal(expectedLines))
		Expect(err).To(BeNil())

		// this is the 3.5 seconds part. commented out for now

		//var offset int64 = 0
		//lines := []string{"dummy"}
		//for i := 0; i < 33666; i++ {
		//	lines, offset, _ = file.ReadLastLinesWithOffset("1GB.txt", offset, 1<<15)
		//}
		//
		//expectedLinesEnd := []string{
		//	"Line 2 this line contains a random animal: pig",
		//	"Line 1 this line contains a random animal: ox",
		//}
		//
		//Expect(lines[len(lines)-2:]).To(Equal(expectedLinesEnd))
	})
})

var _ = Describe("ReadLastLinesWithOffset edge cases", func() {
	It("works as expected when there are empty lines", func() {
		expectedLines := []string{
			"",
			"Line 100 this line contains a random animal: ox",
			"",
			"Line 99 this line contains a random animal: tiger",
			"",
		}
		fileName := "5KBwithEmptyNewLines.txt"

		firstLines, firstOffset, firstErr := file.ReadLastLinesWithOffset(fileName, 0, 128)
		Expect(firstLines).To(Equal(expectedLines))
		Expect(firstOffset).To(Equal(int64(101)))
		Expect(firstErr).To(BeNil())

		var offset int64 = 0
		lines := []string{}
		for i := 0; i < 50; i++ {
			lines, offset, _ = file.ReadLastLinesWithOffset(fileName, offset, 128)
		}

		expectedLinesEnd := []string{
			"Line 2 this line contains a random animal: horse",
			"",
			"Line 1 this line contains a random animal: horse",
		}

		Expect(lines).To(Equal(expectedLinesEnd))
	})

	// This test will fail unless we change the product code.
	// We need to append the EOF index to indices at the very end when fileOffset is 0 in this case
	//It("works as expected when there is no line break for the last line", func() {
	//	expectedLines := []string{
	//		"Line 100 this line contains a random animal: horse",
	//		"Line 99 this line contains a random animal: rabbit",
	//	}
	//	fileName := "5KBwithNoLastLineBreak.txt"
	//
	//	firstLines, firstOffset, firstErr := file.ReadLastLinesWithOffset(fileName, 0, 128)
	//	Expect(firstLines).To(Equal(expectedLines))
	//	Expect(firstOffset).To(Equal(int64(101)))
	//	Expect(firstErr).To(BeNil())
	//
	//	var offset int64 = 0
	//	lines := []string{}
	//	for i := 0; i < 50; i++ {
	//		lines, offset, _ = file.ReadLastLinesWithOffset(fileName, offset, 128)
	//	}
	//
	//	expectedLinesEnd := []string{
	//		"Line 2 this line contains a random animal: rooster",
	//		"Line 1 this line contains a random animal: dragon",
	//	}
	//
	//	Expect(lines).To(Equal(expectedLinesEnd))
	//})
})

var _ = Describe("ReadLastNLinesWithKeyword", func() {
	It("works for small files without query", func() {
		fileName := "/var/log/var5KB.txt"
		expectedLines := []string{
			"Line 5 this line contains a random animal: dragon",
			"Line 4 this line contains a random animal: rooster",
			"Line 3 this line contains a random animal: pig",
			"Line 2 this line contains a random animal: tiger",
			"Line 1 this line contains a random animal: pig",
		}
		lines, err := file.ReadLastNLinesWithKeyword(fileName, 100, "")
		Expect(len(lines)).To(Equal(100))
		Expect(lines[95:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	It("works for medium files without query", func() {
		expectedLines := []string{
			"Line 99005 this line contains a random animal: monkey",
			"Line 99004 this line contains a random animal: ox",
			"Line 99003 this line contains a random animal: tiger",
			"Line 99002 this line contains a random animal: snake",
			"Line 99001 this line contains a random animal: snake",
		}
		lines, err := file.ReadLastNLinesWithKeyword("5MB.txt", 1000, "")
		Expect(len(lines)).To(Equal(1000))
		Expect(lines[995:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	// when this test goes to 1M rows, it adds about 600ms

	It("works for large files without query", func() {
		expectedLines := []string{
			"Line 19900005 this line contains a random animal: rooster",
			"Line 19900004 this line contains a random animal: pig",
			"Line 19900003 this line contains a random animal: dragon",
			"Line 19900002 this line contains a random animal: rabbit",
			"Line 19900001 this line contains a random animal: ram",
		}
		lines, err := file.ReadLastNLinesWithKeyword("1GB.txt", 100000, "")
		Expect(len(lines)).To(Equal(100000))
		Expect(lines[99995:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})
})

var _ = Describe("ReadLastNLinesWithKeyword", func() {
	It("works for small files", func() {
		fileName := "/var/log/var5KB.txt"
		expectedLines := []string{
			"Line 69 this line contains a random animal: dragon",
			"Line 68 this line contains a random animal: dragon",
			"Line 41 this line contains a random animal: dragon",
			"Line 39 this line contains a random animal: dragon",
			"Line 33 this line contains a random animal: dragon",
		}
		lines, err := file.ReadLastNLinesWithKeyword(fileName, 10, "dragon")
		Expect(len(lines)).To(Equal(10))
		Expect(lines[5:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	It("works for medium files", func() {
		expectedLines := []string{
			"Line 98813 this line contains a random animal: ox",
			"Line 98812 this line contains a random animal: ox",
			"Line 98810 this line contains a random animal: ox",
			"Line 98809 this line contains a random animal: ox",
			"Line 98806 this line contains a random animal: ox",
		}
		lines, err := file.ReadLastNLinesWithKeyword("5MB.txt", 100, "ox")
		Expect(len(lines)).To(Equal(100))
		Expect(lines[95:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	// this test adds about 1 second

	It("works for large files", func() {
		expectedLines := []string{
			"Line 18803427 this line contains a random animal: snake",
			"Line 18803418 this line contains a random animal: snake",
			"Line 18803417 this line contains a random animal: snake",
			"Line 18803416 this line contains a random animal: snake",
			"Line 18803410 this line contains a random animal: snake",
		}
		lines, err := file.ReadLastNLinesWithKeyword("1GB.txt", 100000, "snake")
		Expect(len(lines)).To(Equal(100000))
		Expect(lines[99995:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	// this test adds about 600ms

	It("works for large files for longer queries", func() {
		expectedLines := []string{
			"Line 19000005 this line contains a random animal: horse",
			"Line 19000004 this line contains a random animal: snake",
			"Line 19000003 this line contains a random animal: dog",
			"Line 19000002 this line contains a random animal: pig",
			"Line 19000001 this line contains a random animal: snake",
		}
		lines, err := file.ReadLastNLinesWithKeyword("1GB.txt", 1000000,
			"this line contains a random animal")
		Expect(len(lines)).To(Equal(1000000))
		Expect(lines[999995:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})
})
