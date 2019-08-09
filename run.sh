git pull
export GIT_COMMIT=$(git rev-parse --short HEAD)
echo $GIT_COMMIT
oc delete dc winoperator
oc delete is winoperator
oc delete sa winoperator
oc delete project winoperator
sleep 40
oc new-project winoperator
#oc import-image winoperator --from=docker.io/glennswest/winoperator:$GIT_COMMIT --confirm
#oc delete  istag/winoperator:latest
oc create sa winoperator
oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:winoperator:default
kubectl create deployment winoperator --image=docker.io/glennswest/winoperator:$GIT_COMMIT
#oc set volume dc/winoperator --add --name=dbdata --type=hostPath --path=/etc/winoperator --mount-path=/data
#oc new-app --docker-image=glennswest/winoperator:$GIT_COMMIT --token=$(oc sa get-token winoperator) 
#oc run winoperator --tty --stdin --image=glennswest/winoperator:$GIT_COMMIT
#oc policy add-role-to-user admin  system:serviceaccount:winoperator:default
sleep 20
export masterhostname=control-plane-0
#export sshkey=`cat ~/.ssh/id_rsa | base64 | tr -d '\n'`
export workerign=`cat worker.ign | base64 | tr -d '\n'`
export defaultdomain=$(oc describe --namespace=openshift-ingress-operator ingresscontroller/default | grep "Domain:" | cut -d ":" -f 2 | cut -d "." -f 1- | tr -d ' ')
export kubeconfigdata=`cat ~/.kube/config | base64 | tr -d '\n'`
oc set env deployment/winoperator WINMACHINEMAN=http://winmachineman.$defaultdomain SSHKEY=$sshkey MASTERHOST=$masterhostname WORKERIGN=$workerign KUBECONFIGDATA=$kubeconfigdata MYURL=winoperatordata.$defaultdomain
#oc patch dc winoperator -p "spec:
#  template:
#    spec:
#      containers:
#      - name: winoperator
#        tty:   true
#        stdin: true"
#

