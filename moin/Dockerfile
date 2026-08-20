FROM python:2.7-slim-buster AS builder

ARG MOIN_VERSION=1.9.11
ARG MOIN_SHA256=02be31d55f39d4fe0c6253df8b49e01b76d095634cbd1b56d185f66e1e0c3cf5

RUN sed -i \
      -e 's|deb.debian.org/debian|archive.debian.org/debian|g' \
      -e 's|security.debian.org/debian-security|archive.debian.org/debian-security|g' \
      /etc/apt/sources.list && \
    printf 'Acquire::Check-Valid-Until "false";\n' > /etc/apt/apt.conf.d/99archive && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates curl gcc libsasl2-dev libldap2-dev libssl-dev && \
    rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir \
      pip==20.3.4 setuptools==44.1.1 wheel==0.37.1 && \
    pip install --no-cache-dir --prefix=/install \
      python-ldap==3.3.0 uwsgi==2.0.26

RUN curl --fail --location --show-error \
      "https://static.moinmo.in/files/moin-${MOIN_VERSION}.tar.gz" \
      --output /tmp/moin.tar.gz && \
    echo "${MOIN_SHA256}  /tmp/moin.tar.gz" | sha256sum -c - && \
    mkdir /tmp/moin && \
    tar -xzf /tmp/moin.tar.gz --strip-components=1 -C /tmp/moin && \
    iconv -f iso-8859-1 /tmp/moin/setup.cfg > /tmp/moin/setup.cfg.new && \
    mv /tmp/moin/setup.cfg.new /tmp/moin/setup.cfg && \
    pip install --no-cache-dir --prefix=/install /tmp/moin

FROM python:2.7-slim-buster

RUN sed -i \
      -e 's|deb.debian.org/debian|archive.debian.org/debian|g' \
      -e 's|security.debian.org/debian-security|archive.debian.org/debian-security|g' \
      /etc/apt/sources.list && \
    printf 'Acquire::Check-Valid-Until "false";\n' > /etc/apt/apt.conf.d/99archive && \
    apt-get update && \
    apt-get install -y --no-install-recommends libldap-2.4-2 libsasl2-2 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /install/ /usr/local/

RUN mkdir /moin/
WORKDIR /moin/
CMD ["uwsgi", "--http=0.0.0.0:9000", "--wsgi-file", "moin.wsgi"]
