package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type Client struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

// DownloadOssFile Download oss file
func (oc *Client) DownloadOssFile(object string, savePath string) error {
	client, err := oss.New(oc.Endpoint, oc.AccessKeyId, oc.AccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(oc.BucketName)
	if err != nil {
		return err
	}

	err = bucket.GetObjectToFile(object, savePath)
	if err != nil {
		return err
	}
	return nil
}
