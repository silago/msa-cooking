#FROM cooking-go:latest
FROM eu.gcr.io/mobile-223015/cooking-go:1.0.1
#RUN mkdir /go/src/cooking-users
#ADD . /go/src/cooking-users/
ADD . /go/src/github.com/silago/msa-cooking/cooking-users/
RUN go get github.com/silago/social-api

#RUN go get github.com/silago/msa-cooking/cooking-users/

#RUN go get github.com/golang/protobuf/jsonpb
#RUN go get github.com/twitchtv/twirp
WORKDIR /go/src/github.com/silago/msa-cooking/cooking-users/
RUN go get github.com/twitchtv/twirp
RUN go get github.com/gpmgo/gopm/modules/log
RUN go build -o app .
ENTRYPOINT ["/go/src/github.com/silago/msa-cooking/cooking-users/app"]
