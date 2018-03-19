package main

import (
	"fmt"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	// appsv1beta1 "k8s.io/api/apps/v1beta1"
)

var (
	includeFiles = []string{"deployment.yml", "service.yml"}
)

type fileName struct {
	Name string
	Path string
}

func readLocalYAML(f string) string {
	data_yml, err := ioutil.ReadFile(f)
	check(err)
	data_json, err := yaml.ToJSON(data_yml)
	check(err)
	// fmt.Printf("%v\n", string(data_json))
	if strings.Contains(f, "---") {
		return "Multiple"
	}
	return string(data_json)
}

func readRemoteJSON(ns string, d string) string {
	deployment, err := clientset.AppsV1beta1().Deployments(ns).Get(d, metav1.GetOptions{})
	check(err)
	// fmt.Printf("%v\n", deployment.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	return deployment.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]
}

func compareLocalRemote(local string, remote string) {
	objLocalRemote := diff.ObjectReflectDiff(local, remote)
	if objLocalRemote != "<no diffs>" {
		diffLocalRemote := diff.StringDiff(local, remote)
		fmt.Printf("%v\n", diffLocalRemote)
	} else {
		fmt.Printf("No diff!\n")
	}
}

func getFiles(dir string) []string {
	fileList := make([]string, 0)
	e := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if checkStringInIncludedFiles(f.Name()) {
				fileList = append(fileList, path)
			}
			// fmt.Printf("Excluding: %v\n", path)
		}
		return err
	})
	check(e)
	for _, file := range fileList {
		fmt.Println(file)
	}
	return fileList
}

func checkStringInIncludedFiles(s string) bool {
	check := false
	for _, file := range includeFiles {
		if file == s {
			check = true
			break
		}
	}
	return check
}

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func main() {
	kubeClientSetUp()
	// pods, err := clientset.CoreV1().Pods("dms-dev").List(metav1.ListOptions{})
	// check(err)
	// fmt.Printf("There are %d pods in dms-dev\n", len(pods.Items))
	// deployments, err := clientset.AppsV1beta1().Deployments("dms-dev").List(metav1.ListOptions{})
	// check(err)
	// fmt.Printf("There are %d deployments in dms-dev\n", len(deployments.Items))
	// local := readLocalYAML("deployment.yml")
	// remote := readRemoteJSON("dms-dev", "dms")
	// data_cluster, err := ioutil.ReadFile("cluster.txt")
	// data_convert, err := ioutil.ReadFile("convert.txt")
	// check(err)
	// compareLocalRemote(local, remote)
	fileList := getFiles("apps")
	for _, file := range fileList {
		data_local := readLocalYAML(file)
		if data_local != "Multiple" {
			re_ns := regexp.MustCompile("namespace\":\"([a-zA-Z0-9_-]*)\"")
			re_d := regexp.MustCompile("kind\":\"Deployment\".*name\":\"([a-zA-Z0-9_-]*)\"")
			fmt.Printf("local file: %v\n\n", data_local)
			fmt.Printf("namespace: %v\n\n", re_ns.FindStringSubmatch(data_local)[1])
			fmt.Printf("deployment: %v\n\n", re_d.FindStringSubmatch(data_local))
			// data_remote := readRemoteJSON(re_ns.FindStringSubmatch(data_local)[1], re_d.FindStringSubmatch(data_local)[1])
			// fmt.Printf("remote file: %v\n\n", data_remote)
		}
	}

}
