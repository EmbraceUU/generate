package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func loadDuplicateSet(f *excelize.File) error {
	rows, err := f.GetRows(GenerateConfig.ExistedSheet)
	if err != nil {
		return fmt.Errorf("load existed record failed, err: %s", err.Error())
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if row == nil || len(row) == 0 {
			continue
		}

		// TODO 转换fileName也需要抽先出来, 适用于不同的模板
		fileName := transferFileName(row[0])
		duplicateSet[fileName] = struct{}{}
	}

	return err

	//data, err := ioutil.ReadFile("")
	//if err != nil {
	//	return err
	//}
	//var duplicate Duplicate
	//err = json.Unmarshal(data, &duplicate)
	//if err != nil {
	//	return err
	//}
	//
	//for _, v := range duplicate.RECORDS {
	//	v.FileName = transferFileName(v.FileName)
	//	duplicateSet[v.FileName] = struct{}{}
	//}
	//
	//return err
}

// 加载组件文件
func loadComponent(f *excelize.File) error {
	// 读取sheet列表
	sheetList := f.GetSheetList()
	if sheetList == nil || len(sheetList) == 0 {
		return fmt.Errorf("check failed. Excel is empty. ")
	}

	for si, sheet := range sheetList {
		// 从第二个开始
		if si == 0 {
			continue
		}

		sheetRows, errRI := f.GetRows(sheet)
		if errRI != nil {
			return fmt.Errorf("check failed. %s", errRI.Error())
		}

		// 读取优先级
		compPriority = append(compPriority, sheet)

		// 读取文件名和概率
		var frs []FileRarity
		for i, row := range sheetRows {
			if i == 0 {
				continue
			}

			var fr FileRarity
			for vi, v := range row {
				switch vi {
				case 0:
					fr.FileName = v
				case 1:
					fr.AttributeName = v
				case 2:
					fr.Rarity, _ = strconv.ParseFloat(v, 64)
				}
			}

			// 因为略过了0位
			componentUsage[fmt.Sprintf("%s-%d", sheet, i-1)] = UsageCount{
				Available: fr.Rarity * float64(GenerateConfig.RandNum) / float64(100),
				Remain:    fr.Rarity * float64(GenerateConfig.RandNum) / float64(100),
				Rarity:    fr.Rarity,
			}
			frs = append(frs, fr)
		}
		componentMap[sheet] = frs
	}
	return nil
}
