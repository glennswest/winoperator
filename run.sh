oc delete dc winoperator
sleep 5
oc run --rm -i winoperator --image=docker.io/glennswest/winoperator:0.2 

