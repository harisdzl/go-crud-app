package image_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/image_entity"
)

type ImageRepository interface {
	SaveImage(*image_entity.Image) (*image_entity.Image, map[string]string)
	GetImage(int64) (*image_entity.Image, error)
	GetAllImagesOfProduct(int64) ([]image_entity.Image, error)
	UpdateImage(*image_entity.Image) (*image_entity.Image, error)
	DeleteImage(string, string) error
	UpdateImageURL(string, *image_entity.Image) (*image_entity.Image)
}