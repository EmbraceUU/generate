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
