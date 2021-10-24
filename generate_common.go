package main

import "fmt"

func GenerateImage(i interface{}) error {
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
