package encryption

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
)

func AESEncrypt(valueStr string, configKey string) []byte {
  return AESEncryptKey(valueStr, ConfigKey(configKey))
}
func AESEncryptKey(valueStr string, key []byte) []byte {
  iv := make([]byte, aes.BlockSize)
  rand.Read(iv)
  aesCipher, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }
  ctr := cipher.NewCTR(aesCipher, iv)
  value := []byte(valueStr)
  length := len(value)
  toReturn := make([]byte, length + aes.BlockSize)  //  add room for the iv
  copy(toReturn, iv)
  ctr.XORKeyStream(toReturn[aes.BlockSize:], value)
  return toReturn
}

func AESDecrypt(value []byte, configKey string) string {
  return AESDecryptKey(value, ConfigKey(configKey))
}
func AESDecryptKey(value []byte, key []byte) string {
  iv := value[:aes.BlockSize]
  aesCipher, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }
  ctr := cipher.NewCTR(aesCipher, iv)
  value = value[aes.BlockSize:]
  toReturn := make([]byte, len(value))
  ctr.XORKeyStream(toReturn, value)
  return string(toReturn)
}
