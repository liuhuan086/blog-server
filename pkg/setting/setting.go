package setting

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

// 用于初始化本地项目配置的基础属性，即设定配置文件名称，路径，类型及相对路径
func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("../blog-server/configs/")
	//vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp: vp}, nil
}
