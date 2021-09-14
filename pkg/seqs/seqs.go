package seqs

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Seq string

type SeqRecord struct {
	Name string
	Seq  Seq
}

func (seq *Seq) Length() int {
	return len(*seq)
}

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

func ReadFastaRecords(input io.Reader, output chan SeqRecord, errs chan error) {
	defer close(output)
	defer close(errs)

	scanner := bufio.NewScanner(input)

	var seq Seq
	name := ""

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
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
