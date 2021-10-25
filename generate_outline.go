package main

import (
	"fmt"
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"image"
	"image/color"
	"strings"
)

const OutLineLength = 1

func GenerateOutline(i interface{}) error {
	var (
		err     error
		olLeft  *image.RGBA
		olRight *image.RGBA
		olUp    *image.RGBA
		olDown  *image.RGBA
	)

	param := i.(GenerateImageParam)
	fs := param.Fs
	fileName := param.FileName

	fmt.Println("generate image begin, ", fileName)

	defer func() {
		addCount()
		fmt.Println("generate current finished proportion: ", currentProportion(), " count: ", currentCount())
	}()

	// 生成一个除background以外的临时图片
	portion, err := OverlayRGBA(fs[1:])
	if err != nil {
		return err
	}

	fn := func(c color.RGBA) color.RGBA {
		return color.RGBA{A: c.A}
	}
	olPortion := adjust.Apply(portion, fn)
	// 重复n次, 进行描边
	for index := 0; index < 30; index++ {
		olRight = transform.Translate(olPortion, OutLineLength, 0)
		olLeft = transform.Translate(olPortion, -OutLineLength, 0)
		olUp = transform.Translate(olPortion, 0, OutLineLength)
		olDown = transform.Translate(olPortion, 0, -OutLineLength)

		olPortion, err = OverlayRGBAByFile(olRight, olLeft, olUp, olDown)
		if err != nil {
			return err
		}
	}

	// 读取背景
	bg, err := LoadImage(fs[0])
	if err != nil {
		return err
	}

	result, err := OverlayRGBAByFile(*bg, olPortion, portion)
	if err != nil {
		return err
	}

	testPath := fmt.Sprintf("%s%s/%s.png", GenerateConfig.RootPath, GenerateConfig.OutDir, fileName)
	err = imgio.Save(testPath, result, imgio.PNGEncoder())
	if err != nil {
		return err
	}
	return nil
}

func OverlayRGBAByFile(fs ...image.Image) (*image.RGBA, error) {
	bg := fs[0]
	portion := clone.AsRGBA(bg)
	for index, img := range fs {
		if index == 0 {
			continue
		}

		portion = blend.Normal(portion, img)
	}

	return portion, nil
}

func OverlayRGBA(fs []string) (*image.RGBA, error) {
	var images []*image.Image
	for _, fn := range fs {
		img, err := LoadImage(fn)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	if images == nil || len(images) == 0 {
		return nil, fmt.Errorf("overlay image failed, None of them work")
	}

	bg := images[0]
	portion := clone.AsRGBA(*bg)
	for index, img := range images {
		if index == 0 {
			continue
		}

		portion = blend.Normal(portion, *img)
	}

	return portion, nil
}

func transferFileName(fileName string) string {
	fileName = strings.Replace(fileName, "bg", "a", 1)
	fileName = strings.Replace(fileName, "e", "f", 1)
	fileName = strings.Replace(fileName, "d", "e", 1)
	fileName = strings.Replace(fileName, "c", "d", 1)
	fileName = strings.Replace(fileName, "b", "c", 1)
	fileName = strings.Replace(fileName, "a", "b", 1)

	splits := strings.Split(fileName, "a")
	return "a" + splits[1] + splits[0]
}

//func OutLine(source *image.RGBA, length int, fileName string) error {
//	right := transform.Translate(source, length, 0)
//	err := imgio.Save(fmt.Sprintf("%s/%s-temp2.png", tempPath, fileName), right, imgio.PNGEncoder())
//	if err != nil {
//		return err
//	}
//
//	left := transform.Translate(source, -length, 0)
//	err = imgio.Save(fmt.Sprintf("%s/%s-temp3.png", tempPath, fileName), left, imgio.PNGEncoder())
//	if err != nil {
//		return err
//	}
//
//	up := transform.Translate(source, 0, length)
//	err = imgio.Save(fmt.Sprintf("%s/%s-temp4.png", tempPath, fileName), up, imgio.PNGEncoder())
//	if err != nil {
//		return err
//	}
//
//	down := transform.Translate(source, 0, -length)
//	err = imgio.Save(fmt.Sprintf("%s/%s-temp5.png", tempPath, fileName), down, imgio.PNGEncoder())
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func BendOutline(fileName string) (*image.RGBA, error) {
//	var fs []string
//	fs = append(fs, fmt.Sprintf("%s/%s-temp2.png", tempPath, fileName))
//	fs = append(fs, fmt.Sprintf("%s/%s-temp3.png", tempPath, fileName))
//	fs = append(fs, fmt.Sprintf("%s/%s-temp4.png", tempPath, fileName))
//	fs = append(fs, fmt.Sprintf("%s/%s-temp5.png", tempPath, fileName))
//
//	err := OverlayImage(fs, fmt.Sprintf("%s/%s-temp-outline.png", tempPath, fileName))
//	if err != nil {
//		return nil, err
//	}
//
//	img, err := LoadImage(fmt.Sprintf("%s/%s-temp-outline.png", tempPath, fileName))
//	if err != nil {
//		return nil, err
//	}
//
//	return clone.AsRGBA(*img), nil
//}
