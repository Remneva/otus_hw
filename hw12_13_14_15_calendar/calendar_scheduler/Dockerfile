FROM golang:1.15.2 as build

ENV BIN_FILE /opt/calendar/calendar-scheduler
ENV CODE_DIR /usr/local/go/src

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .

ENV GOROOT /usr/local/go/
ADD . /usr/local/go/src/

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build -mod=mod -v\
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} calendar_scheduler/

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/calendar/calendar-scheduler"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar-scheduler/config.toml
COPY ./calendar_scheduler/config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}