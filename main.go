/*

Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package main

import (
	"flag"
	"strings"

	"github.com/golang/glog"
	"github.com/google/go-microservice-helpers/server"
	"github.com/google/go-microservice-helpers/tracing"

	pb "github.com/google/pubkeystore/api"
	"github.com/google/pubkeystore/server"
)

type staticKeysFlag []string

func (sk *staticKeysFlag) String() string {
	if sk == nil {
		return ""
	}
	return strings.Join(*sk, ",")
}

func (sk *staticKeysFlag) Set(value string) error {
	*sk = append(*sk, value)
	return nil
}

var (
	staticKeys staticKeysFlag
)

func main() {
	flag.Var(&staticKeys, "static-key", "openssh key to announce")
	flag.Parse()
	defer glog.Flush()

	err := tracing.InitTracer(*serverhelpers.ListenAddress, "keyserver")
	if err != nil {
		glog.Fatalf("failed to init tracing interface: %v", err)
	}

	grpcServer, credstoreClient, err := serverhelpers.NewServer()
	if err != nil {
		glog.Fatalf("failed to init GRPC server: %v", err)
	}
	if credstoreClient == nil {
		glog.Fatalf("failed to init credstore")
	}

	svr := server.NewServer(staticKeys)

	pb.RegisterPubkeyStoreServer(grpcServer, &svr)

	err = serverhelpers.ListenAndServe(grpcServer, nil)
	if err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}
