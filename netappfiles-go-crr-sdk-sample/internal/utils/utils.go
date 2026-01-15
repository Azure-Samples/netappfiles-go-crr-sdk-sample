// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// Package that provides some general functions.

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Azure-Samples/netappfiles-go-crr-sdk-sample/netappfiles-go-crr-sdk-sample/internal/models"
	"golang.org/x/term"
)

// PrintHeader prints a header message
func PrintHeader(header string) {
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))
}

// ConsoleOutput writes to stdout.
func ConsoleOutput(message string) {
	log.Println(message)
}

// Contains checks if there is a string already in an existing splice of strings
func Contains(array []string, element string) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}

// GetBytesInTiB converts a value from bytes to tebibytes (TiB)
func GetBytesInTiB(size uint64) uint32 {
	return uint32(size / 1024 / 1024 / 1024 / 1024)
}

// GetTiBInBytes converts a value from tebibytes (TiB) to bytes
func GetTiBInBytes(size uint32) uint64 {
	return uint64(size * 1024 * 1024 * 1024 * 1024)
}

// ReadAzureBasicInfoJSON reads the Azure Authentication json file json file and unmarshals it.
func ReadAzureBasicInfoJSON(path string) (*models.AzureBasicInfo, error) {
	infoJSON, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("failed to read file: %v", err)
		return &models.AzureBasicInfo{}, err
	}
	var info models.AzureBasicInfo
	json.Unmarshal(infoJSON, &info)
	return &info, nil
}

// FindInSlice returns index greater than -1 and true if item is found
// Code from https://golangcode.com/check-if-element-exists-in-slice/
func FindInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// GetPassword gets a password
func GetPassword(prompt string) string {
	fmt.Print(prompt)
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return strings.TrimSpace(string(bytePassword))
}

// GetSubscriptionIdFromEnv gets the Azure Subscription ID from environment variable
func GetSubscriptionIdFromEnv() (string, error) {
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		return "", fmt.Errorf("AZURE_SUBSCRIPTION_ID environment variable is not set")
	}
	return subscriptionID, nil
}

// GetSubscriptionIdFromAzCli gets the Azure Subscription ID using the Azure CLI
func GetSubscriptionIdFromAzCli() (string, error) {
	cmd := exec.Command("az", "account", "show", "--query", "id", "-o", "tsv")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running az account show: %v\n", err)
		return "", err
	}
	subscriptionID := strings.TrimSpace(string(output))
	return subscriptionID, nil
}

func GetSubscriptionId() (string, error) {
	var errs []string

	subscriptionID, err := GetSubscriptionIdFromEnv()
	if err == nil {
		return subscriptionID, nil
	}
	errs = append(errs, fmt.Sprintf("GetSubscriptionIdFromEnv error: %v", err))

	subscriptionID, err = GetSubscriptionIdFromAzCli()
	if err == nil {
		return subscriptionID, nil
	}
	errs = append(errs, fmt.Sprintf("GetSubscriptionIdFromAzCli error: %v", err))

	return "", fmt.Errorf(strings.Join(errs, "; "))
}
