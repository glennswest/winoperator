git pull
export GIT_COMMIT=$(git rev-parse --short HEAD)
echo $GIT_COMMIT
oc delete dc winoperator
oc delete is winoperator
oc delete sa winoperator
oc delete project winoperator
sleep 3
oc new-project winoperator
oc set volume dc/winoperator --add --name=dbdata --type=hostPath --path=/etc/winoperator --mount-path=/data
oc import-image winoperator --from=docker.io/glennswest/winoperator:$GIT_COMMIT --confirm
oc delete  istag/winoperator:latest
#oc create sa winoperator
oc new-app glennswest/winoperator:$GIT_COMMIT --token=$(oc sa get-token winoperator) 
oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:winoperator:default
oc policy add-role-to-user admin  system:serviceaccount:winoperator:default

export masterhostname=$(hostname)
oc set env dc/winoperator MASTERHOST=$masterhostname

