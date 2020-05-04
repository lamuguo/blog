# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/lamuguo/blog

WORKDIR /go/src/github.com/lamuguo/blog

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go build -o /tmp/server .

# Run the outyet command by default when the container starts.
ENTRYPOINT /tmp/server

# Document that the service listens on port 8080.
EXPOSE 8080
