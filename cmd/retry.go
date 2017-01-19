package cmd

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Retryer func() error

func Retry(count int, wait time.Duration, task Retryer) error {
	var err error

	for i := 0; i < count; i++ {
		err = task()
		if err == nil {
			return nil
		}
		left := count - i - 1
		if left > 0 {
			log.Warnf("Task failed. %d retries left. Retrying in %v.", left, wait)
			time.Sleep(wait)
		} else {
			log.Error("Task failed. No more retries left.")
		}
	}

	return fmt.Errorf("retry failed %d times: %v", count, err)
}
