package kubernetes

import (
	"testing"
	"github.com/rebuy-de/kubernetes-deployment/util"
	"net"
	"log"
	"bufio"
	"os"
	"os/user"
	"strings"
	"time"
)

const local_k8r_config = "/.kube/config"
const local_deployment_test_file = "./test-fixtures/deployment.yaml"

func checkIfPortIsOpen(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", host+":"+string(port), time.Duration(2)*time.Second)
	if err != nil {
		log.Println("Connection error:", err)
		return false
	} else {
		conn.Close()
		return true
	}
}

func getIpAndPortFromK8sFile(filepath string)(host string, port string){
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "server:"){
			data := strings.Split(strings.Split(scanner.Text(),"//")[1],":")
			host = data[0]
			port = data[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return host,port
}

func getHomeFolder() string{
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	return  usr.HomeDir
}

func TestKubernetes_New_NoConfig(t *testing.T) {
	//git, cleanupParent := testCreateTempGitDir(t)
	//defer cleanupParent()
	_, e := New("")
	util.AssertNoError(t, e)
}

func TestKubernetes_New_WithConfig(t *testing.T) {
	// Test should still pass we are providing only config file string (no evaluation at this point)
	_, e := New("/some/none/existing/file.yml")
	util.AssertNoError(t, e)
}

func TestKubernetes_Apply_NoFile(t *testing.T) {
	host,port := getIpAndPortFromK8sFile(getHomeFolder()+ local_k8r_config)
	if checkIfPortIsOpen(host,port) {
		k, _ := New("/some/none/existing/file.yml")
		util.AssertHasError(t, k.Apply(""))
	}
}

func TestKubernetes_Apply_FileDoseNotExists(t *testing.T) {
	host,port := getIpAndPortFromK8sFile(getHomeFolder()+local_k8r_config)
	if checkIfPortIsOpen(host,port) {
		k, _ := New("/some/none/existing/file.yml")
		util.AssertHasError(t, k.Apply("/some/none/existing/file.yml"))
	}
}

func TestKubernetes_Apply_FilExists(t *testing.T) {
	host,port := getIpAndPortFromK8sFile(getHomeFolder()+local_k8r_config)
	if checkIfPortIsOpen(host,port) {
		k, _ := New("")
		util.AssertNoError(t, k.Apply(local_deployment_test_file))
	}
}

func TestKubernetes_Get_FilExists(t *testing.T) {
	host,port := getIpAndPortFromK8sFile(getHomeFolder()+local_k8r_config)
	if checkIfPortIsOpen(host,port) {
		k, _ := New("")
		out, err := k.Get(local_deployment_test_file)
		util.AssertNoError(t, err)
		util.AssertStringIsNotEmpty(t, string(out[:]), "Testing json response")
		util.AssertStringContains(t, string(out[:]), "apiVersion", "Testing json response")
	}
}

