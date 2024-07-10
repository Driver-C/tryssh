package control

import "time"

const (
	TypeUsers                 = "users"
	TypePorts                 = "ports"
	TypePasswords             = "passwords"
	TypeCaches                = "caches"
	TypeKeys                  = "keys"
	sshClientTimeoutWhenLogin = 5 * time.Second
)
