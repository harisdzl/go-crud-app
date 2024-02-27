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

//	@Summary		Save Image
//	@Description	Saves an image along with its metadata.
//	@Tags			Image
//	@Accept			mpfd
//	@Produce		json
//	@Param			caption		formData	string					false	"Image caption"
//	@Param			product_id	formData	int64					true	"Product ID"
//	@Param			image		formData	file					true	"Image file to upload"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/images [post]
func (im *Image) SaveImage(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	rawImage := image_entity.Image{}

	rawImage.Caption = c.PostForm("caption")
	rawImage.ProductID, _ = strconv.ParseInt(c.PostForm("product_id"), 10, 64) 

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

	im.ImageRepo = application.NewImageApplication(im.Persistence, c)
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

//	@Summary		Get Image
//	@Description	Retrieves an image by its ID.
//	@Tags			Image
//	@Accept			json
//	@Produce		json
//	@Param			image_id	path		int						true	"Image ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/images/{image_id} [get]
func (im *Image) GetImage(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	imageID, err := strconv.ParseInt((c.Param("image_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	im.ImageRepo = application.NewImageApplication(im.Persistence, c)

	image, err := im.ImageRepo.GetImage(imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Image %v obtained", imageID), image))
}

//	@Summary		Get All Images Of Product
//	@Description	Retrieves all images associated with a product.
//	@Tags			Image
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int						true	"Product ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/products/{product_id}/images [get]
func (im *Image) GetAllImagesOfProduct(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	im.ImageRepo = application.NewImageApplication(im.Persistence, c)
	productID, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	allImages, err := im.ImageRepo.GetAllImagesOfProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	results := map[string]interface{}{
		"results" : allImages,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("All images for product %v obtained", productID), results))
}

//	@Summary		Delete Image
//	@Description	Deletes an image by its ID.
//	@Tags			Image
//	@Accept			json
//	@Produce		json
//	@Param			image_id	path		string					true	"Image ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/images/{image_id} [delete]
func (im *Image) DeleteImage(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	im.ImageRepo = application.NewImageApplication(im.Persistence, c)

	deleteErr := im.ImageRepo.DeleteImage("images", c.Param("image_id"))
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Image deleted successfully", ""))
}