package main

import (
    "fmt"
    "crypto/sha1"
    "encoding/base64"
)

func generate_hash(s string) string {
    bv := []byte(s)
    hasher := sha1.New()
    hasher.Write(bv)
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
    return sha
}

func main() {
     fmt.Printf(generate_hash("angryMonkey"))
}
