package main

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
