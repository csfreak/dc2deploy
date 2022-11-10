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
	"github.com/csfreak/dc2deploy/pkg/k8s"
	klog "k8s.io/klog/v2"
)

func DoLive() error {
	err := k8s.Init(Options.LiveKubeconfig)
	if err != nil {
		return fmt.Errorf("unable to create kubernetes client: %w", err)
	}

	dc, err := k8s.LoadDC(Options.LiveDC, Options.LiveNamespace)
	if err != nil {
		return fmt.Errorf("unable to create load %s: %w", Options.LiveDC, err)
	}

	warnings := convert.CheckFeatures(dc)
	if warnings != nil {
		for _, w := range warnings {
			if Options.LiveDryRun && Options.Verbosity < 2 {
				w.Log(2)
			} else {
				w.Log(Options.Verbosity)
			}
		}

		if Options.IgnoreWarnings {
			klog.V(2).InfoS("ignoring warnings")
		} else {
			return fmt.Errorf("use --ignore-warnings to continue")
		}
	}

	deploy, err := convert.ToDeploy(dc)
	if err != nil {
		return fmt.Errorf("unable to convert to deploy: %w", err)
	}

	if Options.LiveDryRun {
		o, err := convert.ToOuput(deploy, string(Options.OutputFileType))
		if err != nil {
			return fmt.Errorf("unable to marshal object: %w", err)
		}

		fmt.Print("-----------\n\n")
		fmt.Println(string(o))
		fmt.Println("-----------")

		return nil
	}

	return nil
}
