package iserver

import "fmt"

var _idx = identification{
	Name: "unknown", IP: "0.0.0.0",
}

type identification struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Pid  int    `json:"pid"`
	UUID string `json:"uuid"`
}

func GetIdentification() string {
	return _idx.Name
}

func GetServerInstance() string {
	return fmt.Sprintf("%s:%d", _idx.IP, _idx.Port)
}
