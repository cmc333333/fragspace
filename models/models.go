package models

import (
  "rand"
  "time"
)

func init() {
  rand.Seed(time.Nanoseconds())
}

type User struct {
  Public []byte
  Private []byte
  Nickname string
}

type Device struct {
  Public []byte
  Private []byte
  Nickname string
}
func random32() []byte {
  slice := make([]byte, 32)
  for i := 0; i < 8; i++ {
    randomValue := rand.Int()
    for j := 0; j < 4; j++ {
      slice[i*4 + j] = byte(randomValue >> uint(j*8))
    }
  }
  return slice
}
func NewDevice() *Device {
  device := &Device { 
    Public: random32(),
    Private: random32(),
  }
  return device
}
