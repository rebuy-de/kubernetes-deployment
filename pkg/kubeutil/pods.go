package kubeutil

import (
	"fmt"

	"k8s.io/client-go/pkg/api/v1"
)

const (
	ErrImagePull = "ErrImagePull"
	Error        = "Error"
)

func PodWarnings(pod *v1.Pod) error {
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

func containerWarnings(status v1.ContainerStatus) error {
	if status.State.Waiting != nil && status.State.Waiting.Reason == ErrImagePull {
		return fmt.Errorf("failed to pull docker image")
	}

	if status.State.Terminated != nil && status.State.Terminated.Reason == Error {
		return fmt.Errorf(
			"failed to start container (Container: %v, ExitCode: %v, Restarts: %v)",
			status.Name, status.State.Terminated.ExitCode, status.RestartCount)
	}

	return nil
}
