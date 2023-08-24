package file_test

import (
	"cribl/logmonitor/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadLastLinesWithOffsetP", func() {
	It("works as expected for 100 line file (5KB)", func() {
		expectedLines := [][]byte{
			[]byte("Line 100 this line contains a random animal: dragon\n"),
			[]byte("Line 99 this line contains a random animal: mouse\n"),
			[]byte("ns a random animal: horse\n"),
		}
		fileName := "5KB.txt"

		firstLines, firstErr := file.ReadLastLinesWithOffsetP(fileName, 0, 128)

		Expect(firstLines).To(Equal(expectedLines))
		Expect(firstErr).To(BeNil())

		lines, _ := file.ReadLastLinesWithOffsetP(fileName, 128*38, 128)

		expectedLinesEnd := [][]byte{
			[]byte("Line 3 this line con"),
			[]byte("Line 2 this line contains a random animal: tiger\n"),
			[]byte("Line 1 this line contains a random animal: pig\n"),
		}

		Expect(lines).To(Equal(expectedLinesEnd))
	})

	It("works as expected when trying to read 0 byte", func() {
		fileName := "5KB.txt"
		firstLines, firstErr := file.ReadLastLinesWithOffsetP(fileName, 0, 0)

		Expect(firstLines).To(Equal([][]byte{}))
		Expect(firstErr).To(BeNil())
	})

	It("works as expected for 100k line file (5MB)", func() {
		fileName := "5MB.txt"
		expectedLines := [][]byte{
			[]byte("Line 100000 this line contains a random animal: pig\n"),
			[]byte("Line 99999 this line contains a random animal: mouse\n"),
		}
		firstLines, firstErr := file.ReadLastLinesWithOffsetP(fileName, 0, 1024)
		Expect(firstLines[:2]).To(Equal(expectedLines))
		Expect(firstErr).To(BeNil())

		lines, _ := file.ReadLastLinesWithOffsetP(fileName, 128*41064, 128)
		expectedLinesEnd := [][]byte{
			[]byte("Line 2 this line contai"),
			[]byte("Line 1 this line contains a random animal: rabbit\n"),
		}

		Expect(lines).To(Equal(expectedLinesEnd))
	})

	It("works as expected for 20M line file (1GB)", func() {
		expectedLines := [][]byte{
			[]byte("Line 20000000 this line contains a random animal: mouse\n"),
			[]byte("Line 19999999 this line contains a random animal: horse\n"),
		}
		firstLines, err := file.ReadLastLinesWithOffsetP("1GB.txt", 0, 1<<15)
		Expect(firstLines[:2]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})
})

var _ = Describe("CombineLines", func() {
	It("works as expected when joining segmented lines", func() {
		fileName := "5KB.txt"
		firstLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 128)
		secondLines, _ := file.ReadLastLinesWithOffsetP(fileName, 128, 128)
		combinedLines := file.CombineLines(firstLines, secondLines)

		comboLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 256)
		Expect(comboLines).To(Equal(combinedLines))
	})

	It("works as expected when first input is empty", func() {
		fileName := "5KB.txt"
		firstLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 0)
		secondLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 128)
		combinedLines := file.CombineLines(firstLines, secondLines)

		comboLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 128)
		Expect(comboLines).To(Equal(combinedLines))
	})

	It("works as expected when second input is empty", func() {
		fileName := "5KB.txt"
		firstLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 128)
		secondLines, _ := file.ReadLastLinesWithOffsetP(fileName, 128, 0)
		combinedLines := file.CombineLines(firstLines, secondLines)

		comboLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 128)
		Expect(comboLines).To(Equal(combinedLines))
	})

	// first lines are
	// "full line\n"
	// "full line\n"
	// second lines are
	// "full line\n"
	// "full line\n"
	// "segmented line\n"
	It("works as expected when joining full lines", func() {
		fileName := "5KB.txt"
		firstLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 102)
		secondLines, _ := file.ReadLastLinesWithOffsetP(fileName, 102, 128)
		Expect(secondLines[0][len(secondLines[0])-1]).To(Equal(byte('\n')))

		combinedLines := file.CombineLines(firstLines, secondLines)

		comboLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 230)
		Expect(comboLines).To(Equal(combinedLines))
	})

	// first lines are
	// "full line\n"
	// "full line\n"
	// "\n"
	// second lines are
	// "full line" (no line break)
	// "full line\n"
	// "segmented line\n"
	It("works as expected when joining full lines", func() {
		fileName := "5KB.txt"
		firstLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 103)
		secondLines, _ := file.ReadLastLinesWithOffsetP(fileName, 103, 128)
		Expect(firstLines[2]).To(Equal([]byte{'\n'}))

		combinedLines := file.CombineLines(firstLines, secondLines)

		comboLines, _ := file.ReadLastLinesWithOffsetP(fileName, 0, 231)
		Expect(comboLines).To(Equal(combinedLines))
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
})

var _ = Describe("ReadLastNLinesWithKeywordP", func() {
	It("works for small files without query", func() {
		fileName := "5KB.txt"
		expectedLines := [][]byte{
			[]byte("Line 5 this line contains a random animal: dragon\n"),
			[]byte("Line 4 this line contains a random animal: rooster\n"),
			[]byte("Line 3 this line contains a random animal: pig\n"),
			[]byte("Line 2 this line contains a random animal: tiger\n"),
			[]byte("Line 1 this line contains a random animal: pig\n"),
		}
		lines, err := file.ReadLastNLinesWithKeywordP(fileName, 100, "")
		Expect(len(lines)).To(Equal(100))
		Expect(lines[95:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	It("works for medium files without query", func() {
		expectedLines := [][]byte{
			[]byte("Line 99005 this line contains a random animal: monkey\n"),
			[]byte("Line 99004 this line contains a random animal: ox\n"),
			[]byte("Line 99003 this line contains a random animal: tiger\n"),
			[]byte("Line 99002 this line contains a random animal: snake\n"),
			[]byte("Line 99001 this line contains a random animal: snake\n"),
		}
		lines, err := file.ReadLastNLinesWithKeywordP("5MB.txt", 1000, "")
		Expect(len(lines)).To(Equal(1000))
		Expect(lines[995:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	It("works for large files without query", func() {
		expectedLines := [][]byte{
			[]byte("Line 19900005 this line contains a random animal: rooster\n"),
			[]byte("Line 19900004 this line contains a random animal: pig\n"),
			[]byte("Line 19900003 this line contains a random animal: dragon\n"),
			[]byte("Line 19900002 this line contains a random animal: rabbit\n"),
			[]byte("Line 19900001 this line contains a random animal: ram\n"),
		}
		lines, err := file.ReadLastNLinesWithKeywordP("1GB.txt", 1000000, "")
		Expect(len(lines)).To(Equal(1000000))
		Expect(lines[99995:100000]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})
})

var _ = Describe("ReadLastNLinesWithKeywordP", func() {
	It("works for small files", func() {
		fileName := "/var/log/var5KB.txt"
		expectedLines := [][]byte{
			[]byte("Line 69 this line contains a random animal: dragon\n"),
			[]byte("Line 68 this line contains a random animal: dragon\n"),
			[]byte("Line 41 this line contains a random animal: dragon\n"),
			[]byte("Line 39 this line contains a random animal: dragon\n"),
			[]byte("Line 33 this line contains a random animal: dragon\n"),
		}
		lines, err := file.ReadLastNLinesWithKeywordP(fileName, 10, "dragon")
		Expect(len(lines)).To(Equal(10))
		Expect(lines[5:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	It("works for small files when we reach the end of the file", func() {
		fileName := "5KB.txt"
		lines, err := file.ReadLastNLinesWithKeywordP(fileName, 100, "snake")
		Expect(lines[len(lines)-1]).To(Equal([]byte("Line 11 this line contains a random animal: snake\n")))
		Expect(err).To(BeNil())
	})

	It("works for small files when we reach the end of the file and the rolling last line contains keyword", func() {
		fileName := "5KB.txt"
		lines, err := file.ReadLastNLinesWithKeywordP(fileName, 100, "pig")
		Expect(lines[len(lines)-1]).To(Equal([]byte("Line 1 this line contains a random animal: pig\n")))
		Expect(err).To(BeNil())
	})

	It("works for medium files", func() {
		expectedLines := [][]byte{
			[]byte("Line 98813 this line contains a random animal: ox\n"),
			[]byte("Line 98812 this line contains a random animal: ox\n"),
			[]byte("Line 98810 this line contains a random animal: ox\n"),
			[]byte("Line 98809 this line contains a random animal: ox\n"),
			[]byte("Line 98806 this line contains a random animal: ox\n"),
		}
		lines, err := file.ReadLastNLinesWithKeywordP("5MB.txt", 100, "ox")
		Expect(len(lines)).To(Equal(100))
		Expect(lines[95:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	It("works for large files", func() {
		expectedLines := [][]byte{
			[]byte("Line 18803427 this line contains a random animal: snake\n"),
			[]byte("Line 18803418 this line contains a random animal: snake\n"),
			[]byte("Line 18803417 this line contains a random animal: snake\n"),
			[]byte("Line 18803416 this line contains a random animal: snake\n"),
			[]byte("Line 18803410 this line contains a random animal: snake\n"),
		}
		lines, err := file.ReadLastNLinesWithKeywordP("1GB.txt", 100000, "snake")
		Expect(len(lines)).To(Equal(100000))
		Expect(lines[99995:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})

	// this test scans through 20M lines and takes about 5 seconds
	//It("works when we read the end of the file for large files", func() {
	//	lines, err := file.ReadLastNLinesWithKeywordP("1GB.txt", 10000000, "ox")
	//	fmt.Println(string(lines[len(lines)-1]))
	//	Expect(lines[len(lines)-1]).To(Equal([]byte("Line 1 this line contains a random animal: ox\n")))
	//	Expect(err).To(BeNil())
	//})

	// This test adds about 2.6 seconds

	//It("works for large files for bigger number of return entires", func() {
	//	expectedLines := [][]byte{
	//		[]byte("Line 18803427 this line contains a random animal: snake\n"),
	//		[]byte("Line 18803418 this line contains a random animal: snake\n"),
	//		[]byte("Line 18803417 this line contains a random animal: snake\n"),
	//		[]byte("Line 18803416 this line contains a random animal: snake\n"),
	//		[]byte("Line 18803410 this line contains a random animal: snake\n"),
	//	}
	//	lines, err := file.ReadLastNLinesWithKeywordP("1GB.txt", 1000000, "snake")
	//	Expect(len(lines)).To(Equal(1000000))
	//	Expect(lines[99995:100000]).To(Equal(expectedLines))
	//	Expect(err).To(BeNil())
	//})

	It("works for large files for longer queries", func() {
		expectedLines := [][]byte{
			[]byte("Line 19000005 this line contains a random animal: horse\n"),
			[]byte("Line 19000004 this line contains a random animal: snake\n"),
			[]byte("Line 19000003 this line contains a random animal: dog\n"),
			[]byte("Line 19000002 this line contains a random animal: pig\n"),
			[]byte("Line 19000001 this line contains a random animal: snake\n"),
		}
		lines, err := file.ReadLastNLinesWithKeywordP("1GB.txt", 1000000,
			"this line contains a random animal")
		Expect(len(lines)).To(Equal(1000000))
		Expect(lines[999995:]).To(Equal(expectedLines))
		Expect(err).To(BeNil())
	})
})

var _ = Describe("RevertBufferByLineBreak", func() {
	It("works as expected with no line break at the end", func() {
		input := []byte("\n\nhello\nworld\nagain")
		output := file.RevertBufferByLineBreak(input)
		Expect(output).To(Equal([][]byte{
			[]byte("again"),
			[]byte("world\n"),
			[]byte("hello\n"),
			[]byte("\n"),
			[]byte("\n"),
		}))
	})

	It("works as expected with line breaks at the end", func() {
		input := []byte("\n\nhello\nworld\nagain\n\n")
		output := file.RevertBufferByLineBreak(input)
		Expect(output).To(Equal([][]byte{
			[]byte("\n"),
			[]byte("again\n"),
			[]byte("world\n"),
			[]byte("hello\n"),
			[]byte("\n"),
			[]byte("\n"),
		}))
	})
})
