package rackspace

type NewServer struct {
  Server Server
}
type Server struct {
  Name string
  ImageId int
  FlavorId int
  Metadata map[string]string
  Personality []Personality
}
type Personality struct {
  Path string
  Contents string
}
