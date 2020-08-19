FROM golang:latest AS build

WORKDIR	/go/src/fluentbit-go-somewhere
COPY ./ ./
ENV	CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -buildmode=c-shared -o out_gcl.so ./cmd/

FROM fluent/fluent-bit:1.5.2

#edit here if you need to specify aws profile name.
#ENV AWS_PROFILE={INSERT_YOUR_AWS_PROFILE_NAME_HERE}

WORKDIR /fluent-bit/etc
COPY --from=build /go/src/fluentbit-go-somewhere/out_gcl.so ./out_gcl.so
COPY --from=build /go/src/fluentbit-go-somewhere/flb.conf ./flb.conf
CMD ["/fluent-bit/bin/fluent-bit", "-e", "/fluent-bit/etc/out_gcl.so", "-c", "/fluent-bit/etc/flb.conf"]