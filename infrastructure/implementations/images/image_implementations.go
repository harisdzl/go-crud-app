package images

import (
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/image_entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

// To manage new Image repositories in the database

// Image Repository struct
type ImageRepo struct {
	p *base.Persistence
}

func NewImageRepository(p *base.Persistence) *ImageRepo {
	return &ImageRepo{p}
}


func (r *ImageRepo) SaveImage(Image *image_entity.Image) (*image_entity.Image, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", r.p)

	dbErr := map[string]string{}
	err := r.p.DB.Debug().Create(&Image).Error
	
	if err != nil {
		fmt.Println("Failed to create Image")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	cacheRepo.SetKey(fmt.Sprintf("%v_IMAGE", Image.ID), Image, time.Minute * 15)
	
	return Image, nil
}


func (r *ImageRepo) GetImage(imageId int64) (*image_entity.Image, error) {
	var image *image_entity.Image

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	_ = cacheRepo.GetKey(fmt.Sprintf("%v_IMAGE", imageId), &image)
	if image == nil {
		err := r.p.DB.Debug().Where("id = ?", imageId).Take(&image).Error
		if err != nil {
			fmt.Println("Failed to get image")
		}
		if image != nil && image.ID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprintf("%v_IMAGE", imageId), image, time.Minute * 15)
		}
	}


	return image, nil
}


func (r *ImageRepo) GetAllImagesOfProduct(productId int64) ([]image_entity.Image, error) {
	var image []image_entity.Image

	err := r.p.DB.Debug().Where("product_id = ?", productId).Find(&image).Error

	if err != nil {
		fmt.Println("Failed to get Image")
	}

	return image, nil
}


func (r *ImageRepo) UpdateImage(image *image_entity.Image) (*image_entity.Image, error) {
	cacheRepo := cache.NewCacheRepository("Redis", r.p)

	err := r.p.DB.Debug().Where("id = ?", image.ID).Updates(&image).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = cacheRepo.SetKey(fmt.Sprintf("%v_IMAGE", image.ID), image, time.Minute * 15)

	return image, nil
}

func (r *ImageRepo) DeleteImage(imageId int64) error {
	var image *image_entity.Image	

	err := r.p.DB.Debug().Where("id = ?", imageId).Delete(&image).Error
	if err != nil {
		return err
	}

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	

	cacheRepo.DelKey(fmt.Sprintf("%v_IMAGE", imageId))
	if err != nil {
		return errors.New("database error, please try again")
	}

	fileName := fmt.Sprint(image.ID)
	_, deleteErr := r.p.DbSupabase.RemoveFile("images", []string{fileName})

	if deleteErr != nil {
		return errors.New("Supabase error, please try again")
	}
	return nil
}
