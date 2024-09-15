package utils

import (
	"ChaiLabs/config"
	"ChaiLabs/constants"
	"encoding/hex"
	"fmt"
	"net/http"
)

func GetFileTypeFromBuffer(buffer []byte) constants.FileType {
	if len(buffer) < 512 {
		return constants.FileType{Ext: "unknown", Mime: "application/octet-stream"}
	}
	// First attempt with magic number
	magicNumber := hex.EncodeToString(buffer[:4])

	var parsedMagicNumber string
	for _, v := range magicNumber {
		if v >= 'a' && v <= 'z' {
			parsedMagicNumber = parsedMagicNumber + string(v-32)
		} else if v >= 'A' && v <= 'Z' {
			parsedMagicNumber = parsedMagicNumber + string(v)
		} else if v >= '0' && v <= '9' {
			parsedMagicNumber = parsedMagicNumber + string(v)
		}
	}

	fileTypeInfo, exists := constants.MagicNumbers[constants.MagicNumber(parsedMagicNumber)]
	if !exists {
		// Fallback to using MIME type detection
		mime := http.DetectContentType(buffer)
		return constants.FileType{Ext: "unknown", Mime: mime}
	}
	return fileTypeInfo
}

func ValidateFileType(file []byte, allowedFileTypes []string) bool {
	fileType := GetFileTypeFromBuffer(file)

	fmt.Printf("File Type Images: %+v\n", fileType)
	for _, allowedType := range allowedFileTypes {
		if fileType.Ext == allowedType {
			return true
		}
	}
	return false
}

func ValidateFileSize(file []byte) bool {
	oneMBInBytes := 1024 * 1024
	fileSizeInMB := len(file) / oneMBInBytes
	return fileSizeInMB <= config.MaxImageUploadSizeInMb()
}

type ValidationResult struct {
	AreValid bool
	Message  string
}

func ValidateFiles(files [][]byte, allowedFileTypes []string, maxSizeInMB int) ValidationResult {
	result := ValidationResult{
		AreValid: true,
		Message:  "",
	}

	for _, file := range files {
		if !ValidateFileType(file, allowedFileTypes) {
			fileType := GetFileTypeFromBuffer(file)
			result.AreValid = false
			result.Message = fmt.Sprintf("Invalid file type, %s is not allowed. Allowed file types are: %s.",
				fileType.Ext, allowedFileTypes)
			break
		}

		if !ValidateFileSize(file) {
			result.AreValid = false
			result.Message = fmt.Sprintf("File size limit exceeded, Max allowed file size is %dMB.", maxSizeInMB)
			break
		}
	}
	return result
}
