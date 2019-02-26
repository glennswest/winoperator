package main;

import (
	"fmt"
        "log"
        "net/http"
        "time"
        "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
        )


func main() {
    //Configure cluster info
    config, err := rest.InClusterConfig()
    if err != nil {
	panic(err.Error())
	}

    //Create a new client to interact with cluster and freak if it doesn't work
    kubeClient, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalln("Client not created sucessfully:", err)
    }
    //Create a cache to store nodes
    var nodeStore cache.Store
    //Watch for Nodes
    nodeStore = watchNodes(kubeClient, nodeStore)
    //Keep alive
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func watchNodes(client *client.Client, store cache.Store) cache.Store {
    //Define what we want to look for (Nodes)
    watchlist := cache.NewListWatchFromClient(client, "Nodes", api.NamespaceAll, fields.Everything())
    resyncPeriod := 30 * time.Minute
    //Setup an informer to call functions when the watchlist changes
    eStore, eController := framework.NewInformer(
        watchlist,
        &api.Pod{},
        resyncPeriod,
        framework.ResourceEventHandlerFuncs{
            AddFunc:    nodeCreated,
            DeleteFunc: nodeDeleted,
        },
    )
    //Run the controller as a goroutine
    go eController.Run(wait.NeverStop)
    return eStore
}

func nodeCreated(obj interface{}) {
    node := obj.(*api.Node)
    fmt.Println("Node created: " + node.ObjectMeta.Name)
}

func nodeDeleted(obj interface{}) {
    node := obj.(*api.Node)
    fmt.Println("Node deleted: " + node.ObjectMeta.Name)
}

