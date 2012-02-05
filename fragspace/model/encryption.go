package model

import (
  "bytes"
  "gob"
  "crypto/aes"
  "crypto/cipher"
  "crypto/hmac"
  "crypto/rand"
  "encoding/base64"
  "encoding/hex"

  "libs/conf"
  "libs/passwordhash"
)

func encrypt(valueStr string, key []byte) []byte {
  iv := make([]byte, aes.BlockSize)
  rand.Read(iv)
  aesCipher, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }
  ctr := cipher.NewCTR(aesCipher, iv)
  value := []byte(valueStr)
  length := len(value)
  //  make sure length is padded to the block size
  if divisor, remainder := length / aes.BlockSize, length % aes.BlockSize; remainder != 0 {
    length = aes.BlockSize * (divisor + 1)
  }
  toReturn := make([]byte, length + aes.BlockSize)  //  add room for the iv
  copy(toReturn, iv)
  ctr.XORKeyStream(toReturn[aes.BlockSize:], value)
  return toReturn
}
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
func configKey(name string) []byte {
  config, err := conf.ReadConfigFile("config.ini")
  if err != nil {
    panic(err)
  }
  keyStr, err := config.GetString("encryption", name)
  if err != nil {
    panic(err)
  }
  key, err := hex.DecodeString(keyStr)
  if err != nil {
    panic(err)
  }
  return key
}
