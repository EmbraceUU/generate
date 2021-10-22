package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/xuri/excelize/v2"
	"runtime"
	"sync"
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

	minRep = float64(1)
)

const (
	SHEET    = "generate"
	SheetII  = "snap"
	RandNum  = 100
	PoolSize = 20

	FileNameFormatFirstChar = 'a'
	Transfer                = false
	GenerateImage           = true
	RootPath                = "/Users/admin/Desktop/mystery-box/"
	TempDir                 = "temp"
	OutDir                  = "out"

	ConfExcelName = "rarity.xlsx"
)

func main() {
	runtime.GOMAXPROCS(runNum)

	fmt.Println("rand begin, ", startTime)

	var err error

	goPool, _ = ants.NewPoolWithFunc(PoolSize, func(fs interface{}) {
		_ = generateImage(fs)
		wg.Done()
	})
	defer goPool.Release()

	f, errR := excelize.OpenFile(fmt.Sprintf("%s/%s", RootPath, ConfExcelName))
	if errR != nil {
		fmt.Println("open conf excel file failed, err: ", errR.Error())
		return
	}

	// 读取组件过程
	err = loadComponent(f)
	if err != nil {
		fmt.Println("load component failed, err: ", err.Error())
		return
	}

	// Create a new sheet.
	index := f.NewSheet(SHEET)
	snapIndex := f.NewSheet(SheetII)

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
	f.SetActiveSheet(index)
	f.SetActiveSheet(snapIndex)
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}

	if GenerateImage {
		for _, g := range taskQueue {
			wg.Add(1)
			_ = goPool.Invoke(g)
		}
		wg.Wait()
	}

	fmt.Println("generate image over. cost ", float64(time.Now().Unix()-startTime)/3600, " h.")
	return
}
