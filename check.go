package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func CheckRun() {
	var (
		path     = "./"
		confName = "attribute_info.xlsx"
	)

	fmt.Println("Start reviewing media materials. ")
	// 加载目录
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("check failed. open folder %s failed, err: %s. ", path, err.Error()))
		return
	}

	var (
		dirs      int
		confFound bool
	)
	for _, f := range files {
		if f.IsDir() {
			dirs++
			continue
		}

		if strings.EqualFold(f.Name(), confName) {
			confFound = true
		}
	}

	if !confFound {
		fmt.Println("check failed. The Excel file for configuring properties was not found. ")
		return
	}

	confFilePath := fmt.Sprintf("%s/%s", path, confName)
	f, errR := excelize.OpenFile(confFilePath)
	if errR != nil {
		fmt.Println(errR.Error())
		return
	}

	sheetList := f.GetSheetList()
	if sheetList == nil || len(sheetList) == 0 {
		fmt.Println("check failed. Excel is empty. ")
		return
	}

	// 判断sheet页是否比文件夹少
	if len(sheetList)-1 < dirs {
		fmt.Println("check failed. Please check whether there are any missing folders that are not recorded in Excel. ")
		return
	}

	var (
		firstSuffix string
		firstDx     int
		firstDy     int
	)

	for si, sheet := range sheetList {
		// sheet页不能为空
		if sheet == EmptyStr {
			fmt.Println("check failed. Excel has an empty sheet page name. ")
			return
		}

		if si == 0 {
			continue
		}

		var (
			images    = make(map[string]struct{})
			imageInEx = make(map[string]struct{})
			count     float64
		)

		// 判断sheet页是否有对应的文件夹
		sheetPath := fmt.Sprintf("%s/%s", path, sheet)
		fs, errS := ioutil.ReadDir(sheetPath)
		if errS != nil {
			fmt.Println("check failed. open ", sheetPath, " failed, please check whether the directory exists. ")
			return
		}

		for _, v := range fs {
			imageInEx[v.Name()] = struct{}{}

			suffix := getSuffix(v.Name())
			if suffix == EmptyStr {
				fmt.Println("check failed. file ", v.Name(), " has no suffix. ")
				return
			}

			if firstSuffix == EmptyStr {
				firstSuffix = suffix
			} else if firstSuffix != suffix {
				fmt.Println("check failed. The suffix of the image file must be consistent. ")
				return
			}

			filePath := fmt.Sprintf("%s/%s", sheetPath, v.Name())
			fd, errF := os.Open(filePath)
			if errF != nil {
				fmt.Println("check failed. open ", filePath, " failed, ", errF.Error())
				return
			}

			imageFile, errI := png.Decode(fd)
			if errI != nil {
				_ = fd.Close()
				fmt.Println("check failed. open ", filePath, " failed, ", errI.Error())
				return
			}

			if firstDx == 0 {
				firstDx = imageFile.Bounds().Dx()
			} else if firstDx != imageFile.Bounds().Dx() {
				_ = fd.Close()
				fmt.Println("check failed. The size of the image file must be consistent. ")
				return
			}

			if firstDy == 0 {
				firstDy = imageFile.Bounds().Dy()
			} else if firstDy != imageFile.Bounds().Dy() {
				_ = fd.Close()
				fmt.Println("check failed. The size of the image file must be consistent. ")
				return
			}
			_ = fd.Close()
		}

		sheetRows, errRI := f.GetRows(sheet)
		if errRI != nil {
			fmt.Println("check failed. ", errRI.Error())
			return
		}

		for i, row := range sheetRows {
			if i == 0 {
				continue
			}

			for vi, v := range row {
				if vi == 0 {
					images[v] = struct{}{}
				}

				if vi == 2 {
					fv, _ := strconv.ParseFloat(v, 64)
					count += fv
				}

				// 检查表格不能为空
				if vi != 2 && v == EmptyStr {
					fmt.Println("check failed. In sheet ", sheet, " Line ", i, "column ", vi, "is empty. ")
					return
				}
			}
		}

		// 检查概率
		if count != 100 {
			fmt.Println("check failed. In sheet ", sheet, " rarity is not 100. ")
			return
		}

		// 检查文件名必须一致
		if len(images) != len(imageInEx) {
			diff := diffMap(images, imageInEx)
			fmt.Println("check failed. In sheet ", sheet, " The \"file name\" is inconsistent with the file, there are ", strings.Join(diff, ","))
			return
		}

		diff := make([]string, 0)
		for k := range imageInEx {
			_, ok := images[k]
			if !ok {
				diff = append(diff, k)
			}
		}
		if len(diff) > 0 {
			fmt.Println("check failed. In sheet ", sheet, " The \"file name\" is inconsistent with the file, there are ", strings.Join(diff, ","))
			return
		}
	}

	fmt.Println("check passed. ")
	return
}

func diffMap(m1, m2 map[string]struct{}) []string {
	diff := make([]string, 0)
	if len(m1) > len(m2) {
		for k := range m1 {
			if _, ok := m2[k]; ok {
				continue
			}
			diff = append(diff, k)
		}
	}

	if len(m2) > len(m1) {
		for k := range m2 {
			if _, ok := m1[k]; ok {
				continue
			}
			diff = append(diff, k)
		}
	}

	return diff
}

func getSuffix(path string) string {
	splitArr := strings.Split(path, ".")
	if len(splitArr) == 0 {
		return EmptyStr
	}

	return splitArr[len(splitArr)-1]
}
