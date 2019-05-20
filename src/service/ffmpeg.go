package service


import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type VideoInfo struct {
	Path string
	Duration int
	Width int
	Height int
}

func (V *VideoInfo) Parse() error {
	cmd := exec.Command("ffmpeg", "-i", V.Path)
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderrPipe.Close()
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(stderrPipe)
	for {
		line, err := reader.ReadBytes('\r')
		if err != nil || err == io.EOF {
			break
		}
		//匹配视频时长
		V.getDuration(string(line))

		// 匹配尺寸
		V.getSize(string(line))

	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (V *VideoInfo)getDuration(line string){

	reg1 := regexp.MustCompile(`Duration:(.*?),`)
	snatch1 := reg1.FindStringSubmatch(line)
	if len(snatch1) > 1 {
		V.Duration = timeEncode(snatch1[1])
	}
}

func (V *VideoInfo)getSize(line string){

	reg := regexp.MustCompile(`Stream(.*?)(\d{3,4}x\d{3,4})`)
	snatch := reg.FindStringSubmatch(string(line))
	if len(snatch) > 2 {
		size := snatch[2]
		info := strings.Split(size,"x")
		V.Width,_ = strconv.Atoi(info[0])
		V.Height,_ = strconv.Atoi(info[1])

	}
}

func timeEncode(t string) int {
	time := strings.Trim(t, " ")
	hour, _ := strconv.Atoi(time[:2])
	minute, _ := strconv.Atoi(time[3:5])
	second, _ := strconv.Atoi(time[6:8])
	return second + minute*60 + hour*3600
}