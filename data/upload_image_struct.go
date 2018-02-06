package data

import (
	"mime/multipart"
)

// UploadImage is exported because it starts with a capital letter
type UploadImage struct {
	Bucket    string
	APIKey    string
	Region    string
	Image     multipart.File
	ImageByte []byte
	ImageName string
}
