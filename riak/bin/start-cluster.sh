#! /bin/bash

set -e

if env | egrep -q "DOCKER_RIAK_CS_DEBUG"; then
  set -x
fi

if ! env | egrep -q "DOCKER_RIAK_CS_AUTOMATIC_CLUSTERING=1" && \
  env | egrep -q "DOCKER_RIAK_CS_HAPROXY=1"; then
  echo
  echo "It appears that you have enabled HAProxy support, but have"
  echo "not enabled automatic clustering. In order to use Riak and"
  echo "HAProxy, please enable automatic clustering."
  echo

  exit 1
fi

CLEAN_DOCKER_HOST=$(echo "${DOCKER_HOST}" | cut -d'/' -f3 | cut -d':' -f1)
CLEAN_DOCKER_HOST=${CLEAN_DOCKER_HOST:-localhost}
DOCKER_RIAK_CS_CLUSTER_SIZE=${DOCKER_RIAK_CS_CLUSTER_SIZE:-5}

if docker ps -a | egrep "deis/riak" >/dev/null; then
  echo
  echo "It looks like you already have some Riak containers running."
  echo "Please take them down before attempting to bring up another"
  echo "cluster with the following command:"
  echo
  echo "  make stop-cluster"
  echo

  exit 1
fi

echo
echo "Bringing up cluster nodes:"
echo

for index in $(seq -w "1" "99");
do
  if [ "${index}" -gt "${DOCKER_RIAK_CS_CLUSTER_SIZE}" ] ; then
      break
  fi
  if [ "${index}" -gt "1" ] ; then
    docker run -e "DOCKER_RIAK_CS_CLUSTER_SIZE=${DOCKER_RIAK_CS_CLUSTER_SIZE}" \
               -e "DOCKER_RIAK_CS_AUTOMATIC_CLUSTERING=${DOCKER_RIAK_CS_AUTOMATIC_CLUSTERING}" \
               -P --name "riak-cs${index}" --link "riak-cs01:seed" \
               -d deis/riak-cs > /dev/null 2>&1
  else
    docker run -e "DOCKER_RIAK_CS_CLUSTER_SIZE=${DOCKER_RIAK_CS_CLUSTER_SIZE}" \
               -e "DOCKER_RIAK_CS_AUTOMATIC_CLUSTERING=${DOCKER_RIAK_CS_AUTOMATIC_CLUSTERING}" \
               -P --name "riak-cs${index}" -d deis/riak-cs > /dev/null 2>&1
  fi

  CONTAINER_ID=$(docker ps | egrep "riak-cs${index}" | cut -d" " -f1)
  CONTAINER_PORT=$(docker port "${CONTAINER_ID}" 8080 | cut -d ":" -f2)

  until curl -s "http://${CLEAN_DOCKER_HOST}:${CONTAINER_PORT}/riak-cs/ping" | egrep "OK" > /dev/null 2>&1;
  do
    sleep 3
  done

  echo "  Successfully brought up [riak-cs${index}]"
done

echo
echo "  Riak CS credentials:"
echo

for field in admin_key admin_secret ; do
  echo -n "    ${field}: "

  docker exec riak-cs01 egrep "${field}" /etc/riak-cs/app.config | cut -d'"' -f2
done

echo
echo "Please wait approximately 30 seconds for the cluster to stabilize."
echo
