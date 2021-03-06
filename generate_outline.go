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

	Infoln("generate image begin, ", fileName)

	defer func() {
		addCount()
		Infoln("generate current finished proportion: ", currentProportion(), " count: ", currentCount())
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
	// 只处理以前带bg的
	if strings.Contains(fileName, "bg") {
		fileName = strings.Replace(fileName, "bg", "a", 1)
		fileName = strings.Replace(fileName, "e", "f", 1)
		fileName = strings.Replace(fileName, "d", "e", 1)
		fileName = strings.Replace(fileName, "c", "d", 1)
		fileName = strings.Replace(fileName, "b", "c", 1)
		fileName = strings.Replace(fileName, "a", "b", 1)

		splits := strings.Split(fileName, "a")
		return "a" + splits[1] + splits[0]
	}

	return fileName
}
