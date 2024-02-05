package supabase

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/harisquqo/quqo-challenge-1/domain/repository/storage_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	storage_go "github.com/supabase-community/storage-go"
)


type supabaseRepo struct {
	p *base.Persistence
}


func (s supabaseRepo) SaveFile(file multipart.File, fileId string, fileBucket string) (string, error) {
	filePath := fmt.Sprintf("%v", fileId)
	imageType := "image/png"
	imageTypePointer := &imageType
	fileOptions := storage_go.FileOptions{
		ContentType: imageTypePointer, 
	}
	
	_, uploadErr := s.p.DbSupabase.
	UploadFile(fileBucket, filePath, file, fileOptions)

	if uploadErr != nil {
		log.Println(uploadErr)
	}

	publicUrl := s.p.DbSupabase.GetPublicUrl(fileBucket, filePath).SignedURL

	return publicUrl, nil

}


func (s supabaseRepo) DeleteFile(bucketId string, fileName string) error {
	response, err := s.p.DbSupabase.RemoveFile(bucketId, []string{fileName})
	if err != nil {
		log.Println(err)
		return err // Return the actual error
	}

	log.Println(response)

	// If there is no error, return nil
	return nil
}



func NewStorageRepository(p *base.Persistence) storage_repository.StorageRepository {
	return &supabaseRepo{p}
}