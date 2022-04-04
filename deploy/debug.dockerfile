FROM golang:1.17-alpine3.15
WORKDIR /pns
ENV GOPROXY=https://proxy.golang.com.cn,direct
COPY .. .
RUN apk add make
RUN go mod tidy
RUN make all
RUN mv /pns/build/pns /pns/pns
EXPOSE 10086 10087 10088
VOLUME [ "/pns/log" ]
CMD [ "./pns" ]