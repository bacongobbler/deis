#!/usr/bin/env bash
set -x

ETCD_HOST=${HOST:-$(host deis-etcd | awk '{print $4}')}
ETCD_PORT=${ETCD_PORT:-4001}
ETCD="$ETCD_HOST:$ETCD_PORT"
ETCDCTL="etcdctl -C $ETCD"

if [[ "$($ETCDCTL get /deis/migrations/data/0001 2> /dev/null)" != "done" ]];
then
    for i in $($ETCDCTL ls /deis/domains 2> /dev/null);
    do
        for j in $($ETCDCTL get $i);
        do
            $ETCDCTL set /deis/domains/$j "$(basename $i)" 1> /dev/null;
            echo "migrated $j"
        done;
        $ETCDCTL rm $i;
    done
    $ETCDCTL set /deis/migrations/data/0001 "done"
fi
