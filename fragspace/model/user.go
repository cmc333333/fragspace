package model

import (
  "strings"
)

type User struct {
  Email []byte
  EmailHash string
}

func NewUser(email string) *User {
  return &User{
    encrypt(email, configKey("user.email")),
    hash(strings.ToLower(email), configKey("user.emailHash")),
  }
}
func UserEmail(email string) string {
  return hash(strings.ToLower(email), configKey("user.emailHash"))
}
