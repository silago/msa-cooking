#FROM cooking-go:latest
FROM eu.gcr.io/mobile-223015/cooking-go:1.0.1
RUN mkdir /go/src/cooking-users
ADD . /go/src/cooking-users/
RUN go get github.com/silago/social-api
#RUN go get github.com/golang/protobuf/jsonpb
#RUN go get github.com/twitchtv/twirp
WORKDIR /go/src/cooking-users/
##RUN go get "github.com/rs/xid"

#RUN go get -u "github.com/silago/ok-go/api"
#RUN cat /go/src/github.com/silago/ok-go/api/api.go
#RUN go get "github.com/go-redis/redis"
#RUN go get "github.com/nats-io/go-nats"
#RUN go get "github.com/Jeffail/gabs"
RUN go get github.com/twitchtv/twirp
RUN go get github.com/gpmgo/gopm/modules/log
RUN go build -o app .
ENTRYPOINT ["/go/src/cooking-users/app"]
