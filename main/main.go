package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Result struct {
	Status  int
	Message string
	Data    Data
}

type Data struct {
	Page       int
	PageSize   int
	TotalPage  int
	TotalCount int
	Datas      []SubData
}

type SubData struct {
	GameFullName string
	RoomName     string
	Screenshot   string
	Nick         string
	Avatar180    string
	Introduction string
}

func main() {
	//分类：https://www.huya.com//cache10min.php?m=Game&do=ajaxNavGame&callback=huyaNavCategory
	//分页查询：https://www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&gameId=2168&tagAll=0&page=1
	response, _ := http.Get("https://www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&gameId=2168&tagAll=0&page=1")
	body, _ := ioutil.ReadAll(response.Body)

	var result Result
	err := json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("解析 JSON 发生异常: ", err)
		return
	}

	for _, anchorInfo := range result.Data.Datas {
		fmt.Println(anchorInfo.Screenshot)
		fileName := anchorInfo.GameFullName + "_" + anchorInfo.RoomName + "_" + anchorInfo.Nick + ".jpeg"
		writeImageFile(fileName, anchorInfo.Screenshot)
	}
}

func writeImageFile(fileName, imageUrl string) {
	response, err := http.Get(imageUrl)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("执行 Get 请求发生异常: ", err)
		return
	}
	imageFileName := "./static/images/" + fileName
	file, err := os.Create(imageFileName)
	if err != nil {
		fmt.Println("创建文件: [ ", imageFileName, " ] 发生异常: ", err)
		return

	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("下载文件发生异常: ", err)
	}
}
