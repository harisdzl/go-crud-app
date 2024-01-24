package search_repository

type SearchRepository interface {
	InsertDoc(src interface{}) (error)
	UpdateDoc(productId uint, updatedFields interface{}) (error)
	DeleteSingleDoc(productId int64) (error)
	DeleteAllDoc(src []interface{}) (error)
	InsertAllDoc(src []interface{}) (error)
	SearchDocByName(name string, src interface{}) (error)
}
