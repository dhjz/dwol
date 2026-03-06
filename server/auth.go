package server

import (
	"crypto/md5"
	"encoding/hex"
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

func GenerateToken() string {
	muAuth.RLock()
	defer muAuth.RUnlock()
	if !passwordSet {
		return ""
	}
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func ValidateToken(token string) bool {
	muAuth.RLock()
	defer muAuth.RUnlock()
	if !passwordSet {
		return true
	}
	if token == "" {
		return false
	}
	hash := md5.Sum([]byte(password))
	return token == hex.EncodeToString(hash[:])
}
