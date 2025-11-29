package pmongo

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lemonkingstar/spider/pkg/predis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

const (
	transactionNumberRedisKeyNamespace = "cc:v3:" + "transaction:number:"
)

type sessionKey string

func (s sessionKey) genKey() string {
	return transactionNumberRedisKeyNamespace + string(s)
}

// a transaction manager
type TxnManager struct {
	cache predis.Client
}

// InitTxnManager is to init txn manager, set the redis storage
func (t *TxnManager) InitTxnManager(r predis.Client) error {
	t.cache = r
	return nil
}

func (t *TxnManager) GetTxnNumber(sessionID string) (int64, error) {
	key := sessionKey(sessionID).genKey()
	v, err := t.cache.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(v, 10, 64)
}

// GenTxnNumber generate the transaction number from redis.
func (t *TxnManager) GenTxnNumber(sessionID string, ttl time.Duration) (int64, error) {
	// return txnNumber with 1 directly, when our mongodb client option's RetryWrite
	// is set to false.
	key := sessionKey(sessionID).genKey()

	pip := t.cache.Pipeline()
	defer pip.Close()

	ctx := context.Background()
	// we increase by step 1, so that we can calculate how many transaction has already
	// be executed in a same session.
	pip.SetNX(ctx, key, 0, ttl).Result()
	incrBy := pip.IncrBy(ctx, key, 1)
	if ttl == 0 {
		ttl = TransactionDefaultTimeout
	}
	_, err := pip.Exec(ctx)
	if err != nil {
		return 0, err
	}
	num := incrBy.Val()
	// value of (num - num/2 +1) is the real transaction number
	// in a distribute session.
	return num, nil
}

func (t *TxnManager) RemoveSessionKey(sessionID string) error {
	key := sessionKey(sessionID).genKey()
	return t.cache.Del(context.Background(), key).Err()
}

//func (t *TxnManager) ReloadSession(sess mongo.Session, info *mongo.SessionInfo) (mongo.Session, error) {
//	err := mongo.CmdbReloadSession(sess, info)
//	if err != nil {
//		return nil, err
//	}
//	return sess, nil
//}

func (t *TxnManager) PrepareCommit(cli *mongo.Client) (mongo.Session, error) {
	// create a session client.
	sess, err := cli.StartSession()
	if err != nil {
		return nil, fmt.Errorf("start session failed, err: %v", err)
	}
	return sess, nil
}

func (t *TxnManager) PrepareTransaction(cap *TxnCapable, cli *mongo.Client) (mongo.Session, error) {
	// create a session client.
	sess, err := cli.StartSession()
	if err != nil {
		return nil, fmt.Errorf("start session failed, err: %v", err)
	}

	// only for changing the transaction status
	err = sess.StartTransaction()
	if err != nil {
		return nil, fmt.Errorf("start transaction %s failed: %v", cap.SessionID, err)
	}

	//txnNumber, err := t.GenTxnNumber(cap.SessionID, cap.Timeout)
	//if err != nil {
	//	return nil, fmt.Errorf("generate txn number failed, err: %v", err)
	//}

	//// reset the session info with the session id.
	//info := &mongo.SessionInfo{
	//	TxnNubmer: txnNumber,
	//	SessionID: cap.SessionID,
	//}

	//err = mongo.CmdbReloadSession(sess, info)
	//if err != nil {
	//	return nil, fmt.Errorf("reload transaction: %s failed, err: %v", cap.SessionID, err)
	//}

	return sess, nil
}

// GetTxnContext create a session if the ctx is a session context, and the bool value is true.
// so the caller must check the bool, and use session only when the bool is true.
// otherwise the caller should not use the session, should call the mongodb command directly.
// Note: this function is always used with mongo.CmdbReleaseSession(ctx, sessCtx) to release the session connection.
func (t *TxnManager) GetTxnContext(ctx context.Context, cli *mongo.Client) (context.Context, mongo.Session, bool, error) {
	cap, useTxn, err := parseTxnInfoFromCtx(ctx)
	if err != nil {
		return ctx, nil, false, err
	}

	if !useTxn {
		// not use transaction, return directly.
		return ctx, nil, false, nil
	}

	session, err := t.PrepareTransaction(cap, cli)
	if err != nil {
		return ctx, nil, true, err
	}

	// prepare the session context, it tells the driver to run this within a transaction.
	//sessCtx := mongo.CmdbContextWithSession(ctx, session)
	//
	//return sessCtx, session, true, nil

	return nil, session, false, nil
}

// parseTxnInfoFromCtx try to parse transaction info from context,
// it returns the TxnCable, and a bool to indicate whether it's a transaction context or not.
// so the caller can use the returned TxnCapable only when the bool is true. otherwise it will be panic.
func parseTxnInfoFromCtx(txnCtx context.Context) (*TxnCapable, bool, error) {
	id := txnCtx.Value(TransactionIdHeader)
	if id == nil {
		// do not use transaction, and return directly.
		return nil, false, nil
	}

	txnID, ok := id.(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid transaction id value： %v", id)
	}

	// parse timeout
	ttl := txnCtx.Value(TransactionTimeoutHeader)
	if ttl == nil {
		return nil, false, errors.New("transaction timeout value not exist")
	}

	ttlStr, ok := ttl.(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid transaction timeout value: %v", ttl)
	}

	timeout, err := strconv.ParseInt(ttlStr, 10, 64)
	if err != nil {
		return nil, false, fmt.Errorf("invalid transaction timeout value, parse %v failed, err: %v", ttl, err)
	}

	cap := &TxnCapable{
		// timeout is not
		Timeout:   time.Duration(timeout),
		SessionID: txnID,
	}
	return cap, true, nil
}

func (t *TxnManager) AutoRunWithTxn(ctx context.Context, cli *mongo.Client, cmd func(ctx context.Context) error) error {
	cap, useTxn, err := parseTxnInfoFromCtx(ctx)
	if err != nil {
		return err
	}

	if !useTxn {
		// not use transaction, run command directly.
		return cmd(ctx)
	}

	_, err = t.PrepareTransaction(cap, cli)
	if err != nil {
		return err
	}

	// prepare the session context, it tells the driver to run this within a transaction.
	//sessCtx := mongo.CmdbContextWithSession(ctx, session)

	// run the command and check error
	//err = cmd(sessCtx)
	//if err != nil {
	//	// release the session connection.
	//	// Attention: do not use session.EndSession() to do this, it will abort the transaction.
	//	// mongo.CmdbReleaseSession(ctx, session)
	//	return err
	//}
	// release the session connection.
	// Attention: do not use session.EndSession() to do this, it will abort the transaction.
	// mongo.CmdbReleaseSession(ctx, session)
	return nil
}

func GenSessionID() (string, error) {
	// mongodb driver used this as it's mongodb session id, and we use it too.
	id, err := uuid.New()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(id[:]), nil
}

// generate a session id and set it to header.
func GenTxnCableAndSetHeader(header http.Header, opts ...TxnOption) (*TxnCapable, error) {
	sessionID, err := GenSessionID()
	if err != nil {
		return nil, fmt.Errorf("generate session id failed, err: %v", err)
	}
	var timeout time.Duration
	if len(opts) != 0 {
		if opts[0].Timeout < 30*time.Second {
			timeout = TransactionDefaultTimeout
		} else {
			timeout = opts[0].Timeout
		}
	} else {
		// set default value
		timeout = TransactionDefaultTimeout
	}

	header.Set(TransactionIdHeader, sessionID)
	header.Set(TransactionTimeoutHeader, strconv.FormatInt(int64(timeout), 10))

	cap := TxnCapable{
		Timeout:   timeout,
		SessionID: sessionID,
	}
	return &cap, nil
}
