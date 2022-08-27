package xray

import (
	"bufio"
	"fmt"
	"os"
)

// LineReader will accept the name of a file and offset as argument
// and will return a channel from which lines can be read
// one at a time.
func LineReader(filename string) (chan string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	out := make(chan string)
	go func() {
		defer func() {
			if err = fp.Close(); err != nil {
				fmt.Printf("error closing %s: %v\n", filename, err)
			}
		}()

		// we need to close the out channel in order
		// to signal the end-of-data condition
		defer close(out)

		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}
