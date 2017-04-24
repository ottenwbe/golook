# initial draft of docker img for the server
FROM ubuntu:16.04
RUN apt -y update
RUN apt -y install git golang
ENV GOPATH /build/go
ENV GOLOOK $GOPATH/src/github.com/ottenwbe/golook
RUN mkdir -p $GOPATH/src
RUN git clone https://github.com/ottenwbe/golook.git $GOLOOK
RUN cd $GOLOOK; go get ./...; go install
RUN cp $GOPATH/bin/golook /usr/local/bin
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/golook"]