package main

import "fmt"

func GenerateImage(i interface{}) error {
	var err error

	param := i.(GenerateImageParam)
	fs := param.Fs
	fileName := param.FileName

	defer func() {
		addCount()
		Infoln("generate current finished proportion: ", currentProportion(), " count: ", currentCount())
	}()

	Infoln("generate image begin, ", fileName)

	// 生成一个除background以外的临时图片
	testPath := fmt.Sprintf("%s%s/%s.png", GenerateConfig.RootPath, GenerateConfig.OutDir, fileName)
	err = OverlayImage(fs, testPath)
	if err != nil {
		Error("generate image failed, ", fileName, " ", err.Error())
		return err
	}

	Error("generate image finished, ", fileName)
	return nil
}
