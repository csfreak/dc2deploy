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

package command

import (
	"fmt"

	"github.com/csfreak/dc2deploy/pkg/convert"
	"github.com/csfreak/dc2deploy/pkg/writer"
	"github.com/spf13/cobra"
)

func RunE(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	if Options.inputType == LiveIOType {
		return DoLive()
	}

	return DoConvert()
}

func DoConvert() error {
	dc, err := convert.LoadDC(Options.Filename)
	if err != nil {
		return fmt.Errorf("unable to load %s: %w", Options.Filename, err)
	}

	warnings := convert.CheckFeatures(dc)
	if warnings != nil {
		checkWarnings(warnings)

		if !Options.IgnoreWarnings {
			return fmt.Errorf("use --ignore-warnings to continue")
		}
	}

	deploy, err := convert.ToDeploy(dc)
	if err != nil {
		return fmt.Errorf("unable to convert to deploy: %w", err)
	}

	o, err := convert.ToOuput(deploy, string(Options.OutputFileType))
	if err != nil {
		return fmt.Errorf("unable to marshal object: %w", err)
	}

	return writer.WriteFile(Options.OutputFilename, o)
}

func checkWarnings(w []*convert.Warning) {
	var warningLogLevel uint8

	if Options.IgnoreWarnings {
		warningLogLevel = 2

		writer.WriteErr(2, "ignoring warnings")
	} else {
		warningLogLevel = 0
	}

	for _, warning := range w {
		warning.Print(warningLogLevel)
	}
}
