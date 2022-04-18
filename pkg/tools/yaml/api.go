package yaml

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
)

func RemoveComments(in io.Reader) io.Reader {
	scanner := bufio.NewScanner(in)
	tester := regexp.MustCompile(`^\s*#`)
	result := bytes.Buffer{}

	for scanner.Scan() {
		line := scanner.Bytes()

		if !tester.Match(line) {
			result.Write(line)
			result.Write([]byte("\n"))
		}
	}

	return &result
}
