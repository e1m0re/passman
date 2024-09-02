package grpc

type key int

const (
	UserIDMarker key = iota
)

var AnonymousMethods = map[string]bool{
	"/protobuf.auth.service.AuthService/Login":  true,
	"/protobuf.auth.service.AuthService/SignUp": true,
}
