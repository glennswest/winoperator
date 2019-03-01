package main

import (
    "flag"
    "log"
    "os"
    "reflect"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
)

func init() {
  flag.Parse();
}

func get_node(c *kubernetes.Clientset, node_name string){
     selector := "metadata.name=" + node_name;
     nodes, err := c.Core().Nodes().List(v1.ListOptions{FieldSelector: selector})

     if (err != nil){
        log.Printf("Cannot get node info\n");
        return;
        }
     //log.Printf("Nodes = %d\n", nodes.Items);
     theos := nodes.Items[0].Labels["beta.kubernetes.io/os"];
     log.Printf("OS = %s\n",theos);
     log.Printf("Type = %s\n",reflect.TypeOf(nodes.Items[0]))

}

func add_node(c *kubernetes.Clientset, node_name string){
     get_node(c,node_name);
}

func delete_node(c *kubernetes.Clientset,node_name string){
}


func main() {
    log.SetOutput(os.Stdout)
    log.Printf("Building config from flags\n")
    config, err := clientcmd.BuildConfigFromFlags("", "")
    if err != nil {
        log.Printf("Failled: BuildConfigFromFlags\n");
        panic(err.Error())
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Printf("Failed: NewForConfig\n")
        panic(err.Error())
    }

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
            add_node(clientset,nodename);
        },
        DeleteFunc: func(obj interface{}) {
            // "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
            // interface that allows us to get metadata easily
            mObj := obj.(v1.Object)
            nodename := mObj.GetName();
            log.Printf("Node Delete from Store: %s\n", nodename)
            delete_node(clientset,nodename);
        },
    })

    informer.Run(stopper)
}

