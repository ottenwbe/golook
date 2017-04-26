# initial draft of docker img for the server
FROM golang:1.8.1
RUN apt -y update
ENV GOLOOK $GOPATH/src/github.com/ottenwbe/golook
RUN mkdir -p $GOLOOK
COPY . $GOLOOK
RUN cd $GOLOOK; go get ./...; go install
EXPOSE 8383
EXPOSE 8382
ENTRYPOINT ["golook"]