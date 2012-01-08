package handlers

import (
  "fmt"
  "http"
  "models"
  "lib"
  "appengine"
  "appengine/datastore"
)

func init() {
  u := new(userHandler)
  u.Self = u
  http.Handle("/user", u)
}

type userHandler struct { lib.BaseHandler }

type postResponse struct {
  Public string `json:"public"`
  Private string `json:"private"`
}

func (handler *userHandler) Post(w lib.JsonResponse, r *lib.JsonRequest) {
  user := new(models.User)
  err := r.Extract(user)
  if err != nil {
    panic(lib.UserError{"invalid json"})
  }
  if user.Nickname == "" {
    panic(lib.UserError{"nickname cannot be empty"})
  }

  context := appengine.NewContext((*http.Request)(r))
  userKey, err := datastore.Put(context, datastore.NewIncompleteKey(context, "User", nil), user)
  if err != nil {
    panic(lib.ServerError{err.String()})
  }

  auth := models.NewAuth(userKey)
  _, err = datastore.Put(context, datastore.NewIncompleteKey(context, "Auth", nil), auth)
  if err != nil {
    panic(lib.ServerError{err.String()})
  }
  w.Success(&postResponse{
    fmt.Sprintf("%x", auth.Public),
    fmt.Sprintf("%x", auth.Private),
  })
}

/*
func userHandler(w http.ResponseWriter, r *http.Request) {
  //context := appengine.NewContext(r)
  user := models.NewUser()
  //_, err := datastore.Put(context, datastore.NewIncompleteKey(context, "Device", nil), device)
  fmt.Fprintf(w, "%x", user.Public)
}
*/
