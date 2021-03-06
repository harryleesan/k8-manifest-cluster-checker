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
)

var (
	// file names of manifests used for comparison
	includeFiles = []string{"deployment", "service", "ingress"}
	// regex used to search for namespace
	re_namespace = regexp.MustCompile("namespace\":\"([a-zA-Z0-9_-]*)\"")
	// regex used to search for name of resource
	re_search_name = regexp.MustCompile("name\":\"([a-zA-Z0-9_-]*)\"")
	// regex used to identify empty annotations
	re_remove_annotations = regexp.MustCompile("\"annotations\":{},")
)

// checkIngresses() compares the local ingress resources with the remote
// ingress resources.
func checkIngresses(d string) {
	if len(re_search_name.FindStringSubmatch(d)) > 0 {
		// fmt.Printf("namespace: %v\n", re_namespace.FindStringSubmatch(d)[1])
		// fmt.Printf("deployment: %v\n\n", re_deployment.FindStringSubmatch(d)[1])
		// fmt.Printf("local file: %q\n\n", d)
		name := re_search_name.FindStringSubmatch(d)[1]
		data_remote := readRemoteIngressJSON(re_namespace.FindStringSubmatch(d)[1], name)
		if data_remote != "" {
			// fmt.Printf("remote file: %q\n\n", data_remote)
			color.Green("Checking ingress: %v ...\n\n", name)
			compareLocalRemote(strings.Join([]string{name}, " (ingress)"), d, data_remote)
		}
	}
}

// checkDeployments() compares the local deployment resources with the remote
// deployment resources.
func checkDeployments(d string) {
	if len(re_search_name.FindStringSubmatch(d)) > 0 {
		// fmt.Printf("namespace: %v\n", re_namespace.FindStringSubmatch(d)[1])
		// fmt.Printf("deployment: %v\n\n", re_deployment.FindStringSubmatch(d)[1])
		// fmt.Printf("local file: %q\n\n", d)
		name := re_search_name.FindStringSubmatch(d)[1]
		data_remote := readRemoteDeploymentJSON(re_namespace.FindStringSubmatch(d)[1], name)
		if data_remote != "" {
			// fmt.Printf("remote file: %q\n\n", data_remote)
			color.Green("Checking deployment: %v ...\n\n", name)
			compareLocalRemote(strings.Join([]string{name}, " (deployment)"), d, data_remote)
		}
	}
}

// checkServices() compares the local service resources with the remote service
// resources.
func checkServices(d string) {
	// fmt.Printf("namespace: %v\n", re_namespace.FindStringSubmatch(d)[1])
	// fmt.Printf("service: %v\n\n", re_search_name.FindStringSubmatch(d)[1])
	// fmt.Printf("local file: %q\n\n", d)
	name := re_search_name.FindStringSubmatch(d)[1]
	data_remote := readRemoteServiceJSON(re_namespace.FindStringSubmatch(d)[1], name)
	if data_remote != "" {
		// fmt.Printf("remote file: %q\n\n", data_remote)
		color.Green("Checking service: %v ...\n\n", name)
		compareLocalRemote(strings.Join([]string{name}, " (service)"), d, data_remote)
	}
}

// readLocalYAML() reads the local YAML file specified by f and returns the
// content of the file as a string.
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

// readRemoteIngressJSON() reads the remote ingress resource as a string.
func readRemoteIngressJSON(ns string, d string) string {
	ingress, err := clientset.ExtensionsV1beta1().Ingresses(ns).Get(d, metav1.GetOptions{})
	// check(err)
	if err != nil {
		color.Red("Ingress does not exist: %v!\n\n", d)
		return ""
	}
	// fmt.Printf("%v\n", ingress.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	return strings.TrimSuffix(re_remove_annotations.ReplaceAllString(ingress.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"], ""), "\n")
}

// readRemoteDeploymentJSON() reads the remote deployment resource as a string.
func readRemoteDeploymentJSON(ns string, d string) string {
	deployment, err := clientset.AppsV1beta1().Deployments(ns).Get(d, metav1.GetOptions{})
	// check(err)
	if err != nil {
		color.Red("Deployment does not exist: %v!\n\n", d)
		return ""
	}
	// fmt.Printf("%v\n", deployment.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	return strings.TrimSuffix(re_remove_annotations.ReplaceAllString(deployment.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"], ""), "\n")
}

// readRemoteServiceJSON() reads the remote deployment resource as a string.
func readRemoteServiceJSON(ns string, s string) string {
	service, err := clientset.CoreV1().Services(ns).Get(s, metav1.GetOptions{})
	if err != nil {
		color.Red("Service does not exist: %v!\n\n", s)
		return ""
	}
	// fmt.Printf("%+v", service.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	return strings.TrimSuffix(re_remove_annotations.ReplaceAllString(service.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"], ""), "\n")
}

// compareLocalRemote() compares the local manifest with remote manifest.
func compareLocalRemote(t string, local string, remote string) {
	objLocalRemote := diff.ObjectReflectDiff(local, remote)
	if objLocalRemote != "<no diffs>" {
		color.Yellow("Differences found in resource: %v\n\n", remote)
		diffLocalRemote := diff.StringDiff(local, remote)
		fmt.Printf("%v\n\n", diffLocalRemote)
	} else {
		color.Green("No diff in %v\n\n", t)
	}
}

// getFiles() reads in all the file names in the specified directory dir.
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

// checkStringInIncludedFiles() compares the file s with includedFiles.
func checkStringInIncludedFiles(s string) bool {
	check := false
	for _, file := range includeFiles {
		if strings.Contains(s, file) {
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

// main() is the entry point.
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
			} else if strings.Contains(data_local, "Ingress") {
				checkIngresses(data_local)
			}
		}
	}
}
