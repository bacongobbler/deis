# Riak CS
#
# VERSION       0.7.2

FROM ubuntu:14.04

# Environment variables
ENV RIAK_VERSION 2.1.1
ENV RIAK_SHORT_VERSION 2.1

RUN apt-get update && apt-get install -yq \
	wget \
	--no-install-recommends && \
	rm -rf /var/lib/apt/lists/*

# Install Riak
RUN wget http://s3.amazonaws.com/downloads.basho.com/riak/${RIAK_SHORT_VERSION}/${RIAK_VERSION}/ubuntu/trusty/riak_${RIAK_VERSION}-1_amd64.deb && \
    dpkg -i "riak_${RIAK_VERSION}-1_amd64.deb" && \
    rm -f "riak_${RIAK_VERSION}-1_amd64.deb" && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY rootfs /

# Tune Riak configuration settings for the container
RUN sed -i.bak 's/listener.http.internal = 127.0.0.1/listener.http.internal = 0.0.0.0/' /etc/riak/riak.conf && \
    sed -i.bak 's/listener.protobuf.internal = 127.0.0.1/listener.protobuf.internal = 0.0.0.0/' /etc/riak/riak.conf && \
    echo "anti_entropy.concurrency_limit = 1" >> /etc/riak/riak.conf

# Open the HTTP port for Riak
EXPOSE 8098

CMD ["/bin/boot"]
