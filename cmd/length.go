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
	"errors"
	"fmt"
	"io"

	"github.com/lucblassel/fastago/pkg/seqs"
	"github.com/spf13/cobra"
)

var lengthMode string

// lengthCmd represents the length command
var lengthCmd = &cobra.Command{
	Use:   "length",
	Short: "get length of sequences in fasta file",
	Long: `By default this command outputs the length of each sequence. 
	It is also possible to retreive the minimum, maximum or average length.
	set the -m/--mode flag to either:
		- each: will display the length of each sequence
		- average (or mean) will display the average lengh
		- min (or minimum) will display the minimum length
		- max (or maximum) will display the maximum length`,
	RunE: func(cmd *cobra.Command, args []string) error {

		records := make(chan seqs.SeqRecord)
		errs := make(chan error)

		go seqs.ReadFastaRecords(inputReader, records, errs)

		var err error

		switch lengthMode {
		case "each":
			err = getEach(records, errs, outputWriter)
		case "average":
			err = getAverage(records, errs, outputWriter)
		case "mean":
			err = getAverage(records, errs, outputWriter)
		case "min":
			err = getMin(records, errs, outputWriter)
		case "minimum":
			err = getMin(records, errs, outputWriter)
		case "max":
			err = getMax(records, errs, outputWriter)
		case "maximum":
			err = getMax(records, errs, outputWriter)
		default:
			return errors.New(fmt.Sprintf(
				"mode %s not recognized.\n"+
					"The mode must be one of the following values: "+
					"'each' 'average' 'mean' 'min' 'minimum' 'max' 'maximum'", lengthMode))
		}
		return err
	},
}

func init() {
	statsCmd.AddCommand(lengthCmd)

	lengthCmd.Flags().StringVarP(&lengthMode, "mode", "m", "each", "How to display lengths")
}

func getEach(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {

	for records != nil && errs != nil {
		select {
		case record := <-records:
			_, err := fmt.Fprintf(output, "%s\t%d\n", record.Name, record.Seq.Length())
			if err != nil {
				return err
			}
		case err := <-errs:
			return err
		}
	}

	return nil
}

func getAverage(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	total, count := 0, 0

	for records != nil && errs != nil {
		select {
		case record := <-records:
			total += record.Seq.Length()
			count++
		case err := <-errs:
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(output, float32(total)/float32(count))
			return err
		}
	}
	return nil
}

func getMin(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	min := -1

	for records != nil && errs != nil {
		select {
		case record := <-records:
			if min < 0 || record.Seq.Length() < min {
				min = record.Seq.Length()
			}
		case err := <-errs:
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(output, min)
			return err
		}
	}
	return nil
}

func getMax(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	max := -1

	for records != nil && errs != nil {
		select {
		case record := <-records:
			if max < 0 || record.Seq.Length() > max {
				max = record.Seq.Length()
			}
		case err := <-errs:
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(output, max)
			return err
		}
	}
	return nil
}
