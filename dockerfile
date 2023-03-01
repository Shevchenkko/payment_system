# build stage
FROM golang:1.17 AS build

# set env variables
ENV GOOS=linux \
  GOARCH=amd64 \
  CGO_ENABLED=0 

# copy and download mods
WORKDIR /srv/app/pkg
COPY go.mod .
COPY go.sum .
RUN go mod download

# copy and build code
COPY . .
RUN go build -o /srv/app/app main.go

# run stage
FROM golang:1.17-alpine as run
ENV GIN_MODE release
EXPOSE 8080

# copy binary
WORKDIR /srv
RUN mkdir -p /srv
COPY --from=build /srv/app/app /srv/app

# run binary
CMD ["/srv/app"]