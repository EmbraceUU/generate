package main

import (
	"github.com/go-ini/ini"
	"log"
	"runtime"
)

type GenerateSetting struct {
	OutSheet  string // 输出的sheet页名称
	SnapSheet string // 记录快照的sheet页名称
	RandNum   int    // 随机的次数
	PoolSize  int    // 协程池大小

	IfGenerateImage bool   // 生成image的开关
	RootPath        string // 文件根目录
	TempDir         string // 临时文件目录
	OutDir          string // 合成图片输出目录

	ExcelName string // excel名称
}

var GenerateConfig = &GenerateSetting{}

var cfg *ini.File

func SettingSetUp() {
	log.Println("setting setup. ")
	var err error
	cfg, err = ini.Load("app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'app.ini': %v", err)
	}

	mapTo("generate", GenerateConfig)

	if GenerateConfig.ExcelName == EmptyStr {
		GenerateConfig.ExcelName = "Attribute_info.xlsx"
	}

	if GenerateConfig.OutDir == EmptyStr {
		GenerateConfig.OutDir = "out"
	}

	if GenerateConfig.RootPath == EmptyStr {
		GenerateConfig.RootPath = "./"
	}

	if GenerateConfig.PoolSize == 0 {
		GenerateConfig.PoolSize = runtime.NumCPU()
	}

	if GenerateConfig.OutSheet == EmptyStr {
		GenerateConfig.OutSheet = "generate"
	}

	if GenerateConfig.SnapSheet == EmptyStr {
		GenerateConfig.SnapSheet = "snap"
	}
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
