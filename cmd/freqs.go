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
	"io"
	"sort"
	"strings"

	"github.com/lucblassel/fastago/pkg/seqs"
	"github.com/spf13/cobra"
)

var freqsMode string

// freqsCmd represents the frequency command
var freqsCmd = &cobra.Command{
	Use:   "freqs",
	Short: "get frequencies of sequences in fasta file",
	Long: `By default this command outputs the frequencies averaged over sequences. 
	It is also possible to retreive the frequencies per sequence.
	set the -m/--mode flag to:
		- each: will display the frequencies of each sequence`,
	RunE: func(cmd *cobra.Command, args []string) error {

		records := make(chan seqs.SeqRecord)
		errs := make(chan error)

		go seqs.ReadFastaRecords(inputReader, records, errs)

		var err error

		switch freqsMode {
		case "each":
			err = getEachFreqs(records, errs, outputWriter)
		default:
			err = getAverageFreqs(records, errs, outputWriter)
		}

		return err
	},
}

func init() {
	statsCmd.AddCommand(freqsCmd)

	freqsCmd.Flags().StringVarP(&freqsMode, "mode", "m", "", "How to display frequencies")
}

func getFreqs(record seqs.SeqRecord, errs chan error, output io.Writer) error {
	var totalCount float64 = 0 
	var countmap = make(map[rune]float64)

	for _, char := range strings.ToUpper(string(record.Seq)) {
		_, isPresent := countmap[char]
		if(isPresent) {
			countmap[char]+=1
		}else{
			countmap[char]=1	
		}
		totalCount+=1
	}

	for char, count := range countmap {
		countmap[char] = count/totalCount
	}

	fmt.Fprintf(output, "%s\t", record.Name)
	err := printmap(countmap, "  ", "\t\t", errs, output)
	fmt.Fprintf(output, "\n")

	return err
}

func getEachFreqs(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	
	for records != nil && errs != nil {
		select {
		case record := <-records:
			err := getFreqs(record, errs, outputWriter)
			if err != nil {
				return err
			}
		case err := <-errs:
			return err
		}
	}

	return nil
}

func getAverageFreqs(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	var totalCount float64 = 0 
	var countmap = make(map[rune]float64)
	
	for records != nil && errs != nil {
		select {
		case record := <-records:	
			for _, char := range strings.ToUpper(string(record.Seq)) {
				_, isPresent := countmap[char]
				if(isPresent) {
					countmap[char]+=1
				}else{
					countmap[char]=1	
				}
				totalCount+=1
			}
		case err := <-errs:
			if err != nil {
				return err
			}
			for char, count := range countmap {
				countmap[char] = count/totalCount
			}
			err = printmap(countmap, "\t", "\n", errs, output)
			return err
		}
	}

	return nil
}

func printmap(myMap map[rune]float64, sep string, endline string, errs chan error, output io.Writer) error {
	keys := make([]rune, 0, len(myMap))
    for k := range myMap {
        keys = append(keys, k)
    }
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
        _, err := fmt.Fprintf(output, "%s%s%.2f%s", string(k), sep, myMap[k], endline)
		if err != nil {
			return err
		}
    }

	return nil

}
