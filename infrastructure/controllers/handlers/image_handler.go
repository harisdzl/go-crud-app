package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/image_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/image_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/storage"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type Image struct {
	ImageRepo image_repository.ImageRepository
	Persistence *base.Persistence
}



func NewImage(p *base.Persistence) *Image {
	return &Image{
		Persistence: p,
	}
}


func (im *Image) SaveImage(c *gin.Context) {
	rawImage := image_entity.Image{}

	rawImage.Caption = c.PostForm("caption")
	rawImage.ProductID = c.PostForm("productID")

	image, imageErr := c.FormFile("image")
	openedImage, openImageErr := image.Open()
	if imageErr != nil {
		log.Println(imageErr)
	}

	if openImageErr != nil {
		log.Println(openImageErr)
	}


	im.ImageRepo = application.NewImageApplication(im.Persistence)
	processedImage, processedImageErr := im.ImageRepo.SaveImage(&rawImage)
	
	if processedImageErr != nil {
		log.Println(processedImageErr)
	}
	storage := storage.NewStorageRepository("Supabase", im.Persistence)

	publicUrl, saveFileErr := storage.SaveFile(openedImage, fmt.Sprint(processedImage.ID), "images")

	if saveFileErr != nil {
		log.Println(saveFileErr)
	}
	
	savedImage := im.ImageRepo.UpdateImageURL(publicUrl, processedImage)
	finalImage, _ := im.ImageRepo.UpdateImage(savedImage)
	c.JSON(http.StatusOK, finalImage)
}

func (im *Image) GetImage(c *gin.Context) {
	imageId, err := strconv.ParseInt((c.Param("image_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	im.ImageRepo = application.NewImageApplication(im.Persistence)

	image, err := im.ImageRepo.GetImage(imageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, image)
} 
 
func (im *Image) GetAllImagesOfProduct(c *gin.Context) {
	im.ImageRepo = application.NewImageApplication(im.Persistence)
	productId, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	allImages, err := im.ImageRepo.GetAllImagesOfProduct(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allImages)
}


func (im *Image) DeleteImage(c *gin.Context) {
	imageId, err := strconv.ParseInt((c.Param("image_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	im.ImageRepo = application.NewImageApplication(im.Persistence)

	
	deleteErr := im.ImageRepo.DeleteImage(imageId)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, deleteErr.Error())
		return
	}

	c.JSON(http.StatusOK, "Image deleted")
}
