package entity

type BaseResponse[T any] interface {
	BindError(err error)
	BindData(data T)
}
