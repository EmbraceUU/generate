package main

import (
	"bytes"
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

	Infoln("+++++++++++++++++++")
	Infoln(cmd)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	Infoln("+++++++++++++++++++")
	Infoln(out.String())
	Infoln("+++++++++++++++++++")
	Infoln(stderr.String())
	Infoln("+++++++++++++++++++")
	if err != nil {
		Infoln(err.Error())
	}

	return nil
}
