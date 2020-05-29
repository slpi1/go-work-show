package lib

import (
    "sync"
    "sort"

    "github.com/spf13/viper"
)

type db struct {
    Username string
    Password string
    Url string
}

type resource struct {
    CoverNum int
    Root string
    Upload string 
    Person string
    Company string
    AudioPath string
    Exts []string
    Img []string
    Video []string
    Audio []string
}

type logger struct {
    Path string
}

type execConfig struct {
    Old bool
    Queue int
    Worker int
    Wait int
    Convert string
}

type Config struct {
    Debug bool
    Mock bool
    Db db
    Resource resource
    Log logger
    Exec execConfig
}


var once sync.Once
var config Config


func NewConfig() Config {
    once.Do(func() {
        config = Config{}

        config.Debug = viper.GetBool("debug")
        config.Mock = viper.GetBool("mock")
        config.Db = loadDB()
        config.Resource = loadResource()
        config.Log = loadLog()
        config.Exec = loadExec()
    })
    return config
}

func loadDB() db {
    var config = db{}
    config.Username = viper.GetString("db.username")
    config.Password = viper.GetString("db.password")
    config.Url = viper.GetString("db.url")

    return config
}

func loadResource() resource {
    var config = resource{}

    config.CoverNum = viper.GetInt("resource.coverNum")
    config.Root = viper.GetString("resource.root")
    config.Upload = viper.GetString("resource.upload")
    config.Person = viper.GetString("resource.person")
    config.Company = viper.GetString("resource.company")
    config.AudioPath = viper.GetString("resource.audioPath")

    var exts = viper.GetStringSlice("resource.exts")
    var img = viper.GetStringSlice("resource.img")
    var video = viper.GetStringSlice("resource.video")
    var audio = viper.GetStringSlice("resource.audio")

    sort.Strings(exts)
    sort.Strings(img)
    sort.Strings(video)
    sort.Strings(audio)

    config.Exts = exts
    config.Img = img
    config.Video = video
    config.Audio = video
    return config
}
func loadLog() logger {
    var config = logger{}
    config.Path = viper.GetString("log.path")
    return config
}
func loadExec() execConfig {
    var config = execConfig{}
    config.Old = viper.GetBool("exec.old")
    config.Queue = viper.GetInt("exec.queue")
    config.Worker = viper.GetInt("exec.worker")
    config.Wait = viper.GetInt("exec.wait")
    config.Convert = viper.GetString("exec.convert")
    return config
}