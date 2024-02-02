package storage_repository

import "mime/multipart"


type StorageRepository interface {
	SaveFile(file multipart.File, fileId string, fileBucket string) (string, error)
	DeleteFile(bucketId string, fileName string) (error)
}
