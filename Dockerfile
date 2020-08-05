FROM golang:1.14-stretch

ENV GOARCH=amd64

ENV GOOS=linux

ENV CC=gcc

ENV CXX=g++

ADD build.sh /usr/local/bin

RUN echo "deb http://mirrors.aliyun.com/debian stretch main non-free contrib" >>/etc/apt/sources.list \
    && echo "deb-src http://mirrors.aliyun.com/debian stretch main non-free contrib" >>/etc/apt/sources.list \
    && echo "deb http://mirrors.aliyun.com/debian stretch-updates main non-free contrib" >>/etc/apt/sources.list \
    && echo "deb-src http://mirrors.aliyun.com/debian stretch-updates main non-free contrib" >>/etc/apt/sources.list \
    && chmod +x /usr/local/bin/build.sh && apt update && apt install gcc g++ -y && apt autoclean && rm -rf /var/cache/apt

