package pmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/lemonkingstar/spider/pkg/iserver"
	"github.com/lemonkingstar/spider/pkg/pnet"
	"github.com/lemonkingstar/spider/pkg/predis"
	"github.com/lemonkingstar/spider/pkg/psafe"
)

var (
	masterNodeKey   = "%s:node:master"
	electionNodeKey = "%s:node:election"
	lockKey         = "%s:node:lock"
	// 节点标识
	nodeIdentity = ""

	client predis.Client
)

// ID get the node identity
func ID() string { return nodeIdentity }

// Get check is master node
func Get() bool {
	nid, err := client.Get(context.Background(), masterNodeKey).Result()
	if err != nil {
		logger.Errorf("get master error: %s", err.Error())
		return false
	}
	return nid == nodeIdentity
}

func Unregister() {
	client.Del(context.Background(), fmt.Sprintf("%s:%s", electionNodeKey, nodeIdentity))
}

func Run(prefix string) {
	if prefix == "" {
		prefix = iserver.SpiderApp
	}
	masterNodeKey = fmt.Sprintf(masterNodeKey, prefix)
	electionNodeKey = fmt.Sprintf(electionNodeKey, prefix)
	lockKey = fmt.Sprintf(lockKey, prefix)
	nodeIdentity = buildSnowflakeUUID()
	if client = predis.Default(); client == nil {
		logger.Fatalf("default client not found, exit.")
		return
	}
	register()

	psafe.Go(func() {
		registerTicker := time.NewTicker(20 * time.Second)
		electionTicker := time.NewTicker(9 * time.Second)
		defer func() {
			registerTicker.Stop()
			electionTicker.Stop()
		}()

		for {
			select {
			case <-registerTicker.C:
				register()
			case <-electionTicker.C:
				setMaster()
			}
		}
	})
}

// register current node info
func register() {
	ip, _ := pnet.GetInternalIP()
	hostname, _ := os.Hostname()
	hb := time.Now().UnixNano() / 1e6

	ctx := context.Background()
	value := map[string]interface{}{
		"ip": ip, "hostname": hostname, "hb": hb, "node": nodeIdentity,
	}
	d, _ := json.Marshal(value)
	nodeKey := fmt.Sprintf("%s:%s", electionNodeKey, nodeIdentity)
	if err := client.Set(ctx, nodeKey, string(d), 33*time.Second).Err(); err != nil {
		logger.Errorf("register node info error: %s, %s", err.Error(), string(d))
	}
}

func setMaster() {
	ctx := context.Background()
	if b, _ := client.Lock(ctx, lockKey); !b {
		return
	}
	defer client.Unlock(ctx, lockKey)

	election := func() {
		// election master
		nodes := NodeList()
		if len(nodes) == 0 {
			logger.Infof("election nodes is empty")
			return
		}
		sort.Strings(nodes)
		// 直接更新master节点
		if err := client.Set(ctx, masterNodeKey, nodes[0], time.Hour).Err(); err != nil {
			logger.Errorf("set master error: %s, %s", err.Error(), nodes[0])
		}
	}
	election()
	time.Sleep(900 * time.Millisecond)
}

func Exec(f func()) {
	pc, _, _, _ := runtime.Caller(1)
	n := runtime.FuncForPC(pc).Name()
	if !Get() {
		logger.Infof("normal node, stop exec. [%s]", n)
		return
	}
	logger.Infof("master node, exec. [%s]", n)
	f()
}

// StrictExec 严格获取锁执行
// 针对某些不能多个节点同时执行的场景
func StrictExec(key string, f func()) {
	pc, _, _, _ := runtime.Caller(1)
	n := runtime.FuncForPC(pc).Name()
	if !Get() {
		logger.Infof("normal node, stop exec. [%s]", n)
		return
	}
	logger.Infof("master node, exec. [%s]", n)
	ctx := context.Background()
	if b, _ := client.Lock(ctx, key); !b {
		return
	}
	defer client.Unlock(ctx, key)
	f()
}

// NodeList get all node list
func NodeList() []string {
	nodes := make([]string, 0)
	retry := 0
	var keys []string
	var err error
	var cursor uint64 = 0
	pattern := fmt.Sprintf("%s:", electionNodeKey)
	for {
		retry++
		keys, cursor, err = client.Scan(context.Background(), cursor, pattern+"*", 1000).Result()
		if err != nil {
			logger.Errorf("scan node info error: %s", err.Error())
			return nodes
		}
		if len(keys) != 0 || cursor == 0 || retry > 99 {
			break
		}
	}
	for _, key := range keys {
		node := strings.TrimPrefix(key, pattern)
		nodes = append(nodes, node)
	}
	return nodes
}

// NodeExist check node exist
func NodeExist(node string) bool {
	key := fmt.Sprintf("%s:%s", electionNodeKey, node)
	if err := client.Get(context.Background(), key).Err(); err != nil {
		return false
	}
	return true
}
