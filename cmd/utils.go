package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func LookupStringEnv(envName string, defVal string) string {
	if envVal, exists := os.LookupEnv(envName); exists {
		return envVal
	}

	return defVal
}

func LookupIntEnv(envName string, defVal int) int {
	if envVal, exists := os.LookupEnv(envName); exists {
		if intVal, err := strconv.Atoi(envVal); err == nil {
			return intVal
		}
	}

	return defVal
}

func getCurrentNamespace() string {
	if ns, ok := os.LookupEnv("POD_NAMESPACE"); ok {
		return ns
	}

	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}

	return "default"
}
