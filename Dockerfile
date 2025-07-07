FROM golang:1.23 AS build

COPY . /app

WORKDIR /app

RUN go get -v

RUN go build -v -o hf-provisioner-digitalenergy

FROM golang:1.23 AS run

COPY --from=build /app/hf-provisioner-digitalenergy /hf-provisioner-digitalenergy

WORKDIR /

CMD ["/hf-provisioner-digitalenergy"]