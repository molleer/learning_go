FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v

ENV GO_USER_SERVICE_SECRET noTMwMrtsxtYfEFt+VaTXG3mEswCOMVwKpAhjRRWy40=

CMD ["app"]