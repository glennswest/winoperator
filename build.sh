export GIT_COMMIT=$(git rev-parse --short HEAD)
#eval $(minishift docker-env)
docker build --no-cache -t glennswest/winoperator .
docker tag glennswest/winoperator:latest docker.io/glennswest/winoperator:$GIT_COMMIT
docker push docker.io/glennswest/winoperator:$GIT_COMMIT
