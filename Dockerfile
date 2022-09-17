FROM golang:alpine
COPY . /infra
WORKDIR /infra
RUN sh -c "go mod tidy"
ENTRYPOINT [ "go", "run", "main.go" ]