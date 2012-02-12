package oauth2

import (
  "crypto/rand"
  "encoding/hex"
  "url"

  "fragspace/encryption"
)

type Client struct {
  Id string
  Secret []byte
  Redirect []byte
  Name []byte
  Email []byte
}

func NewClient(redirect string, name string, email string) *Client {
  idBuff, secretBuff := make([]byte, 32), make([]byte, 128)
  if _, err := rand.Read(idBuff); err != nil {
    panic(err)
  }
  if _, err := rand.Read(secretBuff); err != nil {
    panic(err)
  }

  return &Client{
    hex.EncodeToString(idBuff),
    secretBuff,
    encryption.AESEncrypt(redirect, "oauthclient.redirect"),
    encryption.AESEncrypt(name, "oauthclient.name"),
    encryption.AESEncrypt(email, "oauthclient.email"),
  }
}

func (client *Client) redirectUrl(code string) string {
  //  We assume that the string is of the proper format; safe assumption since we reject invalid uris when saving
  baseUrl, _ := url.Parse(encryption.AESDecrypt(client.Redirect, "oauthclient.redirect"))
  query := baseUrl.Query()
  query.Add("code", code)
  baseUrl.RawQuery = query.Encode()
  return baseUrl.String()
}
