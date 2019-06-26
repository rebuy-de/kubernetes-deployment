package kubeutil

import (
	"fmt"

	core_v1 "k8s.io/api/core/v1"
)

const (
	ErrImagePullReason = "ErrImagePull"
	ErrorReason        = "Error"
)

func PodWarnings(pod *core_v1.Pod) error {
	if pod.ObjectMeta.DeletionTimestamp != nil {
		// ignore Pods with pending deletion
		return nil
	}

	for _, cs := range pod.Status.InitContainerStatuses {
		err := containerWarnings(cs)
		if err != nil {
			return err
		}
	}

	for _, cs := range pod.Status.ContainerStatuses {
		err := containerWarnings(cs)
		if err != nil {
			return err
		}
	}

	return nil
}

type ErrImagePull struct{}

func (err ErrImagePull) Error() string {
	return "failed to pull docker image"
}

type ErrCrash struct {
	Name         string
	ExitCode     int32
	RestartCount int32
}

func (err ErrCrash) Error() string {
	return fmt.Sprintf(
		"failed to start container (Container: %v, ExitCode: %v, Restarts: %v)",
		err.Name, err.ExitCode, err.RestartCount)
}

func containerWarnings(status core_v1.ContainerStatus) error {
	if status.State.Waiting != nil && status.State.Waiting.Reason == ErrImagePullReason {
		return ErrImagePull{}
	}

	if status.State.Terminated != nil && status.State.Terminated.Reason == ErrorReason {
		return ErrCrash{
			Name:         status.Name,
			ExitCode:     status.State.Terminated.ExitCode,
			RestartCount: status.RestartCount,
		}
	}

	return nil
}
