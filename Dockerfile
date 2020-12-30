# development
FROM golang:1.15.6 AS development

# $GOPATHの問題があるので注意
WORKDIR /go/app
COPY src/ .
# RUN GO111MODULE=off go get github.com/oxequa/realize
RUN go install
# CMD ["realize", "start", "--server"]
CMD ["go", "run", "."]

# production builder
FROM golang:1.15.6 AS builder

# $GOPATHの問題があるので注意
WORKDIR /go/app
ADD src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/taurus .

# production runtime
FROM scratch AS production
COPY --from=builder /app/taurus /app

ENTRYPOINT ["/app/taurus"]
