#eval $(minishift docker-env)
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main .
chmod +x main
docker build -t winoperator .
