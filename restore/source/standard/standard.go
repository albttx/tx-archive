package standard

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gnolang/gno/tm2/pkg/amino"
	"github.com/gnolang/gno/tm2/pkg/std"
	"github.com/gnolang/tx-archive/types"
)

type Standard struct {
	file    *os.File
	scanner *bufio.Scanner
}

// NewStandardSource creates a new standard JSON source
func NewStandardSource(filePath string) (*Standard, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to open standard input file, %w",
			err,
		)
	}

	return &Standard{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (s *Standard) Next(ctx context.Context) (*std.Tx, error) {
	// Read the line
	if s.scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, io.EOF
		default:
			// Parse the JSON
			var tx types.TxData

			txt := s.scanner.Text()
			fmt.Println(txt)

			if err := amino.UnmarshalJSON([]byte(txt), &tx); err != nil {
				return nil, fmt.Errorf(
					"unable to unmarshal amino JSON, %w",
					err,
				)
			}

			return &tx.Tx, nil
		}
	}

	// Check for scanning errors
	if err := s.scanner.Err(); err != nil {
		return nil, fmt.Errorf(
			"unable to read standard input file, %w",
			err,
		)
	}

	return nil, io.EOF
}

func (s *Standard) Close() error {
	// Attempt to gracefully close the file
	if closeErr := s.file.Close(); closeErr != nil {
		return fmt.Errorf(
			"unable to gracefully close standard file, %w",
			closeErr,
		)
	}

	return nil
}
