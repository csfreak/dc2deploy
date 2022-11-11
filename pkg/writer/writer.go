/*
Copyright © 2022 Jason Ross

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package writer

import (
	"fmt"
	"os"
	"strings"
)

var Level uint8 = 0

func WriteOut(level uint8, fstring string, args ...interface{}) {
	if !strings.HasSuffix(fstring, "\n") {
		fstring += "\n"
	}

	switch {
	case level == 0:
		fmt.Fprintf(os.Stdout, fstring, args...)
	case level <= Level:
		fmt.Fprintf(os.Stderr, fstring, args...)
	}
}

func WriteErr(level uint8, fstring string, args ...interface{}) {
	if !strings.HasSuffix(fstring, "\n") {
		fstring += "\n"
	}

	if level <= Level {
		fmt.Fprintf(os.Stderr, fstring, args...)
	}
}

func WriteFile(path string, data []byte) error {
	if path == "-" {
		WriteOut(0, string(data))
		return nil
	}

	return os.WriteFile(path, data, 0644)
}
