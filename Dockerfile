ARG APP_SRC=/usr/src/app
ARG BUILD_FILE=lunch-roulette

FROM golang:1.13-alpine
ARG BUILD_FILE
ARG APP_SRC

WORKDIR $APP_SRC

COPY go.mod .
COPY go.sum .

RUN apk add git &&\
    apk --update add ca-certificates &&\
    go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go test ./... &&\
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $BUILD_FILE .


FROM bash:4.3.48
ARG BUILD_FILE
ENV BUILD_FILE $BUILD_FILE
ARG APP_SRC

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 ${APP_SRC}/${BUILD_FILE} $BUILD_FILE

RUN addgroup -S appuser && adduser -S appuser -G appuser -u 1000 &&\
    chown -R appuser $BUILD_FILE
USER appuser

CMD ./$BUILD_FILE
