FROM thub.autohome.com.cn/project/dba-centos7-python3.6.6:1.0.0
MAINTAINER autohome <yanchangjie@autohome.com.cn>

# install system requirement
RUN mkdir -p /data/sysdir/{servers,applogs,pids,data,bak} && \
    mkdir -p /data/sysdir/servers/autodts && \
    pip install gevent
COPY ./ /data/sysdir/servers/autodts
WORKDIR /data/sysdir/servers/autodts
# ENV FRDBENV=test
EXPOSE 80
ENTRYPOINT ["python3", "/data/sysdir/servers/autodts/autodts/run_wsgi.py"]
