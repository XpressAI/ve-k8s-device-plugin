FROM golang:1.16.4-alpine3.13
RUN mkdir -p /go/src/NECVectorEngine/k8s-device-plugin
ADD . /go/src/NECVectorEngine/k8s-device-plugin
WORKDIR /go/src/NECVectorEngine/k8s-device-plugin/cmd/k8s-device-plugin
RUN go install 


FROM alpine:3.13
WORKDIR /root/
COPY --from=0 /go/bin/k8s-device-plugin .
CMD ["./k8s-device-plugin", "-logtostderr=true", "-stderrthreshold=INFO", "-v=5"]
