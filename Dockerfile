# Copyright 2015 bradley childs, All rights reserved.
#

FROM centos:7
MAINTAINER bradley childs, bchilds@gmail.com
RUN yum update -y ; yum install automake fuse-devel gcc-c++ git libcurl-devel libxml2-devel make openssl-devel -y
RUN mkdir -p /root
WORKDIR /root
RUN git clone https://github.com/s3fs-fuse/s3fs-fuse.git
WORKDIR /root/s3fs-fuse
RUN ./autogen.sh  
RUN ./configure 
RUN make 
RUN make install
RUN mkdir -p /mnt/mountpoint
RUN echo $'#!/bin/sh\n \
           echo $S3User:$S3Secret > /etc/passwd-s3fs\n \
           chmod 600 /etc/passwd-s3fs\n \
           exec s3fs "$@"' > /root/s3fs.sh 
RUN chmod +x /root/s3fs.sh

ENV S3User $S3User
ENV S3Secret $S3Secret

ENTRYPOINT ["/root/s3fs.sh"]
# CMD ["--help"]
