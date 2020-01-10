FROM golang:latest AS build

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch AS run
COPY --from=build /go/src/app/main .
CMD ["./main"]