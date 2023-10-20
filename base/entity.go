package base

import (
	"encoding/json"
	"time"
)

type Result struct {
	Code       interface{} `json:"errCode"`
	ErrMessage interface{} `json:"errMessage"`
}

type CheckResult struct {
	Code       interface{} `json:"code"`
	ErrMessage interface{} `json:"msg"`
}

func (r Result) ToJSONStr() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func BuildResult(code int, msg string) Result {
	return Result{
		Code:       code,
		ErrMessage: msg,
	}
}

type DownloadDocResult struct {
	Result
	DownloadDocVo `json:"data"`
}

type DownloadDocVo struct {
	Url     string    `json:"url"`
	Expires time.Time `json:"expires"`
}

type AddOnlyAttachmentResult struct {
	Result
	AddOnlyAttachmentVo `json:"data"`
}

type AddOnlyAttachmentVo struct {
	FileUrl string `json:"fileUrl"`
	FileId  int64  `json:"fileId"`
}

type Message struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}
