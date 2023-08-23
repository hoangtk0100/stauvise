package validator

import (
	"errors"
)

var (
	ErrPasswordTooShort          = errors.New("password too short, password must be at least 8 characters")
	ErrFileExtensionNotSupported = errors.New("file extension not supported")

	fileExts = []string{
		".mp4", ".mov", ".avi", ".mkv", ".flv", ".webm",
		".mp3", ".aac", ".wav", ".flac", ".ogg",
		".ts", ".m2ts", ".3gp", ".wmv",
	}
)

func ValidatePassword(input string) error {
	if len(input) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func ValidateFileExt(input string) error {
	for _, ext := range fileExts {
		if input == ext {
			return nil
		}
	}

	return ErrFileExtensionNotSupported
}
