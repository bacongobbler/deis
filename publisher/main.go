package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"

	"github.com/deis/deis/publisher/server"
)

var (
	bindAddr        string
	refreshDuration time.Duration
	etcdTTL         time.Duration
	host            string
	dockerAddr      string
	etcdAddr        string
	logLevel        string
)

func init() {
	flag.StringVar(&bindAddr, "bind-addr", "localhost:6060", "The IP address and port to bind to")
	flag.DurationVar(&refreshDuration, "refresh-duration", 10 * time.Second, "The time to wait between etcd refreshes")
	flag.DurationVar(&etcdTTL, "etcd-ttl", refreshDuration * 2, "The TTL for all of the keys in etcd")
	flag.StringVar(&host, "host", "127.0.0.1", "The host where the publisher is running")
	flag.StringVar(&dockerAddr, "docker-addr", "unix:///var/run/docker.sock", "The bind address to the docker API")
	flag.StringVar(&etcdAddr, "etcd-addr", "http://127.0.0.1:4001", "The bind address to the etcd host")
	flag.StringVar(&logLevel, "log-level", "error", "Acceptable values: error, debug")
}

func main() {
	flag.Parse()

	log.Println("booting publisher...")

	dockerClient, err := docker.NewClient(dockerAddr)
	if err != nil {
		log.Fatal(err)
	}
	etcdClient := etcd.NewClient([]string{etcdAddr})

	server := server.New(dockerClient, etcdClient, host, logLevel)

	go server.Listen(etcdTTL)

	go func() {
		msg := fmt.Sprintf("profiler listening on %s", bindAddr)
		if err := http.ListenAndServe(bindAddr, nil); err != nil {
			msg = err.Error()
		}
		log.Println(msg)
	}()

	for {
		go server.Poll(etcdTTL)
		time.Sleep(refreshDuration)
	}
}
