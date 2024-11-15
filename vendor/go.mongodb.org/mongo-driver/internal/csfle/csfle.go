// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package csfle

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	EncryptedCacheCollection      = "ecc"
	EncryptedStateCollection      = "esc"
	EncryptedCompactionCollection = "ecoc"
)

// GetEncryptedStateCollectionName returns the encrypted state collection name associated with dataCollectionName.
func GetEncryptedStateCollectionName(efBSON bsoncore.Document, dataCollectionName string, stateCollection string) (string, error) {
	fieldName := stateCollection + "Collection"
	val, err := efBSON.LookupErr(fieldName)
	if err != nil {
		if !errors.Is(err, bsoncore.ErrElementNotFound) {
			return "", err
		}
		// Return default name.
		defaultName := "enxcol_." + dataCollectionName + "." + stateCollection
		return defaultName, nil
	}

	stateCollectionName, ok := val.StringValueOK()
	if !ok {
		return "", fmt.Errorf("expected string for '%v', got: %v", fieldName, val.Type)
	}
	return stateCollectionName, nil
}
