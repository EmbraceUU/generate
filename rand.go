package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 随机过程
func randProcess() error {
	var err error

	blockingCount := 0
	for i := 0; i < RandNum; {
		fmt.Printf("generate image num %d \n", i)

		if blockingCount > RandNum {
			minRep--
		}

		var (
			fs          = make([]string, 0)
			fileName    string
			targets     = make([]string, 0)
			pitchRecord = make([]string, 0)
		)
		// 记录一下选中的key, 去重后要更新计数器
		for alias, j := FileNameFormatFirstChar, 0; j < len(compPriority); alias, j = alias+1, j+1 {
			currentSheet := compPriority[j]
			curSheetValue := componentMap[currentSheet]
			if len(curSheetValue) == 0 {
				continue
			}

			unDone := true
			for unDone {
				rd := RandInt(len(curSheetValue))
				key := fmt.Sprintf("%s-%d", currentSheet, rd)
				// 判断是否还有库存
				//surplus, ok := componentRep[key]
				//if !ok || surplus < minRep {
				//	continue
				//}
				ug, ok := componentUsage[key]
				if !ok || ug.Remain < minRep {
					continue
				}
				target := curSheetValue[rd].FileName

				// 记录匹配的组件
				pitchRecord = append(pitchRecord, key)
				unDone = false

				fs = append(fs, fmt.Sprintf("%s%s/%s", RootPath, currentSheet, target))
				fileName = fmt.Sprintf("%s%c%d", fileName, alias, rd+1)

				targets = append(targets, curSheetValue[rd].AttributeName)

			}
		}

		if ok := checkDuplicate(fileName); ok {
			fmt.Println("generate file again, Duplicate names appear. ", fileName)
			blockingCount++
			continue
		}

		blockingCount = 0

		// 更新计数器
		for _, v := range pitchRecord {
			//componentRep[v] = componentRep[v] - 1
			ug := componentUsage[v]
			ug.Remain--
			componentUsage[v] = ug
		}

		targets = append(targets, fileName)
		for axis, l := 'A', 0; l < len(targets); axis, l = axis+1, l+1 {
			sheetValue[fmt.Sprintf("%c%d", axis, i+2)] = targets[l]
		}

		taskQueue = append(taskQueue, GenerateImageParam{
			Fs:       fs,
			FileName: fileName,
		})

		i++
	}

	return err
}

func checkDuplicate(fileName string) bool {
	// 先替换一下fileName中的编号
	if Transfer {
		fileName = strings.Replace(fileName, "f", "a", 1)
		fileName = strings.Replace(fileName, "g", "b", 1)
		fileName = strings.Replace(fileName, "h", "c", 1)
		fileName = strings.Replace(fileName, "i", "d", 1)
		fileName = strings.Replace(fileName, "j", "e", 1)
	}

	//if _, ok := duplicateSet[fileName]; ok {
	//	return true
	//}

	splitII := strings.Split(fileName, "b")
	for k := range duplicateSet {
		splitI := strings.Split(k, "b")
		if strings.EqualFold(splitI[1], splitII[1]) {
			return true
		}
	}

	// 将fileName存入duplicateSet, 防止后面重复生成
	duplicateSet[fileName] = struct{}{}
	return false
}

func RandInt(len int) int {
	if len == 0 {
		len = 10000
	}

	// 刷新种子, 防止每次重复
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r1.Intn(len)
}
