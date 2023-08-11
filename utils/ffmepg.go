package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	workDir, _ := os.Getwd()
	snapshotPath = workDir + "/assets/" + snapshotPath
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		ZapLogger.Errorf("ffmpegInput err: %v", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		ZapLogger.Errorf("imagingDecode err: %v", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		ZapLogger.Errorf("imagingSave: %v", err)
		return "", err
	}

	return snapshotName, nil
}
