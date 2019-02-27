oc delete dc winoperator
oc delete is winoperator
oc delete rc winoperator
oc delete sa winoperator
sleep 3
#oc create sa winoperator
oc create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default
#oc adm policy add-role-to-user admin system:serviceaccount:winoperator 
oc new-app docker.io/glennswest/winoperator:latest
oc set volume dc/winoperator --add --name=logs --mount-path=/tmp --path=/data/winoperator
