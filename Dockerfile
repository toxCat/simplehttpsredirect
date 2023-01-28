FROM golang

WORKDIR /go/src/app
COPY . .

RUN go install -v net/http
RUN go install -v net/src
RUN go install -v net

CMD ["app"]
