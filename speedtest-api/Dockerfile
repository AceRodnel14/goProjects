FROM golang:alpine as builder

WORKDIR /src

RUN sed -i -e 's/https/http/' /etc/apk/repositories \
  && apk update

ENV USER=speedtest
ENV UID=1414

RUN adduser \    
    --disabled-password --gecos "" \    
    --home "/" --shell "/sbin/nologin" \    
    --no-create-home --uid "${UID}" "${USER}"

COPY ./go.* ./
RUN go mod download

COPY ./ ./

RUN export ARCH=$(uname -m) && \
    if [ $ARCH == 'armv7l' ]; then export ARCH=arm; fi && \
    wget -O /src/exec/speedtest.tgz https://bintray.com/ookla/download/download_file?file_path=ookla-speedtest-1.0.0-${ARCH}-linux.tgz && \
    tar zxvf /src/exec/speedtest.tgz -C /src/exec && \
    rm -f /srcexec/speedtest.md /srcexec/speedtest.5 /srcexec/speedtest.tgz
RUN mkdir -p /src/resources

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/server ./cmd/speedtest-api/

FROM alpine AS bin

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /src/exec/ /exec/
COPY --from=builder /src/resources /resources/
RUN chown -R speedtest:speedtest /resources

USER speedtest:speedtest

COPY --from=builder /out/server ./server
EXPOSE 9096

ENTRYPOINT  [ "/server" ]