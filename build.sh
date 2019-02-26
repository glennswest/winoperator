oc new-build --strategy docker --name winoperator --from-repo=https://github.com/glennswest/winoperator
oc start-build winoperator --from-dir . --follow
#eval $(minishift docker-env)
#docker build -t winoperator .
