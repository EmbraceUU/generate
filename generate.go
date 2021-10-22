package main

import (
	"fmt"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/xuri/excelize/v2"
	"image"
	"os"
	"sync/atomic"
)

var (
	fulfillment int64
)

func addCount() {
	atomic.AddInt64(&fulfillment, 1)
}

func currentProportion() float64 {
	return float64(atomic.LoadInt64(&fulfillment)) / RandNum
}

func currentCount() int64 {
	return atomic.LoadInt64(&fulfillment)
}

// 生成sheetMap
func generateSheetMapValue(f *excelize.File) error {
	var err error

	header := map[string]string{}
	headerSlice := compPriority
	headerSlice = append(headerSlice, "filename")
	for axis, i := 'A', 1; i <= len(headerSlice); axis, i = axis+1, i+1 {
		header[fmt.Sprintf("%c%d", axis, 1)] = headerSlice[i-1]
	}

	for k, v := range header {
		_ = f.SetCellValue(SHEET, k, v)
	}

	for k, v := range sheetValue {
		_ = f.SetCellValue(SHEET, k, v)
	}

	index := 1
	for k, v := range componentRep {
		_ = f.SetCellValue(SheetII, fmt.Sprintf("A%d", index), k)
		_ = f.SetCellValue(SheetII, fmt.Sprintf("B%d", index), v)
		index++
	}

	return err
}

func generateImage(i interface{}) error {
	var err error

	param := i.(GenerateImageParam)
	fs := param.Fs
	fileName := param.FileName

	defer func() {
		addCount()
		fmt.Println("generate current finished proportion: ", currentProportion(), " count: ", currentCount())
	}()

	fmt.Println("generate image begin, ", fileName)

	// 生成一个除background以外的临时图片
	testPath := fmt.Sprintf("%s%s/%s.png", RootPath, OutDir, fileName)
	err = OverlayImage(fs, testPath)
	if err != nil {
		fmt.Println("generate image failed, ", fileName, " ", err.Error())
		return err
	}

	fmt.Println("generate image finished, ", fileName)
	return nil
}

func OverlayImage(fs []string, dst string) error {
	if fs == nil || len(fs) <= 1 {
		return fmt.Errorf("overlay image source file is nil. ")
	}

	if dst == "" {
		return fmt.Errorf("overlay image filename is nil. ")
	}

	var images []*image.Image
	for _, fn := range fs {
		img, err := LoadImage(fn)
		if err != nil {
			return err
		}
		images = append(images, img)
	}

	if images == nil || len(images) == 0 {
		return fmt.Errorf("overlay image failed, None of them work")
	}

	bg := images[0]
	result := clone.AsRGBA(*bg)
	for i, img := range images {
		if i == 0 {
			continue
		}

		result = blend.Normal(result, *img)
	}

	return imgio.Save(dst, result, imgio.PNGEncoder())
}

func LoadImage(filePath string) (*image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	bg, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return &bg, nil
}
