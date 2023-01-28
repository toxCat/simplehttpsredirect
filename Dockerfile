FROM golang

WORKDIR /go/src/app
COPY . .

RUN go install -v net/http

CMD ["app"]
