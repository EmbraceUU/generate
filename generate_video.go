package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func GenerateVideo() error {
	const (
		CommandName = "/Users/admin/Tools/ffmpeg"
		PngPath     = "/Users/admin/Desktop/test/black_0%d.png"
		MusicPath   = "/Users/admin/Desktop/test/v1.mp3"
		OutPath     = "/Users/admin/Desktop/test/output1.mp4"
	)

	cmdArguments := []string{"-r", "1", "-i", PngPath, "-i", MusicPath, "-c:v", "libx264", "-c:a", "aac", "-b:a", "192k", "-ar", "22050", "-ac", "2", "-pix_fmt", "yuvj420p", "-shortest", "-y", OutPath}
	cmd := exec.Command(CommandName, cmdArguments...)

	fmt.Println("+++++++++++++++++++")
	fmt.Println(cmd)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	fmt.Println("+++++++++++++++++++")
	fmt.Println(out.String())
	fmt.Println("+++++++++++++++++++")
	fmt.Println(stderr.String())
	fmt.Println("+++++++++++++++++++")
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
