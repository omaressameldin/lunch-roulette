FROM golang:1.12.4-alpine

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
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BUILD_FILE} .


FROM bash:4.3.48

ARG BUILD_FILE
ARG ENV
ARG APP_SRC
ARG PORT
ARG DB_DIRECTORY
ENV BUILD_FILE $BUILD_FILE

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY $ENV .
COPY --from=0 ${APP_SRC}/${BUILD_FILE} $BUILD_FILE

RUN addgroup -S appuser && adduser -S appuser -G appuser -u 1000 &&\
    chown -R appuser $BUILD_FILE $ENV $DB_DIRECTORY
USER appuser

EXPOSE $PORT

CMD ./$BUILD_FILE
