/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"os"

	"github.com/spf13/cobra"
)

var inputFileName string
var outputFileName string
var compression string

var inputReader io.Reader
var outputWriter io.Writer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fastago",
	Short: "Usefule commands to work with fasta files",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initIO)

	rootCmd.PersistentFlags().StringVarP(&inputFileName, "input", "i", "", "input file (default is stdin)")
	rootCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", "", "output file (default is stdout)")
	rootCmd.PersistentFlags().StringVarP(&compression, "compression", "c", "", "compression mode of file")
}

func initIO() {
	var err error
	if inputFileName != "" {
		inputReader, err = os.Open(inputFileName)
		if err != nil {
			log.Fatal("Error opening input file: ", err)
		}
	} else {
		inputReader = os.Stdin
	}

	if outputFileName != "" {
		outputWriter, err = os.Create(outputFileName)
		if err != nil {
			log.Fatal("Error opening output file: ", err)
		}
	} else {
		outputWriter = os.Stdout
	}
}
