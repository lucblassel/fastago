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
	"github.com/lucblassel/fastago/pkg/seqs"
	"strings"

	"github.com/spf13/cobra"
)

// upperCmd represents the upper command
var upperCmd = &cobra.Command{
	Use:   "upper",
	Short: "Uppercase all sequence nucleotides",
	RunE: func(cmd *cobra.Command, args []string) error {

		records := make(chan seqs.SeqRecord)
		errs := make(chan error)

		go seqs.ReadFastaRecords(inputReader, records, errs)

		for records != nil && errs != nil {
			select {
			case record := <-records:
				output, err := record.Seq.FormatSeq(outputLineWidth)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintf(outputWriter, ">%s\n%s\n", record.Name, strings.ToUpper(output))
				if err != nil {
					return err
				}
			case err := <-errs:
				return err
			}
		}

		return nil
	},
}

// init adds the command to the root
func init() {
	transformCmd.AddCommand(upperCmd)
}
