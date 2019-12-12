FROM alpine:latest
RUN apk --no-cache add ca-certificates

FROM scratch
EXPOSE 80
COPY main /main
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /go/lib/time/zoneinfo.zip
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/main"]
