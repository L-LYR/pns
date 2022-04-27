FROM golang:1.17-alpine3.15
WORKDIR /pns
ENV GOPROXY=https://proxy.golang.com.cn,direct
COPY . .
RUN apk add --update alpine-sdk
RUN go mod tidy
RUN make all
WORKDIR /pns/build
EXPOSE 10086 10087 10088
VOLUME [ "/pns/build/log" ]
CMD [ "./pns" ]