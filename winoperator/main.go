package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "fmt"

    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    
    fmt.Printf("Building config from flags\n")
    config, err := clientcmd.BuildConfigFromFlags("", "")
    if err != nil {
        fmt.Printf("Failled: BuildConfigFromFlags\n");
        panic(err.Error())
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Failed: NewForConfig\n")
        panic(err.Error())
    }

    fmt.Printf("Setting up Informer")
    factory := informers.NewSharedInformerFactory(clientset, 0)
    informer := factory.Core().V1().Nodes().Informer()
    stopper := make(chan struct{})
    defer close(stopper)
    informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            // "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
            // interface that allows us to get metadata easily
            mObj := obj.(v1.Object)
            fmt.Printf("New Node Added to Store: %s", mObj.GetName())
        },
    })

    informer.Run(stopper)
}

