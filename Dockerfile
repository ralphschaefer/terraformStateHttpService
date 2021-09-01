FROM golang:1.16.6-buster
WORKDIR /build
COPY ./ ./
WORKDIR /build
RUN make

FROM busybox:uclibc
WORKDIR /runtime
RUN mkdir -p /share
COPY --from=0 /build/stateservice .
CMD /runtime/stateservice -port 80 -bind 0.0.0.0 -storeto /share
