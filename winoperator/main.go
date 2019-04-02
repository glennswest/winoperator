package main

import (
    "flag"
    "log"
    "os"
    "net"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
    //"github.com/tidwall/gjson"
    "time"
)

func init() {
  flag.Parse();
}

func set_global_variable(vname string, value string){


}

func get_pod_ip(c *kubernetes.Clientset, podname string) string {
        log.Printf("Getting Pod Ip: %s\n",podname)
        
	pods, err := c.Core().Pods(podname).List(v1.ListOptions{})
	if err != nil {
	  // handle error
          log.Printf("get_pod_ip: Error  %v\n",err)
          return ""
          }
        for _, pod := range pods.Items {
	   log.Printf("Pod %s - Ip %s\n",pod.Name, pod.Status.PodIP)
           }
         pod := pods.Items[0]
         return(pod.Status.PodIP)
}

func update_global_variables(c *kubernetes.Clientset){
	

}

func get_node_label(c *kubernetes.Clientset, node_name string,thename string) string {
     selector := "metadata.name=" + node_name;
     nodes, err := c.Core().Nodes().List(v1.ListOptions{FieldSelector: selector})

     if (err != nil){
        log.Printf("Cannot get node info\n");
        return "";
        }
     //log.Printf("Nodes = %d\n", nodes.Items);
    theresult := nodes.Items[0].Labels[thename];
    return theresult
}

func get_node_annotation(c *kubernetes.Clientset, node_name string,thename string) string {
     selector := "metadata.name=" + node_name;
     nodes, err := c.Core().Nodes().List(v1.ListOptions{FieldSelector: selector})

     if (err != nil){
        log.Printf("Cannot get node info\n");
        return "";
        }
     //log.Printf("Nodes = %d\n", nodes.Items);
    theresult := nodes.Items[0].Annotations[thename];
    return theresult
}


func ip_lookup(tip string) string{

        ips, err := net.LookupHost(tip)
	if err != nil {
		log.Printf("Could not get IPs: %v\n", err)
		return("");
	   }
        theip := ""
	for _, ip := range ips {
		log.Printf("%s IN A %s\n", tip,ip)
	    }
        theip = ips[0];
        return(theip);
}

func check_windows_node(c *kubernetes.Clientset, node_name string){
     log.Printf("check_windows_node: %s\n",node_name)
     host_name := get_node_label(c,node_name,"kubernetes.io/hostname")
     log.Printf("hostname = %s\n",host_name);
   
     theip := get_node_annotation(c,node_name,"host/ip")
     if (theip == ""){
        log.Printf("host/ip annotation not set -- lookup up via dns\n");
        theip := ip_lookup(host_name);
        log.Printf("ip = %s\n",theip);
        }
     log.Printf("Host IP: %s\n",theip);

}

func kube_add_node(c *kubernetes.Clientset, node_name string){
     theos := get_node_label(c,node_name,"beta.kubernetes.io/os");
     log.Printf("OS = %s\n",theos)
     switch theos {
         case "linux":
         // Ignore Linux Nodes for now
         case "windows":
               check_windows_node(c,node_name);
         default:
               log.Printf("Undefined OS: %s (Ignored)\n",theos)
         }
              
}

func kube_delete_node(c *kubernetes.Clientset,node_name string){
}


func main() {
    log.SetOutput(os.Stdout)
    log.Printf("Version .001a\n")
    log.Printf("Building config from flags\n")
    config, err := clientcmd.BuildConfigFromFlags("", "")
    if err != nil {
        log.Printf("Failed: BuildConfigFromFlags\n");
        panic(err.Error())
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Printf("Failed: NewForConfig\n")
        panic(err.Error())
    }

    winmachineman_ip := ""
    for winmachineman_ip == "" {
       log.Printf("Waiting on Windows Machine Manager")
       winmachineman_ip := get_pod_ip(clientset,"winmachineman")
       if (winmachineman_ip == ""){
          time.Sleep(10 * time.Second)
          }
       }
    log.Printf("Windows Machine Man found at ip %s\n",winmachineman_ip)

    log.Printf("Setting up Informer\n")
    factory := informers.NewSharedInformerFactory(clientset, 0)
    informer := factory.Core().V1().Nodes().Informer()
    stopper := make(chan struct{})
    defer close(stopper)
    informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            // "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
            // interface that allows us to get metadata easily
            mObj := obj.(v1.Object)
            nodename := mObj.GetName();
            log.Printf("New Node Added to Store: %s\n", nodename)
            kube_add_node(clientset,nodename);
        },
        DeleteFunc: func(obj interface{}) {
            // "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
            // interface that allows us to get metadata easily
            mObj := obj.(v1.Object)
            nodename := mObj.GetName();
            log.Printf("Node Delete from Store: %s\n", nodename)
            kube_delete_node(clientset,nodename);
        },
    })

    informer.Run(stopper)
}

