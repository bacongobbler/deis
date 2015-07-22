FROM debian:jessie

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && \
    apt-get install -yq \
    bash \
    ca-certificates \
    curl \
    dnsutils \
    sudo \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/get-deis/etcdctl-v0.4.9 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.10.0/confd-0.10.0-linux-amd64 \
    && chmod +x /usr/local/bin/confd

# define execution environment
CMD ["/app/bin/boot"]
EXPOSE 8000

# define work environment
WORKDIR /app

# add a deis user
RUN adduser \
    --disabled-password \
    --home /app \
    --shell /bin/bash \
    --gecos "" \
	deis

# create a /app directory for storing application data
RUN mkdir -p /app && chown -R deis:deis /app

# create directory for confd templates
RUN mkdir -p /templates && chown -R deis:deis /templates

ADD requirements.txt /app/requirements.txt

# install required system packages
# HACK: install git so we can install bacongobbler's fork of django-fsm
RUN apt-get update && \
    apt-get install -yq \
    build-essential \
    git \
    libffi-dev \
    libpq-dev \
    libldap2-dev \
    libsasl2-dev \
    libssl-dev \
    python \
    python-dev \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
# install pip \
    && curl -sSL https://raw.githubusercontent.com/pypa/pip/7.0.3/contrib/get-pip.py | python - \
# install dependencies \
    && pip install --disable-pip-version-check --no-cache-dir -r /app/requirements.txt

ADD . /app

# Create static resources
RUN /app/manage.py collectstatic --settings=deis.settings --noinput

ENV DEIS_RELEASE 1.9.0-dev
