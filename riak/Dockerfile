# Riak CS
#
# VERSION       0.7.2

FROM phusion/baseimage:0.9.15

# Environment variables
ENV DEBIAN_FRONTEND noninteractive
ENV RIAK_VERSION 1.4.10
ENV RIAK_SHORT_VERSION 1.4
ENV RIAK_CS_VERSION 1.5.2
ENV RIAK_CS_SHORT_VERSION 1.5
ENV STANCHION_VERSION 1.5.0
ENV STANCHION_SHORT_VERSION 1.5
ENV SERF_VERSION 0.6.3

# Install dependencies
RUN apt-get update -qq && apt-get install -y \
    unzip \
    wget \
    --no-install-recommends && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Riak
RUN wget http://s3.amazonaws.com/downloads.basho.com/riak/${RIAK_SHORT_VERSION}/${RIAK_VERSION}/ubuntu/precise/riak_${RIAK_VERSION}-1_amd64.deb && \
    dpkg -i "riak_${RIAK_VERSION}-1_amd64.deb" && \
    rm -f "riak_${RIAK_VERSION}-1_amd64.deb" && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Riak CS
RUN wget http://s3.amazonaws.com/downloads.basho.com/riak-cs/${RIAK_CS_SHORT_VERSION}/${RIAK_CS_VERSION}/ubuntu/trusty/riak-cs_${RIAK_CS_VERSION}-1_amd64.deb && \
    dpkg -i "riak-cs_${RIAK_CS_VERSION}-1_amd64.deb" && \
    rm -f "riak-cs_${RIAK_CS_VERSION}-1_amd64.deb" && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Stanchion
RUN wget http://s3.amazonaws.com/downloads.basho.com/stanchion/${STANCHION_SHORT_VERSION}/${STANCHION_VERSION}/ubuntu/trusty/stanchion_${STANCHION_VERSION}-1_amd64.deb && \
    dpkg -i "stanchion_${STANCHION_VERSION}-1_amd64.deb" && \
    rm -f "stanchion_${STANCHION_VERSION}-1_amd64.deb" && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Serf
RUN wget https://dl.bintray.com/mitchellh/serf/${SERF_VERSION}_linux_amd64.zip && \
    unzip "${SERF_VERSION}_linux_amd64.zip" -d /usr/bin/ && \
    rm -f "${SERF_VERSION}_linux_amd64.zip"

COPY rootfs /

# Setup the Serf service
RUN adduser --system --disabled-password --no-create-home \
            --quiet --force-badname --shell /bin/bash --group serf && \
    echo "serf ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/99_serf && \
    chmod 0440 /etc/sudoers.d/99_serf

# Tune Riak and Riak CS configuration settings for the container
RUN sed -i.bak "s/riak_cs-VERSION/riak_cs-${RIAK_CS_VERSION}/" /etc/riak/app.config && \
    sed -i.bak 's/\"127.0.0.1\", 8098/\"0.0.0.0\", 8098/' /etc/riak/app.config && \
    sed -i.bak "s/-env ERL_MAX_PORTS 16384/-env ERL_MAX_PORTS 64000/" /etc/riak/vm.args && \
    sed -i.bak "s/##+zdbbl 32768/+zdbbl 96000/" /etc/riak/vm.args && \
    sed -i.bak "s/{cs_ip, \"127.0.0.1\"},/{cs_ip, \"0.0.0.0\"},/" /etc/riak-cs/app.config && \
    sed -i.bak "s/{fold_objects_for_list_keys, false},/{fold_objects_for_list_keys, true},/" /etc/riak-cs/app.config && \
    sed -i.bak "s/{anonymous_user_creation, false},/{anonymous_user_creation, true},/" /etc/riak-cs/app.config && \
    sed -i.bak "s/{stanchion_ip, \"127.0.0.1\"},/{stanchion_ip, \"0.0.0.0\"},/" /etc/stanchion/app.config

# Open the HTTP port for Riak and Riak CS (S3)
EXPOSE 8098 8080

# Leverage the baseimage-docker init system
CMD ["/sbin/my_init", "--quiet"]
