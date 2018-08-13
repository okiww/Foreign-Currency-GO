FROM golang:1.8

WORKDIR /go/src/Foreign-Currency-GO
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && go install -ldflags='-X main.version=1.0.0' \
    && cp docker.yml .env.yml

EXPOSE 8080
VOLUME [ "/var/log/Foreign-Currency-GO" ]

ENTRYPOINT [ "Foreign-Currency-GO", "-migrate", "-log_dir", "/var/log/Foreign-Currency-GO", "-alsologtostderr", "-stderrthreshold", "warning", "-v", "2" ]