package global

import (
	"github.com/go-programming-tour-book/blog-server/pkg/logger"
	"github.com/go-programming-tour-book/blog-server/pkg/setting"
)

var (
	SeverSetting    *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
	JWTSetting      *setting.JWTSettingS
	EmailSetting    *setting.EmailSettingS
)
