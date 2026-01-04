package cfg

import (
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func Init(file string) {
	// 加载配置
	viper.SetConfigFile(file)
	lo.Must0(viper.ReadInConfig())
}

func Viper() *viper.Viper {
	return viper.GetViper()
}

func UnmarshalKey[T any](key string) (v T) {
	viper.UnmarshalKey(key, &v)
	return
}
