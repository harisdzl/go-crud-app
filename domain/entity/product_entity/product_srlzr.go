package product_entity


func ConvertProductInventoryToProduct(productForInventory ProductForInventory) Product {
	var product Product

	product.ID = productForInventory.ID
	product.Name = productForInventory.Name
	product.Description = productForInventory.Description
	product.Price = productForInventory.Price
	product.Category = productForInventory.Category

	return product

}