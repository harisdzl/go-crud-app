package image_entity


func ConvertRawImageToImage(rawImage ImageRaw, url string) Image {
	var processedImage Image

	processedImage.ID = rawImage.ID
	processedImage.Caption = rawImage.Caption
	processedImage.ProductID = rawImage.ProductID
	processedImage.Url = url

	return processedImage
}