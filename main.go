package main

import (
	"fmt"
	"github.com/fatih/color"
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
	includeFiles                  = []string{"deployment.yml", "service.yml"}
	re_namespace                  = regexp.MustCompile("namespace\":\"([a-zA-Z0-9_-]*)\"")
	re_search_name                = regexp.MustCompile("name\":\"([a-zA-Z0-9_-]*)\"")
	re_remove_service_annotations = regexp.MustCompile("\"annotations\":{},")
)

func checkDeployments(d string) {
	if len(re_search_name.FindStringSubmatch(d)) > 0 {
		// fmt.Printf("namespace: %v\n", re_namespace.FindStringSubmatch(d)[1])
		// fmt.Printf("deployment: %v\n\n", re_deployment.FindStringSubmatch(d)[1])
		// fmt.Printf("local file: %q\n\n", d)
		data_remote := readRemoteDeploymentJSON(re_namespace.FindStringSubmatch(d)[1], re_search_name.FindStringSubmatch(d)[1])
		if data_remote != "" {
			// fmt.Printf("remote file: %q\n\n", data_remote)
			compareLocalRemote(d, data_remote)
		}
	}
}

func checkServices(d string) {
	fmt.Printf("namespace: %v\n", re_namespace.FindStringSubmatch(d)[1])
	fmt.Printf("service: %v\n\n", re_search_name.FindStringSubmatch(d)[1])
	fmt.Printf("local file: %q\n\n", d)
	data_remote := readRemoteServiceJSON(re_namespace.FindStringSubmatch(d)[1], re_search_name.FindStringSubmatch(d)[1])
	if data_remote != "" {
		fmt.Printf("remote file: %q\n\n", data_remote)
		compareLocalRemote(d, data_remote)
	}
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

func readRemoteDeploymentJSON(ns string, d string) string {
	deployment, err := clientset.AppsV1beta1().Deployments(ns).Get(d, metav1.GetOptions{})
	// check(err)
	if err != nil {
		color.Red("Deployment does not exist: %v!\n\n", d)
		return ""
		// fmt.Printf("Deployment does not exist!\n\n")
	}
	// fmt.Printf("%v\n", deployment.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	return strings.TrimSuffix(deployment.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"], "\n")
}

func readRemoteServiceJSON(ns string, s string) string {
	service, err := clientset.CoreV1().Services(ns).Get(s, metav1.GetOptions{})
	if err != nil {
		color.Red("Service does not exist: %v!\n\n", s)
		return ""
	}
	// fmt.Printf("%+v", service.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	return strings.TrimSuffix(re_remove_service_annotations.ReplaceAllString(service.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"], ""), "\n")
}

func compareLocalRemote(local string, remote string) {
	objLocalRemote := diff.ObjectReflectDiff(local, remote)
	if objLocalRemote != "<no diffs>" {
		color.Yellow("Differences found in resource: %v!\n\n", remote)
		diffLocalRemote := diff.StringDiff(local, remote)
		fmt.Printf("%v\n\n", diffLocalRemote)
	} else {
		color.Blue("No diff!\n\n")
	}
}

func getFiles(dir string) []string {
	fileList := make([]string, 0)
	e := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if checkStringInIncludedFiles(f.Name()) {
				fileList = append(fileList, path)
			}
		}
		return err
	})
	check(e)
	// for _, file := range fileList {
	// 	fmt.Println(file)
	// }
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
	argDir := os.Args[1]
	kubeClientSetUp()
	fileList := getFiles(argDir)
	for _, file := range fileList {
		data_local := readLocalYAML(file)
		if data_local != "Multiple" {
			if strings.Contains(data_local, "Deployment") {
				checkDeployments(data_local)
			} else if strings.Contains(data_local, "Service") {
				checkServices(data_local)
			}
		}
	}
}
