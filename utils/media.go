package utils

import (
	"net/http"
	"strings"
)

const (
	Png     = 1
	Jpg     = 2
	Bmp     = 3
	Unknown = 4
	Avi     = 5
	Mp4     = 6
	Mp3     = 7
	Wav     = 8
)

var MediaTypeSuffixMap = map[int]string{
	Png:     "png",
	Jpg:     "jpeg",
	Bmp:     "bmp",
	Avi:     "avi",
	Mp4:     "mp4",
	Mp3:     "mp3",
	Wav:     "wav",
	Unknown: "unk",
}

var MediaTypeContentTypeMap = map[int]string{
	Png:     "image/png",
	Jpg:     "image/jpeg",
	Bmp:     "image/bmp",
	Avi:     "video/avi",
	Mp4:     "video/mp4",
	Mp3:     "audio/mpeg",
	Wav:     "audio/wave",
	Unknown: "application/octet-stream",
}

func GetImageType(data []byte) int {
	t := GetMediaType(data)
	if IsImage(t) {
		return t
	}

	return Unknown
}

func GetVideoType(data []byte) int {
	t := GetMediaType(data)
	if IsVideo(t) {
		return t
	}

	return Unknown
}

func GetAudioType(data []byte) int {
	t := GetMediaType(data)
	if IsAudio(t) {
		return t
	}

	return Unknown
}

func GetMediaType(data []byte) int {
	contentType := http.DetectContentType(data)

	if strings.HasPrefix(contentType, "image/png") {
		return Png
	}

	if strings.HasPrefix(contentType, "image/jpeg") {
		return Jpg
	}

	if strings.HasPrefix(contentType, "image/bmp") {
		return Bmp
	}

	if strings.HasPrefix(contentType, "video/avi") {
		return Avi
	}

	if strings.HasPrefix(contentType, "video/mp4") {
		return Mp4
	}

	if strings.HasPrefix(contentType, "audio/mpeg") {
		return Mp3
	}

	if strings.HasPrefix(contentType, "audio/wave") {
		return Wav
	}

	return Unknown
}

func IsImage(fileType int) bool {
	switch fileType {
	case Png, Jpg, Bmp:
		return true
	default:
		return false
	}
}

func IsVideo(fileType int) bool {
	switch fileType {
	case Mp4, Avi:
		return true
	default:
		return false
	}
}

func IsAudio(fileType int) bool {
	switch fileType {
	case Mp3, Wav:
		return true
	default:
		return false
	}
}
