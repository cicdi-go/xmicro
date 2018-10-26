package main

import (
	"crypto/tls"
	"flag"
	"github.com/cicdi-go/xmicro/services"
	"github.com/docker/libkv/store"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"strings"
	"time"
)

var (
	addr     = flag.String("addr", "0.0.0.0:8973", "server address")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName("User", new(services.User), "")
	s.Serve("tcp", *addr)
}

func addRegistryPlugin(s *server.Server) {
	etcdAddrArr := strings.Split(*etcdAddr, ",")
	cer, err := tls.LoadX509KeyPair("ssl/etcd.pem", "ssl/etcd-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    etcdAddrArr,
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
		Options: &store.Config{
			TLS: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cer},
			},
		},
	}
	err = r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
