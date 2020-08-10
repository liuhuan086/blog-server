package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-server/global"
	"github.com/go-programming-tour-book/blog-server/internal/service"
	"github.com/go-programming-tour-book/blog-server/pkg/app"
	"github.com/go-programming-tour-book/blog-server/pkg/convert"
	"github.com/go-programming-tour-book/blog-server/pkg/errcode"
	"github.com/go-programming-tour-book/blog-server/pkg/upload"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
