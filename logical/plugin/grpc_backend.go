package plugin

import (
	"math"

	"google.golang.org/grpc"
)

var defaultGRPCCallOpts []grpc.CallOption

func init() {
	defaultGRPCCallOpts = []grpc.CallOption{
		grpc.MaxCallSendMsgSize(math.MaxInt32),
		grpc.MaxCallRecvMsgSize(math.MaxInt32),
	}
}
