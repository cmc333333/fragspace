package models

import (
  "crypto/rand"
  "appengine/datastore"
)

type User struct {
  Nickname string `json:"nickname"`
}
type Auth struct {
  User *datastore.Key
  Public []byte
  Private []byte
}
func random32() []byte {
  slice := make([]byte, 32)
  rand.Read(slice)
  return slice
}
func NewAuth(user *datastore.Key) *Auth {
  auth := &Auth {
    User: user,
    Public: random32(),
    Private: random32(),
  }
  return auth
}
