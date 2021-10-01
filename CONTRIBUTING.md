# Contributing to fastago

Welcome and thank you for considering contributing to this project!   
This is the first time I maintain any type of project like this, and I decided to jump the gun for hacktoberfest, so please bear with me through the growing pains. 

## Code of conduct
I think that having a code of conduct is quite important for good online interactions, [here is the one](https://github.com/auth0/open-source-template/blob/master/CODE-OF-CONDUCT.md) you agree to follow if you choose to contribute to this project. 

## Getting started

If you wish to contribute but have no idea where to start, please check the open issues to see if you find inspiration.  
Contributions should be made with Pull Requests (PRs) and issues. Some general guidelines: 
 - Search existing PRs / issues before creating your own.
 - Issues can be used to: 
   - report a problem (use the `bug` tag)
   - request a new feature (use the `enhancment` tag)
 - If you find an issue that addresses a problem you're having, please add your remarks to that issue instead of starting a new one. 
 - **PR:** Only fix/add the functionality in question OR address wide-spread whitespace/style issues, not both.
 - **PR:** Address a single concern in the least number of changed lines as possible.

Overall I would suggest you follow the ["fork-and-pull" git workflow](https://github.com/susam/gitpr):

 1. Fork the repository to your own Github account
 1. Clone the project to your machine
 1. Create a branch locally with a succinct but descriptive name
 1. Commit changes to the branch
 1. Following any formatting and testing guidelines specific to this repo
 1. Push changes to your fork
 1. Open a PR in our repository


## Developer's guide
### Prerequisites
First of all you need to [install `go 1.16`](https://golang.org/doc/install) if you haven't already.  
Then `git clone` or better yet `go get` this repo and you can get started.   
This is a CLI tool built on [cobra](https://github.com/spf13/cobra), which you can install with `go get -u github.com/spf13/cobra`.

### Understanding the FASTA file format
the FASTA file format is a widely used file format to represent a collection of biological sequences in bioinformatics. It is a very simple format that is easy to work with.  
The idea is simple a file is made up of records containing: 
 - The sequence Name 
 - The sequence itself

The sequence name is on its own line started by the ">" character, sometimes it can have additional metadata after, but always contained on a single line.  
All the following lines, as long as it doesnt start with ">" represents the sequence. Sometimes a empty line can be added between records to improve legibility.  
The following two files would be totaly equivalent:  
```
>Seq1
ATGGCAGTCAGTTGCAGTCATTTAGCAGTCA
>Seq2
CGATCGAGTCAGTCGAGCGGAGTCGATGCAG
>Seq3
CGATCGAGGCTAGCTAGCAGGTCAGCAGTGC
```

```
>Seq1
ATGGCAGTCAGT
TGCAGTCATTTA
GCAGTCA

>Seq2
CGATCGAGTCAG
TCGAGCGGAGTC
GATGCAG

>Seq3
CGATCGAGGCTA
GCTAGCAGGTCA
GCAGTGC
```
They both represent 3 sequences, called `Seq1`, `Seq2` and `Seq3`.  

A FASTA file can represent DNA, RNA or proteins, sometimes this is represented by different file extensions: `.fasta`, `.fa`, `.fna`, `.ffn`, `.faa`, `.frn`.

More information can be found [here](https://en.wikipedia.org/wiki/FASTA_format)

Overall it is very simple to work with because we can simply stream over the file, sequence by sequence, and do operations on it without having to load the whole file into memory. 

### Adding commands
You can add a new command using cobra:
```shell
# add a base command
cobra add <commandName>
# add a subcommand
cobra add <subCommandName> -p "<parentCommandName>"
```

This will create a `.go` file in the `cmd` directory that will contain a stub for you to fill in.

### Rules for new commands
I have two main rules when adding a new command: 
 - You must at least include a short description of the command in the `Short` field of the `&cobra.Command` struct. If you want to add a longer description in `Long` that's great as well but not mandatory.
 - You must use `RunE` instead of `Run`. This means that the function run by the cobra command, must return an `error` type. You must therefore use a `func(cmd *cobra.Command, args []string) error` in this field.

Thank you for respecting these rules!

If you have any other question about creating / modifying commands, please refer to the [cobra user guide](https://github.com/spf13/cobra/blob/master/user_guide.md).

### Streaming over the FASTA file
The [`seqs` package](https://github.com/lucblassel/fastago/blob/main/pkg/seqs/seqs.go) is used to read the input fasta file. from an input `io.Reader` stream, it reads `SeqRecord` structs that contain the sequence and the sequence name. These `SeqRecord` structs will be sent in a channel that must be given to the `seqs.ReadFastaRecords` function. An error channel must also be given as to this function. 
