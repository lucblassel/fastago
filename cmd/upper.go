/*
Copyright © 2021 LUC BLASSEL

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/lucblassel/fastago/pkg/stream"
	"github.com/spf13/cobra"
)

// upperCmd represents the upper command
var upperCmd = &cobra.Command{
	Use:   "upper",
	Short: "Uppercase all sequence nucleotides",
	RunE: func(cmd *cobra.Command, args []string) error {

		lines := make(chan stream.Line)
		errs := make(chan error)

		go stream.StreamFasta(inputReader, lines, errs)
		var err error

		select {
		case line := <-lines:
			if line.IsName {
				_, err = fmt.Fprintln(outputWriter, line.Line)
			} else {
				_, err = fmt.Fprintln(outputWriter, strings.ToUpper(line.Line))
			}
			if err != nil {
				return err
			}
		case err := <-errs:
			return err
		}
		return nil
	},
}

func init() {
	transformCmd.AddCommand(upperCmd)
}
