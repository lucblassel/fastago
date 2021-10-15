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
	"github.com/spf13/cobra"
	"io"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "count the number of sequences in a fasta file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		count, err := countSeqs(inputReader)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(outputWriter, count)
		return err
	},
}

// init adds the command to the root
func init() {
	statsCmd.AddCommand(countCmd)
}

// countSeqs returns the number of seqRecord elements in the input stream
func countSeqs(input io.Reader) (int, error) {
	records := make(chan seqs.SeqRecord)
	errs := make(chan error)
	go seqs.ReadFastaRecords(input, records, errs)

	count := 0

	for records != nil && errs != nil {
		select {
		case <-records:
			count++
		case err := <-errs:
			if err != nil {
				return 0, err
			}
			return count, nil

		}
	}

	return count, nil
}