package data

import (
	"mime/multipart"
)

// UploadImage is exported because it starts with a capital letter
type UploadImage struct {
	Bucket    string
	AccessKey string
	Region    string
	Image     multipart.File
	ImageByte []byte
	ImageName string
}
