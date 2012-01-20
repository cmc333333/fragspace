package rackspace

type NewServer struct {
  Server Server `json:"server"`
}
type Server struct {
  Id string `json:"id,omitempty"`
  Name string `json:"name"`
  ImageId int `json:"imageId"`
  FlavorId int `json:"flavorId"`
  HostId string `json:"hostId,omitempty"`
  Progress int `json:"progress,omitempty"`
  Status string `json:"status,omitempty"`
  AdminPass string `json:"adminPass,omitempty"`
  Metadata map[string]string `json:"metadata,omitempty"`
  Personality []Personality `json:"personality,omitempty"`
  Addresses map[string][]string `json:"addresses,omitempty"`
}
type Personality struct {
  Path string `json:"path"`
  Contents string `json:"contents"`
}
func (proxy *Proxy) NewServer(img int, flavor int) *NewServer {
  toPost := NewServer{
    Server: Server{
      Name: "Test",
      ImageId: img,
      FlavorId: flavor,
    },
  }
  response := new(NewServer)
  proxy.post("/servers", toPost, response)
  return response
}
