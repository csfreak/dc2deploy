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

	"github.com/csfreak/dc2deploy/pkg/writer"
)

var Options *CommandOptions

type CommandOptions struct {
	inputType      IOType   `default:"FileIOType"`
	outputType     IOType   `default:"FileIOType"`
	Filename       string   `default:""`
	OutputFilename string   `default:""`
	OutputFileType FileType `default:"YAMLFileType"`
	LiveDryRun     bool     `default:"false"`
	LiveNamespace  string   `default:""`
	LiveDC         string   `default:""`
	LiveKubeconfig string   `default:""`
	IgnoreWarnings bool     `default:"false"`
	Verbosity      uint8    `default:"0"`
}

type IOType string
type FileType string

const (
	FileIOType   IOType   = "file"
	LiveIOType   IOType   = "live"
	JSONFileType FileType = "json"
	YAMLFileType FileType = "yaml"
)

func SetCommandOptions(c *CommandOptions) error {
	if Options == nil {
		Options = &CommandOptions{}
	}

	if c.LiveDC != "" {
		Options.LiveDC = c.LiveDC
		Options.LiveNamespace = c.LiveNamespace
		Options.LiveKubeconfig = c.LiveKubeconfig
		Options.LiveDryRun = c.LiveDryRun
		Options.inputType = LiveIOType

		if c.LiveDryRun {
			Options.outputType = FileIOType
			Options.OutputFilename = "-"
		} else {
			Options.outputType = LiveIOType
		}

		if c.Filename != "-" || c.OutputFilename != "-" {
			return fmt.Errorf("cannot specify filename or outfile on live operation")
		}
	} else {
		Options.Filename = c.Filename
		Options.inputType = FileIOType
		Options.outputType = FileIOType

		if c.LiveDryRun ||
			c.LiveKubeconfig != "" ||
			c.LiveNamespace != "" ||
			c.LiveDC != "" {
			return fmt.Errorf("cannot specify input filename and live options")
		}

		if c.OutputFilename != "" {
			Options.OutputFilename = c.OutputFilename
		}
	}

	if Options.outputType == FileIOType {
		Options.OutputFileType = c.OutputFileType
	}

	Options.IgnoreWarnings = c.IgnoreWarnings
	Options.Verbosity = c.Verbosity

	if Options.Verbosity > 4 {
		Options.Verbosity = 4
	}

	writer.Level = Options.Verbosity

	return nil
}
