package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
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
	responseContextData := entity.ResponseContext{Ctx: c}
	rawImage := image_entity.Image{}

	rawImage.Caption = c.PostForm("caption")
	rawImage.ProductID = c.PostForm("productID")

	image, imageErr := c.FormFile("image")
	openedImage, openImageErr := image.Open()
	if imageErr != nil {
		log.Println(imageErr)
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, imageErr.Error(), ""))
		return
	}

	if openImageErr != nil {
		log.Println(openImageErr)
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, openImageErr.Error(), ""))
		return
	}

	im.ImageRepo = application.NewImageApplication(im.Persistence)
	processedImage, processedImageErr := im.ImageRepo.SaveImage(&rawImage)

	if processedImageErr != nil {
		log.Println(processedImageErr)
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, processedImageErr["db_error"], ""))
		return
	}

	storage := storage.NewStorageRepository("Supabase", im.Persistence)
	publicURL, saveFileErr := storage.SaveFile(openedImage, fmt.Sprint(processedImage.ID), "images")

	if saveFileErr != nil {
		log.Println(saveFileErr)
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, saveFileErr.Error(), ""))
		return
	}

	savedImage := im.ImageRepo.UpdateImageURL(publicURL, processedImage)
	finalImage, _ := im.ImageRepo.UpdateImage(savedImage)
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Image saved successfully", finalImage))
}

func (im *Image) GetImage(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	imageID, err := strconv.ParseInt((c.Param("image_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	im.ImageRepo = application.NewImageApplication(im.Persistence)

	image, err := im.ImageRepo.GetImage(imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Image %v obtained", imageID), image))
}

func (im *Image) GetAllImagesOfProduct(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	im.ImageRepo = application.NewImageApplication(im.Persistence)
	productID, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	allImages, err := im.ImageRepo.GetAllImagesOfProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("All images for product %v obtained", productID), allImages))
}

func (im *Image) DeleteImage(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	im.ImageRepo = application.NewImageApplication(im.Persistence)

	log.Println(c.Param("image_id"))
	deleteErr := im.ImageRepo.DeleteImage("images", c.Param("image_id"))
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Image deleted successfully", ""))
}