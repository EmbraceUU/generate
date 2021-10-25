package main

import (
	"fmt"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/panjf2000/ants/v2"
	"github.com/xuri/excelize/v2"
	"image"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	sheetValue   = make(map[string]string)
	duplicateSet = make(map[string]struct{})

	// key -> sheetName  value -> Rows
	componentMap   = make(map[string][]FileRarity)
	componentUsage = make(map[string]UsageCount)
	compPriority   = make([]string, 0)

	goPool *ants.PoolWithFunc
	wg     sync.WaitGroup

	runNum = runtime.NumCPU()

	startTime = time.Now().Unix()

	taskQueue = make([]GenerateImageParam, 0)

	minRep = float64(-1)

	fulfillment int64
)

func GenerateRun() {
	runtime.GOMAXPROCS(runNum)

	fmt.Println("rand begin, ", startTime)

	var err error

	goPool, _ = ants.NewPoolWithFunc(GenerateConfig.PoolSize, func(fs interface{}) {
		_ = GenerateOutline(fs)
		wg.Done()
	})
	defer goPool.Release()

	f, errR := excelize.OpenFile(fmt.Sprintf("%s/%s", GenerateConfig.RootPath, GenerateConfig.ExcelName))
	if errR != nil {
		fmt.Println("open conf excel file failed, err: ", errR.Error())
		return
	}

	err = loadDuplicateSet()
	if err != nil {
		fmt.Println("load duplicate failed, err: ", err.Error())
		return
	}

	// 读取组件过程
	err = loadComponent(f)
	if err != nil {
		fmt.Println("load component failed, err: ", err.Error())
		return
	}

	// Create a new sheet.
	_ = f.NewSheet(GenerateConfig.OutSheet)
	_ = f.NewSheet(GenerateConfig.SnapSheet)

	// 随机过程, 组合组件列表
	err = randProcess()
	if err != nil {
		fmt.Println("rand process failed, err: ", err.Error())
		return
	}

	// 生成sheet map
	err = generateSheetMapValue(f)
	if err != nil {
		fmt.Println("generate sheet map failed, err: ", err.Error())
		return
	}

	// 保存到excel文件
	//f.SetActiveSheet(index)
	//f.SetActiveSheet(snapIndex)
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}

	if GenerateConfig.IfGenerateImage {
		for _, g := range taskQueue {
			wg.Add(1)
			_ = goPool.Invoke(g)
		}
		wg.Wait()
	}

	fmt.Println("generate image over. cost ", float64(time.Now().Unix()-startTime)/3600, " h.")
}

func addCount() {
	atomic.AddInt64(&fulfillment, 1)
}

func currentProportion() float64 {
	return float64(atomic.LoadInt64(&fulfillment)) / float64(GenerateConfig.RandNum)
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
		_ = f.SetCellValue(GenerateConfig.OutSheet, k, v)
	}

	for k, v := range sheetValue {
		_ = f.SetCellValue(GenerateConfig.OutSheet, k, v)
	}

	snapHeader := map[string]string{
		"A1": "AttributeName",
		"B1": "Remain",
		"C1": "Usage",
		"D1": "Usage Ratio",
		"E1": "Rarity",
	}

	for k, v := range snapHeader {
		_ = f.SetCellValue(GenerateConfig.SnapSheet, k, v)
	}

	index := 2
	for k, ug := range componentUsage {
		_ = f.SetCellValue(GenerateConfig.SnapSheet, fmt.Sprintf("A%d", index), k)
		_ = f.SetCellValue(GenerateConfig.SnapSheet, fmt.Sprintf("B%d", index), ug.Remain)

		_ = f.SetCellValue(GenerateConfig.SnapSheet, fmt.Sprintf("C%d", index), ug.Available-ug.Remain)
		_ = f.SetCellValue(GenerateConfig.SnapSheet, fmt.Sprintf("D%d", index), (ug.Available-ug.Remain)/float64(GenerateConfig.RandNum)*100)
		_ = f.SetCellValue(GenerateConfig.SnapSheet, fmt.Sprintf("E%d", index), ug.Rarity)
		index++
	}

	return err
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
