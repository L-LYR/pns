FROM golang:1.17
WORKDIR /pns
ENV GOPROXY=https://proxy.golang.com.cn,direct
COPY .. .
RUN go mod tidy
RUN make all
EXPOSE 10086 10087 10088
VOLUME [ "/pns/log" ]
CMD [ "/pns/build/pns" ]