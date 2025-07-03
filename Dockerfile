FROM golang:1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -ldflags "-s -w" -o /app_bin main.go
RUN upx --brute /app_bin

FROM scratch

COPY --from=builder /app_bin /app_bin

EXPOSE 9876

CMD [ "/app_bin" ]
