package kubernetes

import (
	"testing"
	"github.com/rebuy-de/kubernetes-deployment/util"
	"strings"
)

const local_k8r_config = "/.kube/config"

func TestKubernetes_New_NoConfig(t *testing.T) {
	_, e := New("")
	util.AssertNoError(t, e)
}

func TestKubernetes_New_WithConfig(t *testing.T) {
	// Test should still pass we are providing only config file string (no evaluation at this point)
	_, e := New("/some/none/existing/file.yml")
	util.AssertNoError(t, e)
}

func TestKubernetes_Apply_WithoutConfig(t *testing.T) {
	k, _ := New("")
	k.KubectlPath = "/bin/echo"
	out, error := k.Apply("/some/none/existing/file.yml")
	util.AssertNoError(t, error)
	util.AssertStringEquals(t, "apply -f /some/none/existing/file.yml --record", strings.TrimSpace(string(out[:])), "")
}

func TestKubernetes_Apply_WithConfig(t *testing.T) {
	k, _ := New(local_k8r_config)
	k.KubectlPath = "/bin/echo"
	out, error := k.Apply("/some/none/existing/file.yml")
	util.AssertNoError(t, error)
	util.AssertStringEquals(t, "--kubeconfig=/.kube/config apply -f /some/none/existing/file.yml --record", strings.TrimSpace(string(out[:])), "")
}

func TestKubernetes_Get_WithoutConfig(t *testing.T) {
	k, _ := New("")
	k.KubectlPath = "/bin/echo"
	out, error := k.Get("/some/none/existing/file.yml")
	util.AssertNoError(t, error)
	util.AssertStringEquals(t, "get -f /some/none/existing/file.yml -o json", strings.TrimSpace(string(out[:])), "")
}
func TestKubernetes_Get_WithConfig(t *testing.T) {
	k, _ := New(local_k8r_config)
	k.KubectlPath = "/bin/echo"
	out, error := k.Get("/some/none/existing/file.yml")
	util.AssertNoError(t, error)
	util.AssertStringEquals(t, "--kubeconfig=/.kube/config get -f /some/none/existing/file.yml -o json", strings.TrimSpace(string(out[:])), "")
}


