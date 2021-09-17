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
	"bufio"
	"errors"
	"fmt"
	"github.com/lucblassel/fastago/pkg/seqs"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var mapFile string
var regexRenamer string
var replaceGroup string

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename sequences according to map file or regex",
	RunE: func(cmd *cobra.Command, args []string) error {

		if regexRenamer != "" {

			if replaceGroup == "" {
				return errors.New("if using regex renaming the --replace flag must be specified")
			}

			return renameFromRegex(regexRenamer, replaceGroup)
		}

		if mapFile != "" {
			renamer, err := readMap(mapFile)
			if err != nil {
				return err
			}
			return renameFromMap(renamer)
		}

		return errors.New("you must specify a regular expression or a map file to rename sequences")
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&mapFile, "map", "m", "", "Tab separated file mapping old names to new names. 1 operation per line")
	renameCmd.Flags().StringVarP(&regexRenamer, "regex", "r", "", "Regex to match part of the sequence name")
	renameCmd.Flags().StringVarP(&replaceGroup, "replace", "p", "", "Replace matched element with this")
}

func readMap(filename string) (map[string]string, error) {
	names := make(map[string]string)
	input, err := os.Open(filename)
	defer input.Close()

	if err != nil {
		return names, err
	}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "\t")
		names[split[0]] = split[1]
	}

	if err := scanner.Err(); err != nil {
		return names, err
	}

	return names, nil
}

func renameFromMap(renamer map[string]string) error {

	records := make(chan seqs.SeqRecord)
	errs := make(chan error)

	go seqs.ReadFastaRecords(inputReader, records, errs)

	for records != nil && errs != nil {
		select {
		case record := <-records:
			newName := record.Name
			if val, ok := renamer[newName]; ok {
				newName = val
			}
			output, err := record.Seq.FormatSeq(outputLineWidth)
			_, err = fmt.Fprintf(outputWriter, ">%s\n%s\n", newName, output)
			if err != nil {
				return err
			}
		case err := <-errs:
			return err
		}
	}

	return nil
}

func renameFromRegex(expression string, replace string) error {
	regex, err := regexp.Compile(expression)
	if err != nil {
		return err
	}

	records := make(chan seqs.SeqRecord)
	errs := make(chan error)

	go seqs.ReadFastaRecords(inputReader, records, errs)

	for records != nil && errs != nil {
		select {
		case record := <-records:
			newName := regex.ReplaceAllString(record.Name, replace)
			output, err := record.Seq.FormatSeq(outputLineWidth)
			_, err = fmt.Fprintf(outputWriter, ">%s\n%s\n", newName, output)
			if err != nil {
				return err
			}
		case err := <-errs:
			return err
		}
	}

	return nil

}