// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// Sample package that is used to obtain an authorizer token
// and to return unmarshall the Azure authentication file
// created by az ad sp create create-for-rbac command-line
// into an AzureAuthInfo object.

package iam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Azure-Samples/netappfiles-go-crr-sdk-sample/netappfiles-go-crr-sdk-sample/internal/models"
	"github.com/Azure-Samples/netappfiles-go-crr-sdk-sample/netappfiles-go-crr-sdk-sample/internal/utils"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// GetAuthorizer gets an authorization token to be used within ANF client
func GetAuthorizer() (azcore.TokenCredential, string, error) {

	authorizer, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("%v", err))
		return nil, "", err
	}

	subscriptionID, err := utils.GetSubscriptionId()
	if err != nil {
		return nil, "", err
	}

	return authorizer, subscriptionID, nil
}

// readAuthJSON reads the Azure Authentication json file json file and unmarshals it.
func readAuthJSON(path string) (*models.AzureAuthInfo, error) {
	infoJSON, err := ioutil.ReadFile(path)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("failed to read file: %v", err))
		return &models.AzureAuthInfo{}, err
	}
	var authInfo models.AzureAuthInfo
	json.Unmarshal(infoJSON, &authInfo)
	return &authInfo, nil
}
