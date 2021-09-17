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
- **addid** [🏳](#addid) : add a prefix or a suffix to sequence names 
- **rename** [🏳](#rename) : rename sequences with either a regex or a map file
- **stats** : get statistics and information on the sequences
  - **count** : count sequences in file
  - **length** [🏳](#length) : get length of sequences in file *(can also output the average/min/max)* 
- **subset** [🏳](#subset) : subset the files, keeping only specified sequences. Works with regex, a file of names or positional arguments
- **transform** : apply transformtaion functions to sequences
  - **upper** : transform sequence bases to uppercase
  - **lower** : transform sequence bases to uppercase
- **help** : show usage message
- **version** : get current version of fastago
- **completion** : generate autocompletion script for bash, zsh, fish or powershell *(thank you [cobra](https://github.com/spf13/cobra) 🙏)*
  
## General flags
- `-i` or `--input`: specify the input file. The default is stdin.
- `-o` or `--output`: specify the output file. The default is stdout.
- `-c` or `--compression`: specify the compression method of the input file. If the input on stdin is compressed, this flag must be specified for this tool to work. If this flag is not used and `-i` is, fastago will try to guess the compression from the file extension. Supported compression schemes: `gzip (.gz)`, `bzip2 (.bz2)`, `xzip (.xz)`
- `-w` or `--linewidth`: specify the width at which a sequence will wrap in the output. The default is 80 characters.
- `-h` or `--help` : display a help message.


## Per command flags
### addid
You must specify at least one of the following flagas to run this command: 
 - `-p` or `--prefix` to add your identifier to the beginning of each sequence name
 - `s` or `--suffix` to add your identifier to the end of each sequence name

### rename
There are 2 ways to rename sequences: 
 - The `-m` or `--map` flag allows you to specify a mapping of names to be renamed. On each line of this file you must write the name of the sequence you want to change and the new name, separated by a `tab` character. 
 - The `-r` or `--regex` flag, allows you to specify a regular expression that will match a substring in each sequence name. This match will be replace by the value specified with the `-p`or `--replace` flag. If you provide a regular expression you must also provide a replacement string. More info on Go regular expression syntax [here](https://pkg.go.dev/regexp/syntax).

### stats
#### length 
With the `-m` or `--mode` flag you can choose which information you want to display: 
 - `-m each` : will display the length of each sequence after it's name on a single line
 - `-m average` or `-m mean` will display the average length of all sequences in the file
 - `-m min` or `-m minimum` will display the length of the shortest sequence in the file
 - `-m max` or `-m maximum` will display the length of the longest sequence in the file

The default value for this flag is `each`.

### subset
There are 2 ways to subset your fasta file: 
 - You can use the `-n` or `--names` flag to specify a file of names to keep (1 by line)
 - You can use the `r` or `--regex` flag to specify a regular expression that matches the sequences you want to keep. 

If you specify the `-x` or `--exclude` flag you specify the sequences to exclude instead of the sequences to keep.

