package internal

import "time"

type Image struct {
	ID            string
	UserID        string
	Filename      string
	ContentType   string
	SizeBytes     int64
	Width         int32
	Height        int32
	UploadedAt    time.Time
	OriginalPath  string
	ThumbnailPath string
}
