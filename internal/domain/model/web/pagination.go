package web

type Pagination struct {
	Next      int
	Current   int
	Previous  int
	TotalPage int
	Data      int
}

type Offset int
