package seqs

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// The Seq type is a string representation of a biological sequence
type Seq string

// SeqRecord holds the string representation of a sequence and it's name
type SeqRecord struct {
	Name string
	Seq  Seq
}

// Length returns the number of characters in a sequence
func (seq *Seq) Length() int {
	return len(*seq)
}

// FormatSeq transforms the linear representation of a sequence in multiline representation of the sequence
// each line being of length `width`. This `width` parameter must be an int > 0. If the original sequence
// length is not divisible by `width`, then there will be `Length(sequence) // width` lines and a last line
// with `Length(sequence) % width` characters on it.
func (seq *Seq) FormatSeq(width int) (string, error) {
	if width <= 0 {
		return "", errors.New("width of fasta line must be > 0")
	}
	num := (seq.Length() - 1) / (width + 1)
	chunks := make([]string, 0, num)
	length, start := 0, 0
	for i := 0; i < seq.Length(); i++ {
		if length == width {
			chunks = append(chunks, string(*seq)[start:i])
			length = 0
			start = i
		}
		length++
	}
	chunks = append(chunks, string(*seq)[start:])

	return strings.Join(chunks, "\n"), nil
}

// ReadFastaRecords takes a fasta formated input stream and outputs a collection of SeqRecords to the `output` channels.
func ReadFastaRecords(input io.Reader, output chan SeqRecord, errs chan error) {
	defer close(output)
	defer close(errs)

	scanner := bufio.NewScanner(input)

	var seq Seq
	name := ""

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if line[0] == '>' {
			if name != "" {
				output <- SeqRecord{Name: name, Seq: seq}
			}
			name = strings.TrimSpace(line[1:])
			seq = ""
		} else {
			seq += Seq(strings.TrimSpace(line))
		}

	}

	if err := scanner.Err(); err != nil {
		errs <- err
	}

	output <- SeqRecord{Name: name, Seq: seq}
	errs <- nil
}
