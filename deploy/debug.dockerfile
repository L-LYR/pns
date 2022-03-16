FROM golang:1.17
WORKDIR /pns
ENV GOPROXY=https://proxy.golang.com.cn,direct
COPY .. .
RUN make all
EXPOSE 10086
CMD ["/pns/build/pns"]