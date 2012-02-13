package api

import (
  "http"

  "fragspace/encryption"
  fhttp "fragspace/http"
  "fragspace/httphelpers"
  "fragspace/model"
)

func init() {
  m := new(meHandler)
  m.Self = m
  http.Handle("/api/me", m)
}

type meHandler struct { fhttp.BaseHandler }

type getMeResponse struct {
  Email string `json:"email"`
}

func (handler *meHandler) Get(r *http.Request) fhttp.Response {
  return httphelpers.ReqUser(r, func(user *model.User) fhttp.Response {
    return fhttp.JsonResponse{&getMeResponse{encryption.AESDecrypt(user.Email, "user.email")}}
  })
}
