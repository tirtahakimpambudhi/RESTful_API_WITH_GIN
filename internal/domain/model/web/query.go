package web

import (
	"go_gin/internal/config"
	"strconv"
)

type SearchQuery struct {
	Page   string `form:"page" validate:"numeric"`
	Search string `form:"search"`
}

type TodoListByIDQuery struct {
	ID string `form:"id" validate:"required,numeric"`
}
type TodoListByIDValue struct {
	ID int
}

type TodoListByIDsQuery struct {
	IDs []string `form:"id" validate:"required,dive,numeric"`
}
type TodoListByIDsValue struct {
	IDs []int
}

type GetAllQuery struct {
	Page string `form:"page" validate:"numeric"`
}

type GetAllValue struct {
	Offset Offset
	Page   int
}

type SearchValue struct {
	Offset Offset
	Page   int
	Search string
}

func (q *SearchQuery) ToValue() (search *SearchValue, err error) {
	if q.Page == "" || q.Page == "0" {
		q.Page = "1"
	}
	page, err := strconv.Atoi(q.Page)
	offset := Offset((page - 1) * config.Other.Limit)
	search = &SearchValue{
		Offset: offset,
		Search: q.Search,
		Page:   page,
	}
	return
}

func (g *GetAllQuery) ToValue() (getAll *GetAllValue, err error) {
	if g.Page == "" || g.Page == "0" {
		g.Page = "1"
	}
	page, err := strconv.Atoi(g.Page)
	offset := Offset((page - 1) * config.Other.Limit)
	getAll = &GetAllValue{offset, page}
	return
}

func (t *TodoListByIDQuery) ToValue() (value *TodoListByIDValue, err error) {
	id, err := strconv.Atoi(t.ID)
	value = &TodoListByIDValue{id}
	return
}

func (q *TodoListByIDsQuery) ToValue() (value *TodoListByIDsValue, err error) {
	var IDs []int
	for _, idstr := range q.IDs {
		id, errParse := strconv.Atoi(idstr)
		if errParse != nil {
			value = &TodoListByIDsValue{}
			err = errParse
			return
		}
		IDs = append(IDs, id)
	}
	value = &TodoListByIDsValue{IDs: IDs}
	return
}
