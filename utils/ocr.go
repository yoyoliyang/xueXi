package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"image/color"
	"image/png"

	"github.com/trazyn/uiautomator-go"
)

type red struct {
	r        uint32
	position struct {
		x int
		y int
	}
}

type ocrResult struct {
	WordsResult []words `json:"words_result"`
}

type words struct {
	Words string `json:"words"`
}

// GetRedFontString 获取“查看提示”里面的红色文字
func GetRedFontString(ua *uiautomator.UIAutomator) (string, error) {

	sc, err := ua.GetScreenshot()
	if err != nil {
		return "", err
	}

	bf, err := getRedFontImageFile(sc.Base64)
	if err != nil {
		return "", err
	}

	str, err := baiduOcrHandle(bf)
	if err != nil {
		return "", err
	}
	fmt.Println(str)

	result := &ocrResult{}

	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		return "", errors.New("Error parsing data returned by OCR_API")
	}

	str = ""
	for _, item := range result.WordsResult {
		str += item.Words
	}

	return str, nil

}

// 图像挖取函数，将图片中的红色字体单独取出，供ocr api识别使用
func getRedFontImageFile(base64str string) (*bytes.Buffer, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64str))

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	// 此处设置一个偏移地址将字体图像偏移到0,0位置，从而缩小图像大小
	var offsetX, tempX, offsetY, tempY int
	var redFontImageMaxX, redFontImageMaxY int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r == 65535 && g == 0 && b == 0 {
				tempX = x
				tempY = y
				// 此处获取x轴和y轴最小的偏移量
				if offsetX == 0 && offsetY == 0 {
					offsetX = tempX
					offsetY = tempY
				}
				if offsetX >= tempX {
					offsetX = tempX
				}
				if offsetY >= tempY {
					offsetY = tempY
				}

				// 获取红色字体最大的x轴和y轴坐标
				if x > redFontImageMaxX {
					redFontImageMaxX = x
				}
				if y > redFontImageMaxY {
					redFontImageMaxY = y - offsetY
				}
			}
		}
	}
	fmt.Println(offsetX, offsetY)
	fmt.Println(redFontImageMaxX, redFontImageMaxY)

	upLeft := image.Point{0, 0}

	lowRight := image.Point{redFontImageMaxX + 8, redFontImageMaxY + 8}

	redFontImage := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	fontColor := color.RGBA{255, 0, 0, 0xff}
	whiteColor := color.RGBA{255, 255, 255, 0xff}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ { // 填充背景为白色
			redFontImage.Set(x, y, whiteColor)
			r, g, b, _ := img.At(x, y).RGBA()
			if r == 65535 && g == 0 && b == 0 {
				redFontImage.Set(x-offsetX+4, y-offsetY+4, fontColor)
				//加粗字体
				redFontImage.Set(x-offsetX+4+1, y-offsetY+4+1, fontColor)
			}
		}
	}

	bs := &bytes.Buffer{}

	png.Encode(bs, redFontImage)

	f, err := os.Create("red_font_image.png")
	if err != nil {
		return nil, err
	}
	png.Encode(f, redFontImage)

	return bs, nil

}

func baiduOcrHandle(b *bytes.Buffer) (string, error) {
	base64Str := base64.StdEncoding.EncodeToString(b.Bytes())

	params := url.Values{}
	params.Set("image", base64Str)

	api := `https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic`
	bdToken := "BAIDU_OCR_TOKEN"
	accessToken := os.Getenv(bdToken)
	if accessToken == "" {
		return "", errors.New("No environment variables were found: " + bdToken)
	}

	api += "?access_token=" + accessToken

	resq, err := http.NewRequest("POST", api, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	resq.Header.Add("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(resq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)

	return string(bs), nil

}
