package file_test

import (
	"cribl/logmonitor/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadLastLinesWithOffsetPagination", func() {
	It("works as expected for 100 line file (5KB)", func() {
		expectedLines := []file.LineReturn{
			{
				"Line 100 this line contains a random animal: dragon",
				52,
			},
			{
				"Line 99 this line contains a random animal: mouse",
				102,
			},
		}
		fileName := "/var/log/var5KB.txt"

		firstLines, firstErr := file.ReadLastLinesWithOffsetPagination(fileName, 0, 128)
		Expect(firstLines).To(Equal(expectedLines))
		Expect(firstErr).To(BeNil())

		lines := []file.LineReturn{{"", 0}}
		for i := 0; i < 50; i++ {
			lines, _ = file.ReadLastLinesWithOffsetPagination(fileName, lines[len(lines)-1].Offset, 128)
		}

		expectedLinesEnd := []file.LineReturn{
			{
				"Line 2 this line contains a random animal: tiger",
				4933,
			},
			{
				"Line 1 this line contains a random animal: pig",
				4980,
			},
		}

		Expect(lines).To(Equal(expectedLinesEnd))
	})

})

var _ = Describe("ReadLastNLinesWithKeywordPagination", func() {
	It("works for small files without query", func() {
		fileName := "/var/log/var5KB.txt"
		expectedLines := []string{
			"Line 2 this line contains a random animal: tiger",
			"Line 1 this line contains a random animal: pig",
		}

		newlines := []string{}
		var offset int64 = 0
		var err error
		for i := 0; i < 50; i++ {
			newlines, offset, err = file.ReadLastNLinesWithKeywordPagination(fileName, 2, "", offset)
		}

		Expect(newlines).To(Equal(expectedLines))
		Expect(offset).To(Equal(int64(4980)))
		Expect(err).To(BeNil())
	})

	It("works for medium files without query", func() {
		fileName := "5MB.txt"
		expectedLines := []string{
			"Line 2 this line contains a random animal: pig",
			"Line 1 this line contains a random animal: rabbit",
		}

		newlines := []string{}
		var offset int64 = 0
		var err error
		for i := 0; i < 50; i++ {
			newlines, offset, err = file.ReadLastNLinesWithKeywordPagination(fileName, 2000, "", offset)
		}

		Expect(newlines[len(newlines)-2:]).To(Equal(expectedLines))
		Expect(offset).To(Equal(int64(5256265)))
		Expect(err).To(BeNil())
	})

	It("works for medium files with query", func() {
		fileName := "5MB.txt"
		expectedLines := []string{
			"Line 3 this line contains a random animal: rabbit",
			"Line 1 this line contains a random animal: rabbit",
		}

		newlines := []string{""}
		var offset int64 = 0
		var err error
		for i := 0; i < 5; i++ {
			newlines, offset, err = file.ReadLastNLinesWithKeywordPagination(fileName, 2000, "rabbit", offset)
		}

		Expect(newlines[len(newlines)-2:]).To(Equal(expectedLines))
		Expect(offset).To(Equal(int64(5256265)))
		Expect(err).To(BeNil())
	})
})

var _ = Describe("ReadLastNLinesWithKeywordPaginationInternal", func() {
	It("works for small files and a small buffer size without query", func() {
		fileName := "/var/log/var5KB.txt"
		expectedLines := []string{
			"Line 2 this line contains a random animal: tiger",
			"Line 1 this line contains a random animal: pig",
		}

		newlines := []string{}
		var offset int64 = 0
		var err error
		for i := 0; i < 50; i++ {
			newlines, offset, err = file.ReadLastNLinesWithKeywordPaginationInternal(fileName, 2, "", offset, 60)
		}

		Expect(newlines).To(Equal(expectedLines))
		Expect(offset).To(Equal(int64(4980)))
		Expect(err).To(BeNil())
	})

	It("works for small files and a small buffer size with query", func() {
		fileName := "/var/log/var5KB.txt"
		expectedLines := []string{
			"Line 3 this line contains a random animal: pig",
			"Line 1 this line contains a random animal: pig",
		}

		newlines := []string{"dummy"}
		var offset int64 = 0
		var err error
		for i := 0; i < 4; i++ {
			newlines, offset, err = file.ReadLastNLinesWithKeywordPaginationInternal(fileName, 2, "pig", offset, 60)
		}

		Expect(newlines).To(Equal(expectedLines))
		Expect(offset).To(Equal(int64(4980)))
		Expect(err).To(BeNil())
	})
})
