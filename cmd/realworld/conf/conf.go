package conf

import (
	"github.com/lemonkingstar/spider/pkg/pconf"
)

var (
	bs Bootstrap
)

type Bootstrap struct {
	Server *Server `json:"server,omitempty"`
	Data   *Data   `json:"data,omitempty"`
}

type Server struct {
	Http *ServerHTTP `json:"http,omitempty"`
	Grpc *ServerGRPC `json:"grpc,omitempty"`
}

type ServerHTTP struct {
	Name    string `json:"name,omitempty"`
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
}

type ServerGRPC struct {
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
}

type Data struct {
	Database *DataDatabase `json:"database,omitempty"`
	Redis    *DataRedis    `json:"redis,omitempty"`
}

type DataDatabase struct {
	Driver string `json:"driver,omitempty"`
	Source string `json:"source,omitempty"`
}

type DataRedis struct {
	Network  string `json:"network,omitempty"`
	Addr     string `json:"addr,omitempty"`
	Password string `json:"password,omitempty"`
	Db       int32  `json:"db,omitempty"`
}

func Startup() error {
	return pconf.Unmarshal(&bs)
}

func GetServer() *Server { return bs.Server }
func GetData() *Data     { return bs.Data }
