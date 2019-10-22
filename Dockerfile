FROM golang as build-stage
COPY go.mod /
COPY go.sum /
COPY main.go /
RUN cd / && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-prom-app

FROM alpine
COPY --from=build-stage /go-prom-app /
EXPOSE 8080
CMD ["/go-prom-app"]