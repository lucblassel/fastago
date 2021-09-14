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
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
)

var inputFileName string
var outputFileName string
var inputCompression string

var inputReader io.Reader
var outputWriter io.Writer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fastago",
	Short: "Useful commands to work with fasta files",
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
	cobra.OnInitialize(initReader, initWriter)

	rootCmd.PersistentFlags().StringVarP(&inputFileName, "input", "i", "", "input file (default is stdin)")
	rootCmd.PersistentFlags().StringVarP(&outputFileName, "output", "o", "", "output file (default is stdout)")
	rootCmd.PersistentFlags().StringVarP(&inputCompression, "compression", "c", "", "compression mode of file (can be autodected from file extension)")

}

func initWriter() {
	var err error

	if outputFileName != "" {
		outputWriter, err = os.Create(outputFileName)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		outputWriter = os.Stdout
	}
}

func initReader(){
	var reader io.Reader
	var err error

	if inputFileName != "" {
		reader, err = os.Open(inputFileName)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		reader = os.Stdin
	}

	if inputCompression == "" {
		inputCompression = filepath.Ext(inputFileName)[1:]
	}

	inputReader, err =  deCompress(inputCompression, &reader)
	if err != nil {
		log.Fatal(err)
	}
}



func deCompress(compAlg string, reader *io.Reader) (io.Reader, error) {
	switch compAlg {
	case "gz":
		return gzip.NewReader(*reader)
	case "bz2":
		return bzip2.NewReader(*reader), nil
	case "xz":
		return xz.NewReader(*reader)
	default:
		return *reader, nil
	}
}