package model

import (
  "os"
  "strings"

  "appengine"
  "appengine/datastore"

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
func UserFromKey(keyStr string, context appengine.Context) (*User, os.Error) {
  key, err := datastore.DecodeKey(keyStr)
  if err != nil {
    return nil, err
  }
  user := new(User)
  if err := datastore.Get(context, key, user); err != nil {
    return nil, err
  }
  return user, nil
}
