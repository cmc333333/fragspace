package http

import (
  "http"
  "os"
  "io"
  "io/ioutil"
  "json"
)

type JsonRequest http.Request

func (req *JsonRequest) Extract(into interface{}) os.Error {
  body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1 << 20))  //  1 MB limit
  if err != nil { return err }
  return json.Unmarshal(body, &into)
}
