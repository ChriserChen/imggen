package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// LiveResult is LiveResult
type LiveResult struct {
	Status  int
	Message string
	Data    LiveData
}

// LiveData is data
type LiveData struct {
	Page       int
	PageSize   int
	TotalPage  int
	TotalCount int
	Datas      []AnchorInfo
}

// AnchorInfo is AnchorInfo
type AnchorInfo struct {
	GameFullName string
	RoomName     string
	Screenshot   string
	Nick         string
	Avatar180    string
	Introduction string
}

type CategoryResponse struct {
	Status int
	Result CategoryResult
}

type CategoryResult struct {
	Hot  []Category
	User []Category
}

type Category struct {
	Host  interface{}
	Name  string
	IsHot int
}

func main() {
	// 查询分类
	var categorys = getCategory() // 打印输出
	if categorys == nil {
		fmt.Println("获取类别为空")
		return
	}
	for index, category := range categorys.Hot {
		fmt.Println("Key: ", index, " 类别: ", category.Name)
	}

	var categoryType int
	fmt.Println("请输入类别编号:")
	fmt.Scanln(&categoryType)

	// 分页查询
	category := categorys.Hot[categoryType]
	liveResult := getLiveListByPage(category.Host.(string), 1)

	// 下载图片
	downloadImages(liveResult.Data.Datas)
}

// 分页查询: https://www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&gameId=2168&tagAll=0&page=1
// -1: 查询全部
func getLiveListByPage(gameId string, page int) *LiveResult {
	response, err := http.Get("https://www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&gameId=" + gameId + "&tagAll=0&page=" + strconv.Itoa(page))
	if err != nil {
		fmt.Println("执行 Get 请求发生异常: ", err)
		return nil
	}
	body, _ := ioutil.ReadAll(response.Body)
	var result LiveResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("解析 JSON 发生异常: ", err)
		return nil
	}
	return &result
}

// getCategory: 获取分类
// 分类：https://www.huya.com//cache10min.php?m=Game&do=ajaxNavGame&callback=huyaNavCategory
func getCategory() *CategoryResult {
	response, err := http.Get("https://www.huya.com//cache10min.php?m=Game&do=ajaxNavGame")
	if err != nil {
		fmt.Println("执行 Get 请求发生异常: ", err)
		return nil
	}
	bytes, _ := ioutil.ReadAll(response.Body)
	var categoryResponse CategoryResponse
	err = json.Unmarshal(bytes, &categoryResponse)
	if err != nil {
		fmt.Println("解析 JSON 发生异常: ", err)
		return nil
	}
	return &categoryResponse.Result
}

// 写入文件
func writeImageFile(fileName, imageURL string) {
	response, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("执行 Get 请求发生异常: ", err)
		return
	}
	defer response.Body.Close()

	imageFileName := "./images/" + fileName
	file, err := os.Create(imageFileName)
	if err != nil {
		fmt.Println("创建文件: [ ", imageFileName, " ] 发生异常: ", err)
		return
	}

	absPath, _ := filepath.Abs(imageFileName)
	fmt.Println("文件地址: ",absPath)
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("下载文件发生异常: ", err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 下载文件
func downloadImages(aAnchorInfos []AnchorInfo) {
	isExists := exists("./images/")
	if isExists {
		os.RemoveAll("./images/")
	}
	os.Mkdir("./images/", os.ModePerm)
	for _, anchorInfo := range aAnchorInfos {
		fmt.Println(anchorInfo.Screenshot)
		fileName := anchorInfo.GameFullName + "_" + anchorInfo.RoomName + "_" + anchorInfo.Nick + ".jpeg"
		writeImageFile(fileName, anchorInfo.Screenshot)
	}
}
