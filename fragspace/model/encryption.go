package model

import (
  "bytes"
  "gob"
  "crypto/hmac"
  "encoding/base64"

  "libs/passwordhash"
)

func passwordHash(valueStr string) []byte {
  hashBuffer := bytes.NewBuffer(make([]byte, 0, 1000000))
  gob.NewEncoder(hashBuffer).Encode(passwordhash.New(valueStr))
  return hashBuffer.Bytes()
}
func decodePassword(b []byte) *passwordhash.PasswordHash {
  hash := new(passwordhash.PasswordHash)
  err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(hash)
  if err == nil {
    return hash
  }
  panic(err)
}
func hash(valueStr string, key []byte) string {
  hasher := hmac.NewSHA256(key)
  hasher.Write([]byte(valueStr))
  return base64.StdEncoding.EncodeToString(hasher.Sum())
}
