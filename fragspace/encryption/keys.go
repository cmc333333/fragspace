package encryption

import (
  "encoding/hex"

  "libs/conf"
)

func ConfigKey(name string) []byte {
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
