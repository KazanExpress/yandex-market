package models

type PagingOptions struct {
	Page int32
	Size int32
}

type LimitOptions struct {
	Limit  int32
	Offset int32
}
