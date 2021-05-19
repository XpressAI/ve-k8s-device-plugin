FROM golang:1.16.4-alpine3.13
RUN apk --no-cache add git pkgconfig build-base libdrm-dev
RUN apk --no-cache add hwloc-dev --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community
RUN mkdir -p /go/src/NECVectorEngine/k8s-device-plugin
ADD . /go/src/NECVectorEngine/k8s-device-plugin
WORKDIR /go/src/NECVectorEngine/k8s-device-plugin/cmd/k8s-device-plugin
RUN go install 
#  -ldflags="-X main.gitDescribe=$(git -C /go/src/github.com/NECVectorEngine/k8s-device-plugin/ describe --always --long --dirty)" \ 
# github.com/hazimhasnan/Device-Plugin/cmd/k8s-device-plugin


FROM alpine:3.13
RUN apk --no-cache add ca-certificates libdrm
RUN apk --no-cache add hwloc --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community
WORKDIR /root/
COPY --from=0 /go/bin/k8s-device-plugin .
CMD ["./k8s-device-plugin", "-logtostderr=true", "-stderrthreshold=INFO", "-v=5"]