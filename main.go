package main

import "flag"

func init() {
	SettingSetUp()
}

func main() {
	var mode string
	flag.StringVar(&mode, "mode", Generate, "the mode type, default is generate. ")
	flag.Parse()

	switch mode {
	case Generate:
		GenerateRun()
	case Check:
		CheckRun()
	}
}
