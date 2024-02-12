package uploader

import (
	"bytes"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/yundun"
	errors "github.com/wuntsong-org/wterrors"
	"path/filepath"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
)

var (
	aliyun = newAliyun()
)

func PutImage(data []byte) (string, errors.WTError) {
	return aliyun.PutImage(data)
}

func CopyImage(originUrl string) (string, errors.WTError) {
	return aliyun.CopyImage(originUrl)
}

func GetImage(key string) (string, errors.WTError) {
	return aliyun.GetImage(key)
}

func newAliyun() *aliyunOssUploader {
	return &aliyunOssUploader{
		once:   sync.Once{},
		bucket: nil,
	}
}

type uploader interface {
	PutImage(data []byte) (string, error)
	PutObject(key string, data []byte) (string, error)
	CopyImage(originUrl string) (string, error)
}

// 阿里云oss
type aliyunOssUploader struct {
	once   sync.Once
	bucket *oss.Bucket
}

func (aliyun *aliyunOssUploader) PutImage(data []byte) (string, errors.WTError) {
	fileType := utils.GetMediaType(data)
	if !utils.IsImage(fileType) {
		return "", errors.Errorf("bad picture")
	}

	ok, err := yundun.CheckBaseLinePic(data, fileType)
	if err != nil {
		logger.Logger.Tag("A", err.Error())
		return "", errors.WarpQuick(err)
	} else if !ok {
		return "", errors.Errorf("bad picture")
	}

	key := generateImageKey(data)
	err = aliyun.putObject(fmt.Sprintf("图像文件/%s", key), data)
	if err != nil {
		logger.Logger.Tag("B", err.Error())
		return "", errors.WarpQuick(err)
	}
	return key, nil
}

func (aliyun *aliyunOssUploader) CopyImage(originUrl string) (string, errors.WTError) {
	data, err := download(originUrl)
	if err != nil {
		return "", errors.WarpQuick(err)
	}
	key, err := aliyun.PutImage(data)
	if err != nil {
		return "", errors.WarpQuick(err)
	}
	return key, nil
}

func (aliyun *aliyunOssUploader) putObject(key string, data []byte) errors.WTError {
	bucket := aliyun.getBucket()
	if err := bucket.PutObject(key, bytes.NewReader(data)); err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func (aliyun *aliyunOssUploader) GetImage(key string) (string, errors.WTError) {
	url, err := aliyun.getObject(fmt.Sprintf("图像文件/%s", key))
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	return url, nil
}

func (aliyun *aliyunOssUploader) getObject(key string) (string, errors.WTError) {
	bucket := aliyun.getBucket()
	url, err := bucket.SignURL(key, oss.HTTPGet, 60)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	return url, nil
}

func (aliyun *aliyunOssUploader) getBucket() *oss.Bucket {
	aliyun.once.Do(func() {
		client, err := oss.New(viper.GetString("aliyun.endpoint"), viper.GetString("aliyun.accessKeyId"), viper.GetString("aliyun.accessKeySecret"))
		if err == nil {
			aliyun.bucket, err = client.Bucket(viper.GetString("aliyun.bucket"))
		}
	})
	return aliyun.bucket
}

// generateKey 生成图片Key
func generateImageKey(data []byte) string {
	md5 := utils.MD5Bytes(data)
	return filepath.Join("images", utils.TimeFormat(time.Now(), "2006/01/02/"), md5+".jpg")
}

func download(url string) ([]byte, errors.WTError) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	return rsp.Body(), nil
}
