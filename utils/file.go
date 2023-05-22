package utils

import (
	"github.com/pkg/errors"
	"io"
)

// GetFileSize returns the size of a file in bytes.
func GetFileSize(file io.Reader) (int64, error) {
	// Check if the reader also implements the Seeker interface
	seeker, ok := file.(io.Seeker)
	if !ok {
		return 0, errors.New("file size retrieval not supported")
	}

	// Get the current offset position
	currentPos, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Seek to the end to get the file size
	fileSize, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Restore the original offset position
	_, err = seeker.Seek(currentPos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	// Return the file size
	return fileSize, nil
}
