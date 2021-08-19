package seqs

import (
	"bufio"
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
