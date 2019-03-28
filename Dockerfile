FROM golang:1.12 AS build

WORKDIR /go/src/yap

RUN go get github.com/golang/dep/cmd/dep
COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -vendor-only

COPY util ./util/
COPY eval ./eval/
COPY alg ./alg/
COPY nlp ./nlp/
COPY app ./app/
COPY webapi ./webapi/
COPY main.go ./
RUN go build -o /bin/yap


# new image
FROM golang:1.12
WORKDIR /app

COPY ./data/md_model_temp_i9.b64 ./data/
COPY ./data/bgulex ./data/bgulex
COPY ./conf ./conf/

COPY --from=build /bin/yap ./yap

EXPOSE 8000
ENTRYPOINT ["./yap", "api", "-tagonly"]