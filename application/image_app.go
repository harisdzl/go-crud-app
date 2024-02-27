package application

import (
	"context"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/image_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/image_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/images"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type ImageApp struct {
	p *base.Persistence
	c context.Context
}

func NewImageApplication(p *base.Persistence, c context.Context) image_repository.ImageRepository {
	return &ImageApp{p, c}
}

func (a *ImageApp) SaveImage(image *image_entity.Image) (*image_entity.Image, map[string]string) {
	repoImage := images.NewImageRepository(a.p, a.c)
	return repoImage.SaveImage(image)
}


func (a *ImageApp) GetImage(imageId int64) (*image_entity.Image, error) {
	repoImage := images.NewImageRepository(a.p, a.c)
	return repoImage.GetImage(imageId)
}


func (a *ImageApp) GetAllImagesOfProduct(imageId int64) ([]image_entity.Image, error) {
	repoImage := images.NewImageRepository(a.p, a.c)
	return repoImage.GetAllImagesOfProduct(imageId)
}


func (a *ImageApp) UpdateImage(image *image_entity.Image) (*image_entity.Image, error) {
	repoImage := images.NewImageRepository(a.p, a.c)
	return repoImage.UpdateImage(image)
}


func (a *ImageApp) DeleteImage(bucketId string, fileName string) (error) {
	repoImage := images.NewImageRepository(a.p, a.c)
	return repoImage.DeleteImage(bucketId, fileName)
}

func (a *ImageApp) UpdateImageURL(url string, image *image_entity.Image) (*image_entity.Image) {
	image.Url = url
	return image
}