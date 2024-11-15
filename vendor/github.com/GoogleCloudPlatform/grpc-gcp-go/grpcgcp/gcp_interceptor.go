/*
 *
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpcgcp

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

type key int

var gcpKey key

type gcpContext struct {
	// request message used for pre-process of an affinity call
	reqMsg interface{}
	// response message used for post-process of an affinity call
	replyMsg interface{}
}

// GCPUnaryClientInterceptor intercepts the execution of a unary RPC
// and injects necessary information to be used by the picker.
func GCPUnaryClientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	gcpCtx := &gcpContext{
		reqMsg:   req,
		replyMsg: reply,
	}
	ctx = context.WithValue(ctx, gcpKey, gcpCtx)

	return invoker(ctx, method, req, reply, cc, opts...)
}

// GCPStreamClientInterceptor intercepts the execution of a client streaming RPC
// and injects necessary information to be used by the picker.
func GCPStreamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	// This constructor does not create a real ClientStream,
	// it only stores all parameters and let SendMsg() to create ClientStream.
	cs := &gcpClientStream{
		ctx:      ctx,
		desc:     desc,
		cc:       cc,
		method:   method,
		streamer: streamer,
		opts:     opts,
	}
	cs.cond = sync.NewCond(cs)
	return cs, nil
}

type gcpClientStream struct {
	sync.Mutex
	grpc.ClientStream

	cond          *sync.Cond
	initStreamErr error

	ctx      context.Context
	desc     *grpc.StreamDesc
	cc       *grpc.ClientConn
	method   string
	streamer grpc.Streamer
	opts     []grpc.CallOption
}

func (cs *gcpClientStream) SendMsg(m interface{}) error {
	cs.Lock()
	// Initialize underlying ClientStream when getting the first request.
	if cs.ClientStream == nil {
		ctx := context.WithValue(cs.ctx, gcpKey, &gcpContext{reqMsg: m})
		realCS, err := cs.streamer(ctx, cs.desc, cs.cc, cs.method, cs.opts...)
		if err != nil {
			cs.initStreamErr = err
			cs.Unlock()
			cs.cond.Broadcast()
			return err
		}
		cs.ClientStream = realCS
	}
	cs.Unlock()
	cs.cond.Broadcast()
	return cs.ClientStream.SendMsg(m)
}

func (cs *gcpClientStream) RecvMsg(m interface{}) error {
	// If RecvMsg is called before SendMsg, it should wait until cs.ClientStream
	// is initialized or the initialization failed.
	cs.Lock()
	for cs.initStreamErr == nil && cs.ClientStream == nil {
		cs.cond.Wait()
	}
	if cs.initStreamErr != nil {
		cs.Unlock()
		return cs.initStreamErr
	}
	cs.Unlock()
	return cs.ClientStream.RecvMsg(m)
}
