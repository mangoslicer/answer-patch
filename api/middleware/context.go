package middleware

import (
	"github.com/mangoslicer/answer-patch/datastores"
	"github.com/mangoslicer/answer-patch/models"
	auth "github.com/mangoslicer/answer-patch/services"
)

type Context struct {
	*auth.AuthContext
	RepStore    datastores.RepStoreServices
	ParsedModel models.ModelServices
}

func NewContext() *Context {
	return &Context{new(auth.AuthContext), nil, nil}
}
