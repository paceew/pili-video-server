package taskrunner

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/pace/sample/scheduler/dbops"
)

func deleteVideo(vid string) error {
	path, _ := filepath.Abs(VIDEOS_PATH + vid)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("delete video_id: %v error!\n", err)
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
		return errors.New("all tasks ware done!")
	}

	for _, vid := range res {
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
