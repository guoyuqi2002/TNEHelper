package common

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// starts capturing stdout from the very beginning
func StdoutRedirect(onLine func(string)) {
	orig := os.Stdout
	// create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		panic(fmt.Errorf("failed to create pipe: %w", err))
	}

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			onLine(line)
		}

		if err := scanner.Err(); err != nil && err != io.EOF {
			msg := fmt.Sprintf("error while reading output: %v", err)
			fmt.Fprintln(orig, msg)
			onLine(msg)
		}
	}()
	os.Stdout = w
}
