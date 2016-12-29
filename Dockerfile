FROM ubuntu:16.04

# Install packaged dependencies
RUN apt-get update && apt-get install -y \
	chromium-browser \
	chromium-chromedriver \
	curl \
	ffmpeg \
	xvfb

# Install Go
ENV GO_VERSION 1.7.3
RUN curl -fsSL "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -xzC /usr/local
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

# Install code
ADD . /go/src/github.com/jordanpotter/site-analyzer
RUN go install github.com/jordanpotter/site-analyzer

# Create X11 socket directory
RUN mkdir /tmp/.X11-unix && chmod 1777 /tmp/.X11-unix

# Create mount point to access data
RUN mkdir /data && chmod 777 /data
VOLUME ["/data"]

# Run commands as non-root user
RUN useradd -m -s /bin/bash docker
USER docker
WORKDIR /home/docker

# Run site-analyzer
ENTRYPOINT ["site-analyzer", "-videodir", "/data", "-chromedriver", "/usr/lib/chromium-browser/chromedriver"]
