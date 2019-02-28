export GIT_COMMIT=$(git rev-parse --short HEAD)
rm -r -f tmp
mkdir tmp
echo $GIT_COMMIT > tmp/commit.id
chmod -R 777 ./tmp
#eval $(minishift docker-env)
docker build --no-cache -t glennswest/winoperator:$GIT_COMMIT .
docker tag glennswest/winoperator:$GIT_COMMIT  docker.io/glennswest/winoperator:$GIT_COMMIT
docker push docker.io/glennswest/winoperator:$GIT_COMMIT
