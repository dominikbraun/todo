
FROM golang:1.15-alpine AS build

WORKDIR /src

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /out/todo .

FROM alpine AS final

LABEL maintainer="Dominik Braun <mail@dominikbraun.io>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="todo"
LABEL org.label-schema.description="A simple ToDo REST API."
LABEL org.label-schema.url="https://github.com/dominikbraun/todo"
LABEL org.label-schema.vcs-url="https://github.com/dominikbraun/todo"

COPY --from=build /out/todo /bin/todo

ENTRYPOINT ["/bin/todo"]

EXPOSE 8000