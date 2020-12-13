package rpcRouter

import "context"

type RouterCaller interface {
	RouterCall(ctx context.Context, namespace string, method string, seqId int32, rpcData []byte) ([]byte, error)
}
