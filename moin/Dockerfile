FROM python:2.7

RUN apt-get update && \
    apt-get install -y gcc libsasl2-dev libldap2-dev libssl-dev

RUN pip install --upgrade pip setuptools && \
    pip install python-ldap uwsgi

ADD http://static.moinmo.in/files/moin-1.9.11.tar.gz /tmp/

RUN cd /tmp && \
    tar xf moin-1.9.8.tar.gz && \
    rm moin-1.9.8.tar.gz && \
    cd moin-1.9.8 && \
    iconv -f iso-8859-1 setup.cfg  > setup.cfg.new && \
    mv setup.cfg.new setup.cfg && \
    pip install .

RUN mkdir /moin/
WORKDIR /moin/
CMD uwsgi --http=0.0.0.0:9000 --wsgi-file moin.wsgi
