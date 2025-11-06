package pmongo

import (
	"context"
	"fmt"
)

// CommitTransaction 提交事务
func (c *Client) CommitTransaction(ctx context.Context, cap *TxnCapable) error {
	rid := ctx.Value(ContextRequestIDField)
	reloadSession, err := c.tm.PrepareTransaction(cap, c.dbc)
	if err != nil {
		return err
	}
	// reset the transaction state, so that we can commit the transaction after start the
	// transaction immediately.
	//mongo.CmdbPrepareCommitOrAbort(reloadSession)

	// we commit the transaction with the session id
	err = reloadSession.CommitTransaction(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %s failed, err: %v, rid: %v", cap.SessionID, err, rid)
	}

	err = c.tm.RemoveSessionKey(cap.SessionID)
	if err != nil {
		// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
		// do not return.
	}

	return nil
}

// AbortTransaction 取消事务
func (c *Client) AbortTransaction(ctx context.Context, cap *TxnCapable) error {
	rid := ctx.Value(ContextRequestIDField)
	reloadSession, err := c.tm.PrepareTransaction(cap, c.dbc)
	if err != nil {
		return err
	}
	// reset the transaction state, so that we can abort the transaction after start the
	// transaction immediately.
	//mongo.CmdbPrepareCommitOrAbort(reloadSession)

	// we abort the transaction with the session id
	err = reloadSession.AbortTransaction(ctx)
	if err != nil {
		return fmt.Errorf("abort transaction: %s failed, err: %v, rid: %v", cap.SessionID, err, rid)
	}

	err = c.tm.RemoveSessionKey(cap.SessionID)
	if err != nil {
		// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
		// do not return.
	}

	return nil
}
