FROM golang:1-alpine AS builder

WORKDIR /app

# Uncoment this if you really want your image to be less than ~5MiB 
# RUN apk add --no-cache upx  

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -ldflags "-s -w" -o /app_bin main.go

# Uncoment this if you really want your image to be les than ~5MiB. Takes some time
# RUN upx --brute /app_bin

FROM scratch

COPY --from=builder /app_bin /app_bin

EXPOSE 9876

CMD [ "/app_bin" ]
