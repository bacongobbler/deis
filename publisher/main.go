package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"

	"github.com/deis/deis/publisher/server"
)

var (
	bindAddr   string
	interval   time.Duration
	etcdTTL    time.Duration
	host       string
	dockerAddr string
	etcdAddr   string
	log        = logrus.New()
	logLevel   string
)

func init() {
	flag.StringVar(&bindAddr, "bind-addr", "localhost:6060", "address to listen for incoming HTTP requests")
	flag.DurationVar(&interval, "interval", 10*time.Second, "backend polling interval")
	flag.DurationVar(&etcdTTL, "publish-ttl", interval*2, "backend TTL when publishing keys")
	flag.StringVar(&host, "host", "127.0.0.1", "host address of the machine")
	flag.StringVar(&dockerAddr, "docker-addr", "unix:///var/run/docker.sock", "address to a docker API")
	flag.StringVar(&etcdAddr, "etcd-addr", "http://127.0.0.1:4001", "address to the etcd host")
	flag.StringVar(&logLevel, "log-level", "error", "level which publisher should log messages")
}

func main() {
	flag.Parse()

	log.Info("booting publisher...")
	setLogLevel()

	// run a profiler
	go func() {
		if err := http.ListenAndServe(bindAddr, nil); err != nil {
			log.Warningf("failed to run profiler (%s)", err)
		} else {
			log.Infof("profiler listening on %s", bindAddr)
		}
	}()

	dockerClient, err := docker.NewClient(dockerAddr)
	if err != nil {
		log.Fatal(err)
	}
	etcdClient := etcd.NewClient([]string{etcdAddr})

	server := server.New(dockerClient, etcdClient, host)

	go server.Listen(etcdTTL)

	for {
		go server.Poll(etcdTTL)
		time.Sleep(interval)
	}
}

func setLogLevel() {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("failed to parse log level '%s' (%s)", logLevel, err)
	}
	logrus.SetLevel(lvl)
}
