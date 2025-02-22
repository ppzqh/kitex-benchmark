/*
 * Copyright 2021 CloudWeGo Authors
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
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/echoserver"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = 8002
)

// EchoServerImpl implements the last service interface defined in the IDL.
type EchoServerImpl struct{}

var recorder = perf.NewRecorder("KITEX-MUX@Server")

// Echo implements the EchoServerImpl interface.
func (s *EchoServerImpl) Echo(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	resp := runner.ProcessRequest(recorder, req.Action, req.Msg)

	return &echo.Response{
		Action: resp.Action,
		Msg:    resp.Msg,
	}, nil
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	address := &net.UnixAddr{Net: "tcp", Name: fmt.Sprintf(":%d", port)}
	svr := echoserver.NewServer(new(EchoServerImpl),
		server.WithServiceAddr(address),
		server.WithMuxTransport())

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
