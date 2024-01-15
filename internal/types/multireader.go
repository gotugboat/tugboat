package types

import "io"

// multiReadCloser is a type that combines multiple io.ReadClosers into one.
type MultiReadCloser struct {
	Readers []io.ReadCloser
}

// Read reads from the combined readers.
func (mrc *MultiReadCloser) Read(p []byte) (n int, err error) {
	for _, reader := range mrc.Readers {
		n, err = reader.Read(p)
		if err != nil && err != io.EOF {
			return n, err
		}
		if n > 0 {
			return n, nil
		}
	}
	return 0, io.EOF
}

// Close closes all combined readers.
func (mrc *MultiReadCloser) Close() error {
	var err error
	for _, reader := range mrc.Readers {
		if closeErr := reader.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}
