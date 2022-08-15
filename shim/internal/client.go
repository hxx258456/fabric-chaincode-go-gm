// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"time"

	tls "github.com/hxx258456/ccgo/gmtls"

	"github.com/hxx258456/ccgo/grpc"
	"github.com/hxx258456/ccgo/grpc/credentials"
	"github.com/hxx258456/ccgo/grpc/keepalive"
	peerpb "github.com/hxx258456/fabric-protos-go-gm/peer"
)

const (
	dialTimeout        = 10 * time.Second
	maxRecvMessageSize = 100 * 1024 * 1024 // 100 MiB
	maxSendMessageSize = 100 * 1024 * 1024 // 100 MiB
)

// NewClientConn ...
func NewClientConn(
	address string,
	tlsConf *tls.Config,
	kaOpts keepalive.ClientParameters,
) (*grpc.ClientConn, error) {

	dialOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(kaOpts),
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMessageSize),
			grpc.MaxCallSendMsgSize(maxSendMessageSize),
		),
	}

	if tlsConf != nil {
		creds := credentials.NewTLS(tlsConf)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, address, dialOpts...)
}

// NewRegisterClient ...
func NewRegisterClient(conn *grpc.ClientConn) (peerpb.ChaincodeSupport_RegisterClient, error) {
	return peerpb.NewChaincodeSupportClient(conn).Register(context.Background())
}
