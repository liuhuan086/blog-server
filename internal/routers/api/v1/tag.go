package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-server/global"
	"github.com/go-programming-tour-book/blog-server/internal/service"
	"github.com/go-programming-tour-book/blog-server/pkg/app"
	"github.com/go-programming-tour-book/blog-server/pkg/convert"
	"github.com/go-programming-tour-book/blog-server/pkg/errcode"
)

type Tag struct{}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

func NewTag() Tag {
	return Tag{}
}

func (t Tag) Get(c *gin.Context) {}

func (t Tag) List(c *gin.Context) {
	//params := struct {
	//	Name  string `form:"name" binding:"max=100"`
	//	State uint8  `form:"state, default=1" binding:"oneof=0 1"`
	//}{}
	//response := app.NewResponse(c)
	//valid, errs := app.BindAndValid(c, &params)
	//if valid == true {
	//	global.Logger.Errorf("app.BindAndValid errs: %v", errs)
	//	response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
	//	return
	//  }
	//  response.ToResponse(gin.H{})
	//  return

	param := service.TagListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	pager := app.Pager{
		Page:     app.GetPage(c),
		PageSize: app.GetPageSize(c),
	}

	TotalRows, err := svc.CountTag(&service.CountTagRequest{
		Name:  param.Name,
		State: param.State,
	})
	if err != nil {
		global.Logger.Errorf(c, "svg.CountTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}

	tags, err := svc.GetTagList(&param, &pager)
	if err != nil {
		global.Logger.Errorf(c, "svg.GetTagList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}

	response.ToResponseList(tags, TotalRows)
	return
}

func (t Tag) Create(c *gin.Context) {
	param := service.CreateTagRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&param)

	if err != nil {
		global.Logger.Errorf(c, "svc.CreateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

func (t Tag) Update(c *gin.Context) {
	param := service.UpdateTagRequest{
		ID: convert.StrTo(c.Param("id")).MustInt32(),
	}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)

	if !valid {
		global.Logger.Errorf(c, "svc.CreateTag err: %v", errs)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UpdateTag(&param)

	if err != nil {
		global.Logger.Errorf(c, "svc.UpdateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}
func (t Tag) Delete(c *gin.Context) {
	param := service.DeleteTagRequest{ID: convert.StrTo(c.Param("id")).MustInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)

	if !valid {
		global.Logger.Errorf(c, "svc.DeleteTag err: %v", errs)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}

	svc := service.New(c)
	err := svc.DeleteTag(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.DeleteTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}
