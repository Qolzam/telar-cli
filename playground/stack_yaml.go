package main

// Copyright (c) Alex Ellis 2017. All rights reserved.
// This script was adapted from https://github.com/openfaas/faas-cli/blob/master/commands

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	stack "github.com/openfaas/faas-cli/stack"
	"gopkg.in/yaml.v2"
)

func readFiles(files []string, rootPath string) (map[string]string, error) {
	envs := make(map[string]string)

	for _, file := range files {
		bytesOut, readErr := ioutil.ReadFile(path.Join(rootPath, file))
		if readErr != nil {
			return nil, readErr
		}

		envFile := stack.EnvironmentFile{}
		unmarshalErr := yaml.Unmarshal(bytesOut, &envFile)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		for k, v := range envFile.Environment {
			envs[k] = v
		}

	}
	return envs, nil
}

func readConfigFile(rootPath string, file string) (map[string]string, error) {
	envs := make(map[string]string)
	bytesOut, readErr := ioutil.ReadFile(path.Join(rootPath, file))
	if readErr != nil {
		return nil, readErr
	}

	envFile := stack.EnvironmentFile{}
	unmarshalErr := yaml.Unmarshal(bytesOut, &envFile)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	for k, v := range envFile.Environment {
		envs[k] = v
	}
	return envs, nil
}

func compileEnvironment(envvarOpts []string, yamlEnvironment map[string]string, fileEnvironment map[string]string) (map[string]string, error) {
	envvarArguments, err := parseMap(envvarOpts, "env")
	if err != nil {
		return nil, fmt.Errorf("error parsing envvars: %v", err)
	}

	functionAndStack := mergeMap(yamlEnvironment, fileEnvironment)
	return mergeMap(functionAndStack, envvarArguments), nil
}

func parseMap(envvars []string, keyName string) (map[string]string, error) {
	result := make(map[string]string)
	for _, envvar := range envvars {
		s := strings.SplitN(strings.TrimSpace(envvar), "=", 2)
		if len(s) != 2 {
			return nil, fmt.Errorf("label format is not correct, needs key=value")
		}
		envvarName := s[0]
		envvarValue := s[1]

		if !(len(envvarName) > 0) {
			return nil, fmt.Errorf("empty %s name: [%s]", keyName, envvar)
		}
		if !(len(envvarValue) > 0) {
			return nil, fmt.Errorf("empty %s value: [%s]", keyName, envvar)
		}

		result[envvarName] = envvarValue
	}
	return result, nil
}

func mergeMap(i map[string]string, j map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range i {
		merged[k] = v
	}
	for k, v := range j {
		merged[k] = v
	}
	return merged
}
