//
// common.go
// Copyright (C) 2021 rmelo <Ricardo Melo <rmelo@ludia.com>>
//
// Distributed under terms of the MIT license.
//

package etcd

import (
	uuid "github.com/satori/go.uuid"
)

//uuidGenerator return random uuid that are intended to be used as unique identifiers.
//In a string format
func uuidGenerator() string {
	// Creating UUID Version 4
	uu := uuid.NewV4()

	return uu.String()
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
