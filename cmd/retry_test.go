package cmd

import (
	"fmt"
	"testing"
	"time"
)

type FailN struct {
	n int
}

func (f *FailN) Task() error {
	if f.n > 0 {
		f.n -= 1
		return fmt.Errorf("%d tries left.", f.n+1)
	}
	return nil
}

func TestRetry(t *testing.T) {
	for i := 0; i < 3; i++ {
		failn := &FailN{i}
		err := Retry(3, time.Second/10, failn.Task)
		if err != nil {
			t.Errorf("Task %d failed with: %v", i, err)
		}
	}

	for i := 3; i < 6; i++ {
		failn := &FailN{i}
		err := Retry(3, time.Second/10, failn.Task)
		if err == nil {
			t.Errorf("Task %d should have get an error.", i)
		}
	}

}
