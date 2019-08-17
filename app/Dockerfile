FROM golang:latest 

RUN mkdir /go/src/github.com
RUN mkdir /go/src/github.com/gedelumbung
RUN mkdir /go/src/github.com/gedelumbung/go-ga

RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/gedelumbung/go-ga

WORKDIR /go/src/github.com/gedelumbung/go-ga 

RUN dep ensure

CMD ["go", "run", "main.go"]