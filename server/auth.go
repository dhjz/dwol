package server

import (
	"sync"
)

var (
	password    string
	passwordSet bool
	muAuth      sync.RWMutex
)

func SetPassword(pwd string) {
	muAuth.Lock()
	defer muAuth.Unlock()
	password = pwd
	passwordSet = pwd != ""
}

func CheckPassword(pwd string) bool {
	muAuth.RLock()
	defer muAuth.RUnlock()
	if !passwordSet {
		return true
	}
	return pwd == password
}

func IsPasswordRequired() bool {
	muAuth.RLock()
	defer muAuth.RUnlock()
	return passwordSet
}
