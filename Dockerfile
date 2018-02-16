FROM drhayt/ora

ADD . /go/src/github.com/mangeshhendre/mathsvc

RUN export GOPATH=/go && \
	cd /go/src/github.com/mangeshhendre/mathsvc && \
	CGO_LDFLAGS=-L$ORACLE_HOME CGO_CFLAGS=-I$ORACLE_HOME/sdk/include go install ./...

EXPOSE 5678 7778

ENTRYPOINT ["/go/bin/mathsvc"]
CMD ["--help"]
