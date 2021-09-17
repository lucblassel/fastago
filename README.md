# Fastago

## Presentation
This is a very simple tool to do basic operations on fasta formatted files.   

 It is inspired by [goalign](https://github.com/evolbioinfo/goalign), with the main difference that this tool streams the files and executes operations per sequence. This eliminates the need to load the whole file into memory.  
 This tool is not meant to work on alignments, it is simply a collection of useful "ease of life" functions when working with fasta files. 

## Examples
```bash
# Transform all characters to uppercase and save to "upper.fasta"
$> fastago transform upper --input input.fasta --output upper.fasta

# Count sequences in file
$> fastago stats count --input input.fasta 

# Read from stdin
$> cat input.fasta | fastago stats length  

# chain operations
$> cat input.fasta | fastago transform lower | fastago stats count

# read compressed file
$> fastago stats count --input input.fasta.gz

# pipe in compressed data
$> cat input.fasta.xz | fastago stats count --compression xz

# Extract "Seq_1" and "Seq_2" sequences from input file
$> fastago subset --input input.fasta Seq_1 Seq_2 
```

## Installation
### Binaries
Grab a binary from the [releases](https://github.com/lucblassel/fastago/releases/) section.

### Build from source
Make sure you have [go `1.16` installed](https://golang.org/doc/install). Then clone this repository and build with:  
`go build` or `go build -o <binaryName>` if you want to specify another name than `fastago`. 

## Commands
- **rename** : rename sequences with either a regex or a map file `(TODO)`
- **stats** : get statistics and information on the sequences
  - **count** : count sequences in file
  - **length** : get length of sequences in file. *(can also output the average/min/max)*
- **subset** : subset the files, keeping only specified sequences. Works with regex, a file of names or positional arguments. 
- **transform** : apply transformtaion functions to sequences
  - **upper** : transform sequence bases to uppercase.
  - **lower** : transform sequence bases to uppercase.
- **version** : get current version of fastago
- **completion** : generate autocompletion script for bash, zsh, fish or powershell *(thank you [cobra](https://github.com/spf13/cobra) 🙏)*
  
## General flags
- `-i` or `--input`: specify the input file. The default is stdin.
- `-o` or `--output`: specify the output file. The default is stdout.
- `-c` or `--compression`: specify the compression method of the input file. If the input on stdin is compressed, this flag must be specified for this tool to work. If this flag is not used and `-i` is, fastago will try to guess the compression from the file extension. Supported compression schemes: `gzip (.gz)`, `bzip2 (.bz2)`, `xzip (.xz)`
- `-w` or `--linewidth`: specify the width at which a sequence will wrap in the output. The default is 80 characters.
- `-h` or `--help` : display a help message.


## Per command flags
### stats
#### length
With the `-m` or `--mode` flag you can choose which information you want to display: 
 - `-m each` : will display the length of each sequence after it's name on a single line
 - `-m average` or `-m mean` will display the average length of all sequences in the file
 - `-m min` or `-m minimum` will display the length of the shortest sequence in the file
 - `-m max` or `-m maximum` will display the length of the longest sequence in the file

The default value for this flag is `each`.

### subset
With the `-n` or `--names` flag you can specify a file where you have the names of the sequences you want to keep. The tool expects a name per line in the file.
