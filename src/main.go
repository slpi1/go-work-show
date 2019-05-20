package main

import (
    "errors"
    "fmt"

    "lib"
    "resource"
    "service"

    "github.com/spf13/viper"
)

func main() {

    err := prepare()
    if err != nil {
        lib.Logger().Println("初始化失败")
        return
    }

    // 脚本执行计时
    timer := service.NewTimer()
    timer.Start()

    lib.Logger().Println("准备就绪。。。")

    // 同步数据
    resource.CountSupplier()
}

func prepare() error {
    var err error

    // 加载配置文件
    err = loadConfig()
    if err != nil {
        lib.Logger().Println("配置文件读取失败:" + err.Error())
        return err
    }

    return nil
}

func loadConfig() error {
    viper.SetConfigName("config")
    viper.AddConfigPath("./config")
    err := viper.ReadInConfig()
    if err != nil {
        return err
    }

    config := lib.NewConfig()
    fmt.Println(config)

    if config.Resource.Root == "" || config.Resource.Upload == "" {
        return errors.New("目录配置错误，请确认 resource.root 与 resource.upload 配置项")
    }

    return nil
}