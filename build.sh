#eval $(minishift docker-env)
docker build --no-cache -t glennswest/winoperator .
docker tag glennswest/winoperator:latest docker.io/glennswest/winoperator:latest
docker push docker.io/glennswest/winoperator:latest
