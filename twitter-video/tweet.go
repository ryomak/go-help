package main

import (
	"encoding/base64"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	"os"
	"time"
)

var text string = ""

var chunkSize int = 5 * 1024 * 1024

func GetFileSize(file *os.File) (int, error) {
	fileStatus, err := file.Stat()
	if err != nil {
		return 0, nil
	}
	return int(fileStatus.Size()), nil
}

func main() {
	api := InitTwitterApi()
	fi, err := os.Open("ex.mp4")
	if err != nil {
	}
	defer fi.Close()

	size, err := GetFileSize(fi)
	if err != nil {
	}
	byteMedia, err := ioutil.ReadAll(fi)
	if err != nil {
	}

	chunkedMedia, err := api.UploadVideoInit(size, "video/mp4", "tweet_video")
	if err != nil {
	}
	index := 0
	for i := 0; i < size; i += chunkSize {
		var data string
		if (size-i)/chunkSize > 0 {
			data = base64.StdEncoding.EncodeToString(byteMedia[i : i+chunkSize])
		} else {
			data = base64.StdEncoding.EncodeToString(byteMedia[i:])
		}
		if err = api.UploadVideoAppend(chunkedMedia.MediaIDString, index, data); err != nil {
			fmt.Printf("%v\n", err)
		}
		index++
	}
	video, err := api.UploadVideoFinalize(chunkedMedia.MediaIDString)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	//status
	for {
		videos, err := api.UploadVideoStatus(chunkedMedia.MediaIDString)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if videos.ProcessingInfo.State == "succeeded" {
			break
		}
		fmt.Printf("%d%--%dsec\n", videos.ProcessingInfo.ProgressPercent, videos.ProcessingInfo.CheckAfterSecs)
		time.Sleep(time.Duration(videos.ProcessingInfo.CheckAfterSecs))
	}

	params := url.Values{}
	params.Add("media_ids", video.MediaIDString)
	if _, err := api.PostTweet(text, params); err != nil {

		fmt.Printf("%v\n", err)
	}
}

func InitTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("T_CK"))
	anaconda.SetConsumerSecret(os.Getenv("T_CS"))
	return anaconda.NewTwitterApi(os.Getenv("T_AT"), os.Getenv("T_ATS"))
}
