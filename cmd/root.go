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
package cmd

import (
	"fmt"
	"os"

	"github.com/csfreak/dc2deploy/pkg/command"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/validation"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dc2deploy [name]",
	Short: "Convert Openshift DeploymentConfig to Kuberentes Deployment",
	Long:  `Convert Openshift DeploymentConfig to Kuberentes Deployment. It can source from and output to json, yaml, or kubernetes. Flags and Args match kubectl where possible.`,
	Example: `
From File:
dc2deploy -f dc.yaml --output deploy.yaml
	
From Kubernetes:
dc2deploy dcname -n namespacename --dry-run`,
	Args:    validateArgs,
	PreRunE: validateFlags,
	RunE:    command.RunE,
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("filename", "f", "-", "File containing DeploymentConfig manifest")
	rootCmd.MarkFlagFilename("filename")

	// Live Flags
	rootCmd.Flags().Bool("dry-run", false, "Only print the new object that would be sent")
	rootCmd.Flags().StringP("namespace", "n", "", "Namespace of DeploymentConfig")
	rootCmd.Flags().String("kubeconfig", "", "Path to Kubeconfig")
	rootCmd.MarkFlagsMutuallyExclusive("kubeconfig", "filename")

	// Output Flags
	rootCmd.Flags().String("outfile", "-", "Output filename. Defaults to STDOUT")
	rootCmd.Flags().StringP("output", "o", "yaml", "Output in JSON")

	// Options
	rootCmd.Flags().Bool("ignore-warnings", false, "Ignore Warnings about missing Deployment Features")

	rootCmd.Flags().UintP("verbosity", "v", 0, "Set Verbosity")
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
		return err
	}

	if len(args) == 1 {
		if errs := validation.NameIsDNSSubdomain(args[0], false); errs != nil {
			return fmt.Errorf("invalid deploymentconfig name: %s", errs)
		}
	}

	return nil
}

func validateFlags(cmd *cobra.Command, args []string) error {
	c := &command.CommandOptions{}

	if filename, err := cmd.Flags().GetString("filename"); err == nil {
		c.Filename = filename
	}

	if outfile, err := cmd.Flags().GetString("outfile"); err == nil {
		c.OutputFilename = outfile
	}

	if output, err := cmd.Flags().GetString("output"); err == nil {
		c.OutputFileType = command.FileType(output)
	}

	if namespace, err := cmd.Flags().GetString("namespace"); err == nil {
		c.LiveNamespace = namespace
	}

	if kubeconfig, err := cmd.Flags().GetString("kubeconfig"); err == nil {
		c.LiveKubeconfig = kubeconfig
	}

	if len(args) == 1 {
		c.LiveDC = args[0]
	}

	if dryrun, err := cmd.Flags().GetBool("dry-run"); err == nil {
		c.LiveDryRun = dryrun
	}

	if ignore, err := cmd.Flags().GetBool("ignore-warnings"); err == nil {
		c.IgnoreWarnings = ignore
	}

	if verbosity, err := cmd.Flags().GetUint8("verbosity"); err == nil {
		c.Verbosity = verbosity
	}

	return command.SetCommandOptions(c)
}
