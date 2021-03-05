FROM golang:latest as builder

LABEL maintainer="Sahadat Hossain"

# set the current working directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# build the Go app (API server)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o controller .

########### start a new stage from scratch ###########
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# copy the pre-built binary file from the previous stage
COPY --from=builder /app/controller .

# Expose port 8080 to the outside world
EXPOSE 8080

# command to run the executable
CMD ["./controller"]