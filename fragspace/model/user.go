package model

import (
  "strings"

  "fragspace/encryption"
)

type User struct {
  Email []byte
  EmailHash string
}

func NewUser(email string) *User {
  return &User{
    encryption.AESEncrypt(email, "user.email"),
    hash(strings.ToLower(email), encryption.ConfigKey("user.emailHash")),
  }
}
func UserEmail(email string) string {
  return hash(strings.ToLower(email), encryption.ConfigKey("user.emailHash"))
}
