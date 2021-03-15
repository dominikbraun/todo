FROM golang:1.15-alpine AS build

WORKDIR /src

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /out/todo .

FROM alpine:3.13 AS final

LABEL maintainer="Dominik Braun <mail@dominikbraun.io>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="todo"
LABEL org.label-schema.description="A simple ToDo REST API."
LABEL org.label-schema.url="https://github.com/dominikbraun/todo"
LABEL org.label-schema.vcs-url="https://github.com/dominikbraun/todo"

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.8.0/wait /wait
RUN chmod +x /wait

COPY --from=build /out/todo /bin/todo

ENTRYPOINT /wait && /bin/todo

EXPOSE 8000