package utils

import (
	"io"
	"net/http"
	"net/url"
	"os"
)

// GetParams 定义一个名为 GetParams 的工具函数，接收一个 URL 字符串作为参数，并返回一个 map[string]string 类型的查询参数
func GetParams(urlString string) (map[string]string, error) {
	// 创建一个空的 map[string]string 用于存储结果
	result := make(map[string]string)

	// 解析 URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return result, err
	}
	// 获取 URL 中的查询参数
	queryParams := parsedURL.Query()

	// 将查询参数转为 map[string]string
	for key, values := range queryParams {
		// 只取第一个值作为结果
		if len(values) > 0 {
			result[key] = values[0]
		}
	}

	return result, nil
}

// DownloadFile 定义一个名为 DownloadFile 的工具函数，接收一个 URL 字符串和目标文件路径作为参数，并下载文件保存到本地
func DownloadFile(url string, filePath string) error {
	// 发起 GET 请求获取文件内容
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将文件内容拷贝到本地文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
