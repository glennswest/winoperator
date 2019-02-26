oc delete dc winoperator
sleep 5
oc run --rm -i winoperator --image=winoperator --image-pull-policy=Never

