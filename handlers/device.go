package handlers

import (
  "fmt"
  "http"
  "models"
  "appengine"
  "appengine/datastore"
)

func init() {
  http.HandleFunc("/device", deviceHandler)
}

func deviceHandler(w http.ResponseWriter, r *http.Request) {
  context := appengine.NewContext(r)
  device := models.NewDevice()
  _, err := datastore.Put(context, datastore.NewIncompleteKey(context, "Device", nil), device)
  fmt.Fprint(w, err)
}
