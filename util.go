package main

import "math/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// StringInSlice returns whether a string is contained
// in a slice of strings
func StringInSlice(s string, list []string) bool {
    for _, m := range list {
        if m == s {
            return true
        }
    }
    return false
}

// RandomToken returns a 64 character string of random characters
func RandomToken() string {
    b := make([]byte, 64)
    for i := range b {
        b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
    }
    return string(b)
}
