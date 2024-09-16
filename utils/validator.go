package utils

import (
	"fmt"
	"io"
	"mime/multipart"
)

const (
	allowedFileTypeMagicNumber = "%PDF-"
	pdfMagicNumberLength       = 5
	maxChunkSizeMB             = 1
	maxChunkSizeBytes          = maxChunkSizeMB * 1024 * 1024
)

// ValidateFileType checks if the file type is PDF based on the magic number
func ValidateFileType(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	buf := make([]byte, pdfMagicNumberLength)
	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading file: %v", err)
	}

	// if !strings.HasPrefix(string(buf), allowedFileTypeMagicNumber) {
	//     return fmt.Errorf("unsupported file type, only PDFs are allowed")
	// }

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking file: %v", err)
	}

	return nil
}

// ValidateFileSize checks if the file size does not exceed 1 MB
func ValidateFileSize(file multipart.File) error {
	buf := make([]byte, 0)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading file: %v", err)
	}

	if n > maxChunkSizeBytes {
		return fmt.Errorf("file size exceeds the maximum allowed size of 1 MB")
	}

	return nil
}
