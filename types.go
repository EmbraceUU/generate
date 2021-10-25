package main

type Duplicate struct {
	RECORDS []struct {
		FileName string `json:"file_name"`
	} `json:"RECORDS"`
}

type GenerateImageParam struct {
	Fs       []string
	FileName string
}

type FileRarity struct {
	FileName      string
	AttributeName string
	Rarity        float64
}

type UsageCount struct {
	Available  float64 // 提供数量
	Remain     float64 // 库存
	Usage      float64 // 使用数量
	UsageRatio float64 // 使用占比
	Rarity     float64 // 目标占比
}
