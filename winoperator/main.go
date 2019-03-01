package main

import (
  //  log "github.com/sirupsen/logrus"
    //"os"
    "fmt"
    "flag"

    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
)

func init() {
  flag.Parse();
  // Log as JSON instead of the default ASCII formatter.
  //log.SetFormatter(&log.JSONFormatter{})
  //log.SetFormatter(&log.TextFormatter{})

  // Output to stdout instead of the default stderr
  // Can be any io.Writer, see below for File example
  //log.SetOutput(os.Stdout)

  // Only log the warning severity or above.
  //log.SetLevel(log.TraceLevel)
  //logger := logrus.New()
  //logger.Formatter = &logrus.JSONFormatter{}

  // Use logrus for standard log output
  // Note that `log` here references stdlib's log
  // Not logrus imported under the name `log`.
  //log.SetOutput(logger.Writer())
}

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

