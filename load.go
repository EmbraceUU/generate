package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"strconv"
)

func loadDuplicateSet() error {
	data, err := ioutil.ReadFile("")
	if err != nil {
		return err
	}
	var duplicate Duplicate
	err = json.Unmarshal(data, &duplicate)
	if err != nil {
		return err
	}

	for _, v := range duplicate.RECORDS {
		duplicateSet[v.FileName] = struct{}{}
	}

	return err
}

// 加载组件文件
func loadComponent(f *excelize.File) error {
	// 读取excel文件
	// 读取sheet列表
	// 从第二个开始排序读取
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
			//componentRep[fmt.Sprintf("%s-%d", sheet, i-1)] = fr.Rarity * float64(RandNum) / float64(100)
			componentUsage[fmt.Sprintf("%s-%d", sheet, i-1)] = UsageCount{
				Available: fr.Rarity * float64(RandNum) / float64(100),
				Remain:    fr.Rarity * float64(RandNum) / float64(100),
				Rarity:    fr.Rarity,
			}
			frs = append(frs, fr)
		}
		// 将当前库存复制到快照
		//mapCopy(componentRep, componentRepSnap)
		componentMap[sheet] = frs
	}
	return nil
}

func mapCopy(m1, m2 map[string]float64) {
	for k, v := range m1 {
		m2[k] = v
	}
}
