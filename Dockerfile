#base image
FROM golang:1.22

#create work directory
WORKDIR /app

#in order to download modules, we need go.mod and go.sum in the work directory
COPY go.mod go.sum ./

#download modules
RUN go mod download

#copy source code
COPY ./ ./

#WORKDIR /r2-api-go

#compile
RUN CGO_ENABLED=0 GOOS=linux go build -o /r2-api-go

EXPOSE 8080

#what command to run, when container is started from this image
CMD ["/r2-api-go"]