package taskrunner

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/pili-video-server/scheduler/dbops"
	"github.com/xfrr/goffmpeg/transcoder"
)

func deleteVideo(vid string) error {
	// log.Printf("begin to dele")
	pathVd, _ := filepath.Abs(VIDEOS_PATH + vid)
	pathVd2, _ := filepath.Abs(VIDEOS_PATH2 + vid)
	pathVd3, _ := filepath.Abs(VIDEOS_PATH3 + vid)
	pathVd4, _ := filepath.Abs(VIDEOS_PATH4 + vid)
	pathIc, _ := filepath.Abs(ICON_PATH + vid)

	err := os.Remove(pathVd)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("delete video_id: %v error!\n", err)
		return err
	}

	err = os.Remove(pathIc)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("delete video_id: %v error!\n", err)
		return err
	}

	err = os.Remove(pathVd2)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("delete video_id2: %v error!\n", err)
		return err
	}
	err = os.Remove(pathVd3)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("delete video_id3: %v error!\n", err)
		return err
	}
	err = os.Remove(pathVd4)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("delete video_id4: %v error!\n", err)
		return err
	}

	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadDeletion(5)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		log.Printf("clear all tasks ware done!\n")
		return errors.New("all tasks ware done!")
	}

	for _, vid := range res {
		log.Printf("clear dispatcher tasks : %v\n", vid)
		dc <- vid
	}

	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error
forloop:
	for {
		select {
		case id := <-dc:
			go func(vid interface{}) {
				err := deleteVideo(vid.(string))
				if err != nil {
					errMap.Store(id, err)
					return
				}

				err = dbops.DeleteDeletion(vid.(string))
				if err != nil {
					errMap.Store(id, err)
					return
				}
			}(id)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			log.Printf("vid: %v, err: %v", k, v)
			return false
		}
		return true
	})
	return nil
}

func formatVideo(vid string) error {
	// log.Printf("begin to format video...\n")
	err := formatVideoPart(vid, "720p")
	if err != nil {
		log.Printf("format error :%v\n", err)
		return err
	}

	err = formatVideoPart(vid, "480p")
	if err != nil {
		log.Printf("format error :%v\n", err)
		return err
	}

	err = formatVideoPart(vid, "360p")
	if err != nil {
		log.Printf("format error :%v\n", err)
		return err
	}

	return nil
}

func formatVideoPart(vid string, format string) error {
	inputPath, _ := filepath.Abs(VIDEOS_PATH + vid)
	var outputPath, resolution string
	switch format {
	case "720p":
		outputPath, _ = filepath.Abs(VIDEOS_PATH2 + vid + VIDEOS_FORMAT)
		resolution = "1280x720"
	case "480p":
		outputPath, _ = filepath.Abs(VIDEOS_PATH3 + vid + VIDEOS_FORMAT)
		resolution = "848x480"
	case "360p":
		outputPath, _ = filepath.Abs(VIDEOS_PATH4 + vid + VIDEOS_FORMAT)
		resolution = "640x360"
	}

	trans := new(transcoder.Transcoder)
	err := trans.Initialize(inputPath, outputPath)
	if err != nil {
		log.Printf("format error :%v\n", err)
		return err
	}

	trans.MediaFile().SetResolution(resolution)
	trans.MediaFile().SetPreset("ultrafast")

	done := trans.Run(false)
	err = <-done

	return nil
}

func VideoFormatDispatcher(dc dataChan) error {
	res, err := dbops.ReadUnFormat(5)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		log.Printf("format all tasks ware done!\n")
		return errors.New("all tasks ware done!")
	}

	for _, vid := range res {
		log.Printf("format dispatcher tasks : %v\n", vid)
		dc <- vid
	}

	return nil
}

func VideoFormatExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error
forloop:
	for {
		select {
		case id := <-dc:
			func(vid interface{}) {
				err := formatVideo(vid.(string))
				if err != nil {
					errMap.Store(id, err)
					return
				}

				err = dbops.DeleteUnFormat(vid.(string))
				if err != nil {
					errMap.Store(id, err)
					return
				}
			}(id)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			log.Printf("vid: %v, err: %v", k, v)
			return false
		}
		return true
	})
	return nil
}
