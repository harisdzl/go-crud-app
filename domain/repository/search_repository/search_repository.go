package search_repository

type SearchRepository interface {
	InsertDoc(collectionName string, src interface{}) (error)
	UpdateDoc(productId uint, collectionName string, updatedFields interface{}) (error)
	DeleteSingleDoc(fieldName string, collectionName string, id int64) (error)
	DeleteMultipleDoc(fieldName string, collectionName string, id int64) (error)
	DeleteAllDoc(collectionName string, src []interface{}) (error)
	InsertAllDoc(collectionName string, src []interface{}) (error)
	SearchDocByName(name string, indexName string, src interface{}) (error)
}
