export GIT_COMMIT=$(git rev-parse --short HEAD)
echo $GIT_COMMIT
docker pull docker.io/glennswest/winoperator:$GIT_COMMIT
oc delete dc winoperator
oc delete is winoperator
oc delete rc winoperator
oc delete sa winoperator
oc delete project winoperator
sleep 3
oc new-project winoperator
oc create sa winoperator
oc policy add-role-to-user admin  system:serviceaccount:winoperator:winoperator
oc new-app glennswest/winoperator:$GIT_COMMIT --token=$(oc sa get-token winoperator) 
