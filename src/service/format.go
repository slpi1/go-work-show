package service

import (
	"strings"
	"path"
	"path/filepath"
	"fmt"
	"strconv"
	"os"
	"os/exec"
    "sort"
	"time"

	"lib"
)
var debug bool
var root string
var upload string

var queue chan string
var threads chan string
var workerNum int

// 闲置时长
var wait int

// 队列是否被激活
var queueStatus bool = false
var queueLength int


func InitQueue(){
	config := lib.NewConfig()

	debug = config.Debug

	queueLen := config.Exec.Queue
	workerNum = config.Exec.Worker
	wait = config.Exec.Wait

	queue = make(chan string, queueLen)
	threads = make(chan string, workerNum)
}

func FormatFile(file string){
	thumb := GetThumbPath(file)

	if !Exists(thumb) {	
		queue <- file
		if debug {
			queueLength++
			fmt.Println("queue length:", queueLength)
		}
	}
}

func StartExec(){

	// wait 时长内没有新的任务，就会关闭程序
	timeout := 600 * time.Second
	to := time.NewTimer(timeout)

	taskNum := 0
	for {
		to.Reset(timeout)
		select {
			case file := <-queue:
				taskNum++

				// 获取一个可用的工作进程
				threads <- "on"
				fmt.Printf("++++ start %d\r", taskNum)
				go execOneTask(file, taskNum)

			case <-to.C:
				Finish()
				
		}
	}
}

func Finish(){
	var timer = NewTimer()
	timer.End()
	timer.Diff("start", "end")
	os.Exit(0)
}

func execOneTask(file string, id int){
	defer func() {
		// 执行完毕后释放工作进程
		<-threads
	}()

	thumb := GetThumbPath(file)

	if !Exists(thumb) {	
		Thumb(file)
	}

	preview := GetPreviewPath(file)
	if !Exists(preview) {
		Preview(file)
	}

	// var timer = NewTimer()
	// timer.TimePoint(fmt.Sprintf("%d",id))
	fmt.Printf("---- quit %d\r", id)
}

func Thumb(file string) error {

	thumb := GetThumbPath(file)

	err := AccessParentDir(thumb)
	if err != nil {
		return nil
	}

	if IsImg(file) {
		converImg(file, thumb,  "300x300")	
	}

	if IsVideo(file) {
		makeGif(file, thumb)
	}
	
	return nil
}

func Preview(file string) error {
	var preview string
	tmp := ReplaceRoot(file)
	if IsImg(file) {
		preview = resize(tmp, 1000)

		err := AccessParentDir(preview)
		if err != nil {
			return nil
		}
		converImg(file, preview,  "1000x1000")
	}
	if IsVideo(file) {
		preview =  m3u8(tmp)

		err := AccessParentDir(preview)
		if err != nil {
			return nil
		}
		makeM3u8(file, preview)
	}
	return nil
}

func GetThumbPath(path string)string{
	tmp := ReplaceRoot(path)
	return resize(tmp, 300)
}

func GetPreviewPath(path string)string{
	tmp := ReplaceRoot(path)
	if IsImg(path) {
		return resize(tmp, 1000)
	}

	if IsVideo(path) {
		return m3u8(tmp)
	}
	return ""
}

func ReplaceRoot(path string) string {
	if root == "" {
		config := lib.NewConfig()

		root = config.Resource.Root
		upload = config.Resource.Upload
	}

	return strings.Replace(path, root, upload, 1)
}

func resize(p string, width int) string {
	ext := path.Ext(p)
	resizeExt :=  "_" + strconv.Itoa(width) +  translateExt(ext)

	return strings.Replace(p, ext, resizeExt, 1)
}

func m3u8(p string) string {
	ext := path.Ext(p)

	return strings.Replace(p, ext, ".m3u8", 1)
}

func makeGif(origin string, target string) error {
    video := &VideoInfo{origin,0,0,0}
    video.Parse()

    width := 300
    height := video.Height * width / video.Width 
    resize := fmt.Sprintf("%dx%d", width, height)
    if video.Duration < 5 {
		cmd := exec.Command("ffmpeg.exe","-i",origin,"-vframes", "100", "-to", "3", "-y", "-f", "gif", "-fs", "100000",  "-s", resize, target)

		if err := cmd.Run(); err != nil {
			lib.Logger().Println(origin, "compressFile decode error", err.Error())
			return err
		}
    }else{
		cmd := exec.Command("ffmpeg.exe", "-i", origin, "-y", "-f", "image2", "-ss", "5", "-vframes", "1", "-s", resize, target)

		if err := cmd.Run(); err != nil {
			lib.Logger().Println(origin, "compressFile decode error", err.Error())
			return err
		}
    }
	return nil
}

func makeM3u8(origin string, target string) error {

	cmd := exec.Command("ffmpeg", "-i", origin, "-vcodec", "libx264", "-y", target)
	if err := cmd.Run(); err != nil {
		lib.Logger().Println(origin, "compressFile decode error", err.Error())
		return err
	}
	return nil
}


func converImg(file string, target string, resize string) error {

	if IsGif(file) {
		// convert "{temp}" -coalesce -resize "{resize}" -fuzz 5% +dither -layers Optimize +map "{target}"
		cmd := exec.Command("convert.exe", file, "-coalesce", "-resize", resize, "-fuzz", "5%", "+dither", "-layers", "Optimize", "+map", target)		
		if err := cmd.Run(); err != nil {
			lib.Logger().Println(file, "compressFile decode error", err.Error())
			return err
		}
	}else{

		cmd := exec.Command("convert.exe", "-resize", resize, file, target)		
		if err := cmd.Run(); err != nil {
			lib.Logger().Println(file, "compressFile decode error", err.Error())
			return err
		}
	}
	return nil
}

func AccessParentDir(file string) error {
	dir := filepath.Dir(file)

	if !Exists(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			lib.Logger().Println("Dir Create Failed:", dir)
			return err
		}
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsImg(file string) bool {
	config := lib.NewConfig()
	var imgExts = config.Resource.Img
	ext := strings.ToLower(path.Ext(file))

	if  inArray(ext, imgExts)  {
		return true
	}

	return false
}

func IsVideo(file string) bool {
	config := lib.NewConfig()
	var videoExts = config.Resource.Video
	ext := strings.ToLower(path.Ext(file))

	if inArray(ext, videoExts) {
		return true
	}

	return false
}

func IsGif(file string) bool {

	ext := strings.ToLower(path.Ext(file))

	if ext == ".gif" {
		return true
	}

	return false
}

func translateExt(ext string) string {
	config := lib.NewConfig()
	var imgExts = config.Resource.Img
	var videoExts = config.Resource.Video

	ext = strings.ToLower(ext) 

	if inArray(ext, imgExts) {
		return ext
	}

	if inArray(ext, videoExts) {
		return ".gif"
	}
	return ".jpg"
}

// 检查文件格式是否支持
func CheckFileType(file os.FileInfo) bool {
	config := lib.NewConfig()
	var ext = path.Ext(file.Name())
	var allowExts = config.Resource.Exts

	ext = strings.ToLower(ext)
	return inArray(ext, allowExts)
}

func inArray(search string, arr []string) bool {
	index := sort.SearchStrings(arr, search)
	return (index < len(arr) && arr[index] == search)
}