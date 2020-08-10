package service

import (
	"errors"
	"github.com/go-programming-tour-book/blog-server/global"
	"github.com/go-programming-tour-book/blog-server/pkg/upload"
	"mime/multipart"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File,
	fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)
	uploadSavePath := upload.GetSavePath()

	dst := uploadSavePath + "/" + fileName

	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported")
	}

	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("insufficient file permissions")
	}

	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}

