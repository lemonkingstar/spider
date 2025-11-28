package pmongo

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lemonkingstar/spider/pkg/predis"
)

// DB db operation interface
type DB interface {
	// Table collection 操作
	Table(collection string) Table

	//// NextSequence 获取新序列号(非事务)
	//NextSequence(ctx context.Context, sequenceName string) (uint64, error)
	//
	////NextSequences 批量获取新序列号(非事务)
	//NextSequences(ctx context.Context, sequenceName string, num int) ([]uint64, error)

	// Ping 健康检查
	Ping() error // 健康检查

	// HasTable 判断是否存在集合
	HasTable(ctx context.Context, name string) (bool, error)
	// DropTable 移除集合
	DropTable(ctx context.Context, name string) error
	// CreateTable 创建集合
	CreateTable(ctx context.Context, name string) error

	IsDuplicatedError(error) bool
	IsNotFoundError(error) bool

	Close() error

	// CommitTransaction 提交事务
	CommitTransaction(context.Context, *TxnCapable) error
	// AbortTransaction 取消事务
	AbortTransaction(context.Context, *TxnCapable) error

	// InitTxnManager TxnID management of initial transaction
	InitTxnManager(r predis.Client) error
}

type ReadPreferenceMode string

const (
	// NilMode not set
	NilMode ReadPreferenceMode = ""
	// PrimaryMode indicates that only a primary is
	// considered for reading. This is the default
	// mode.
	PrimaryMode ReadPreferenceMode = "1"
	// PrimaryPreferredMode indicates that if a primary
	// is available, use it; otherwise, eligible
	// secondaries will be considered.
	PrimaryPreferredMode ReadPreferenceMode = "2"
	// SecondaryMode indicates that only secondaries
	// should be considered.
	SecondaryMode ReadPreferenceMode = "3"
	// SecondaryPreferredMode indicates that only secondaries
	// should be considered when one is available. If none
	// are available, then a primary will be considered.
	SecondaryPreferredMode ReadPreferenceMode = "4"
	// NearestMode indicates that all primaries and secondaries
	// will be considered.
	NearestMode ReadPreferenceMode = "5"
)

const (
	// if maxOpenConns isn't configured, use default value
	DefaultMaxOpenConns = 1000
	// if maxOpenConns exceeds maximum value, use maximum value
	MaximumMaxOpenConns = 3000
	// if maxIDleConns is less than minimum value, use minimum value
	MinimumMaxIdleOpenConns = 50
	// if timeout isn't configured, use default value
	DefaultSocketTimeout = 10
	// if timeout exceeds maximum value, use maximum value
	MaximumSocketTimeout = 30
	// if timeout less than the minimum value, use minimum value
	MinimumSocketTimeout = 5
)

// Config config
type Config struct {
	Connect       string
	Address       string
	User          string
	Password      string
	Port          string
	Database      string
	Mechanism     string
	MaxOpenConns  uint64
	MaxIdleConns  uint64
	RsName        string
	SocketTimeout int
}

// BuildURI return mongo uri according to  https://docs.mongodb.com/manual/reference/connection-string/
// format example: mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
func (c Config) BuildURI() string {
	if c.Connect != "" {
		return c.Connect
	}

	if !strings.Contains(c.Address, ":") && len(c.Port) > 0 {
		c.Address = c.Address + ":" + c.Port
	}

	c.User = url.QueryEscape(c.User)
	c.Password = url.QueryEscape(c.Password)
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authMechanism=%s", c.User, c.Password, c.Address, c.Database, c.Mechanism)
	return uri
}

func (c Config) GetMongoConf() MongoConf {
	return MongoConf{
		MaxOpenConns:  c.MaxOpenConns,
		MaxIdleConns:  c.MaxIdleConns,
		URI:           c.BuildURI(),
		RsName:        c.RsName,
		SocketTimeout: c.SocketTimeout,
	}
}

func (c Config) GetMongoClient() (db DB, err error) {
	mongoConf := MongoConf{
		MaxOpenConns:  c.MaxOpenConns,
		MaxIdleConns:  c.MaxIdleConns,
		URI:           c.BuildURI(),
		RsName:        c.RsName,
		SocketTimeout: c.SocketTimeout,
	}
	db, err = NewMgo(mongoConf, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	return
}
