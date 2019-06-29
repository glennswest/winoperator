package main

import (
    "fmt"
    "flag"
    "regexp"
    "log"
    "os"
    "os/signal"
    "syscall"
    "net"
    "net/http"
    "time"
    "strings"
    "bytes"
    "errors"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
    //"github.com/tidwall/gjson"
    "github.com/tidwall/sjson"
    "bufio"
    bolt "go.etcd.io/bbolt"
)



// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func GetBucketAndKey(k string) (string, string){
     idx := strings.IndexByte(k,'.')
     bucket := k[0:idx]
     key := k[idx+1:]
     return bucket,key
}

func GetDbValue(p string) string{
     val := ""
     //log.Printf("GetDbValue(%s)\n",p)
     b,k := GetBucketAndKey(p)
     db, _ := bolt.Open("/data/winoperator", 0600,nil)
     db.View(func(tx *bolt.Tx) error {
              bucket := tx.Bucket([]byte(b))
              if bucket == nil {
                   log.Printf("No Bucket\n")
                   return errors.New("No Key ")
                   }

        val = string(bucket.Get([]byte(k)))
        return nil
    })
     //log.Printf("External value = %s\n",val)
     db.Close()
     return val
}
func SetDbValue(k string,v string){
     //log.Printf("SetDbValue(%s,%s)\n",k,v)
     bucket,key := GetBucketAndKey(k)
     db, _ := bolt.Open("/data/winoperator", 0600,nil)
     db.Update(func(tx *bolt.Tx) error {
           b, err := tx.CreateBucketIfNotExists([]byte(bucket))
           if err != nil {
              log.Printf("Err in create of bucket\n")
              return err
              }
           //log.Printf("Setting %s to %s\n",key,v)
           b.Put([]byte(key),[]byte(v))
           return(nil)
           })
     db.Close()
     return
}



func InitDb(){
     SetDbValue("global.dbversion","1.0")
     SetDbValue("global.User","Administrator")
     SetDbValue("global.Password","Secret2018")
     SetDbValue("global.ocpversion","3.11")
     master_host := os.Getenv("MASTERHOST")
     SetDbValue("global.master",master_host)
     SetDbValue("global.sshuser","root")
     sshkey := os.Getenv("SSHKEY")
     SetDbValue("global.sshkey",sshkey)
     workerign := os.Getenv("WORKERIGN")
     SetDbValue("global.workerign")
}

func SetupDb() {
    _ = os.MkdirAll("/data", 0700)
    dbexists := Exists("/data/winoperator")
    if (dbexists == false){
             log.Printf("Setup Database")
             InitDb()
       } else {
             log.Printf("Using Existing Database")
        }
}



func init() {
  flag.Parse();
}

func get_pod_ip(c *kubernetes.Clientset, podname string) string {
        log.Printf("Getting Pod Ip: %s\n",podname)
        
	pods, err := c.Core().Pods(podname).List(v1.ListOptions{})
	if err != nil {
	  // handle error
          log.Printf("get_pod_ip: Error  %v\n",err)
          return ""
          }
        if (pods == nil){
            log.Printf("get_pod_ip: Pods empty\n")
            return ""
            }
        if (len(pods.Items) == 0){
            log.Printf("get_pod_ip: No Pods\n")
            return ""
            }
        for _, pod := range pods.Items {
	   log.Printf("Pod %s - Ip %s\n",pod.Name, pod.Status.PodIP)
           }
         log.Printf("%v\n",pods)
         pod := pods.Items[0]
         return(pod.Status.PodIP)
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

func ArAdd(d string,aname string,v1 string,v2 string) string{
      s := `{"` + v1 + `":"` + v2 + `"}`
      a := aname + ".-1"
      d,_ = sjson.SetRaw(d,a,s)
      return d
      }

func build_variables(c *kubernetes.Clientset, node_name string) string {
     d := `{"version": 1, "labels": [], "annotations": []}`
     selector := "metadata.name=" + node_name;
     nodes, err := c.Core().Nodes().List(v1.ListOptions{FieldSelector: selector})

     if (err != nil){
        log.Printf("Cannot get node info\n")
        return d;
        }
     for index, element := range nodes.Items[0].Labels {
         log.Printf("%s -> %s", index,element);
         d = ArAdd(d,"labels",index,element)
         }
     for index, element := range nodes.Items[0].Annotations {
         log.Printf("%s -> %s", index,element);
         d = ArAdd(d,"annotations",index,element)
         }
    nul := node_name + ".UserName"
    node_user := GetDbValue(nul)
    node_password := ""
    if (node_user == ""){
       node_user = GetDbValue("global.User")
       node_password = GetDbValue("global.Password")
      } else {
       npk := node_name + ".UserPassword"
       node_password = GetDbValue(npk)
      }
    d = ArAdd(d,"settings","user",node_user)
    d = ArAdd(d,"settings","password",node_password)
    master_host := GetDbValue("global.master")
    d = ArAdd(d,"settings","master",master_host)
    sshuser := GetDbValue("global.sshuser")
    d = ArAdd(d,"settings","sshuser",sshuser)
    sshkey := GetDbValue("global.sshkey")
    d = ArAdd(d,"settings","sshkey",sshkey)
    workerign := GetDbValue("global.workerign")
    d = ArAdd(d,"settings","workerign",workerign)
    //log.Printf("d = %s\n", d);
    return d
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
     v := build_variables(c,node_name)
     log.Printf("Variables =%s\n",v)
     winmachineman_ip := GetMachineManIp(c)
     wmmurl := "http://" + winmachineman_ip + ":8080/machines"
     resp, err := http.Post(wmmurl,"application/json", bytes.NewBuffer([]byte(v)))
     log.Printf("Response = %s %s\n",resp,err)
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


func GetMachineManIp(c *kubernetes.Clientset) string {
    winmachineman_ip := ""
    for winmachineman_ip == "" {
       winmachineman_ip = get_pod_ip(c,"winmachineman")
       log.Printf("IP = %s\n",winmachineman_ip)
       if (winmachineman_ip == ""){
          log.Printf("Waiting on Windows Machine Manager")
          time.Sleep(10 * time.Second)
          }
       }
    return winmachineman_ip
}

func smartsplit(s string) []string {
    r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`) 
    arr := r.FindAllString(s, -1) 
    return(arr)
}

func docli(){
     for {
         handlecli()
         }
}

func handlecli(){
     scanner := bufio.NewScanner(os.Stdin)
     for scanner.Scan() {
          cmdline := scanner.Text();
          process_cli(cmdline)
      }
     return
}

func trimQuotes(s string) string {
    if len(s) >= 2 {
        if s[0] == '"' && s[len(s)-1] == '"' {
            return s[1 : len(s)-1]
        }
    }
    return s
}

func process_cli(text string){
       if (len(text) == 0){
          return
          }
        cmd := smartsplit(text)
        switch(cmd[0]){
           case "help":
              log.Printf("HELP is being used\n")
              fmt.Printf("\n WinOperator Help\n")
              fmt.Printf("Commands: \n")
              fmt.Printf("   help      - This CLI Help\n")
              fmt.Printf("   ls        - List directory content\n")
              fmt.Printf("   set.value - Set a value in the winoperator settings db\n")
              fmt.Printf("   get.value - Get a value in the winoperator settings db\n")
              fmt.Printf("   quit      - Kill the winoperator\n")
              break
           case "set.value":
              if (len(cmd) < 3){
                 fmt.Printf("More arguments needed: set.value variable value\n")
                 return
                 }
              value := trimQuotes(cmd[2])
              SetDbValue(cmd[1],value)
              break
           case "get.value":
              if (len(cmd) < 2){
                 fmt.Printf("More arguments needed: get.value variable\n")
                 return
                 }
              thevalue := GetDbValue(cmd[1])
              fmt.Printf("%s = %s\n",cmd[1],thevalue)
              break
           case "set":
              switch(len(cmd)){
              case 1:
                   process_cli("get.value global.User")
                   process_cli("get.value global.Password")
                   process_cli("get.value global.ocpversion")
                   process_cli("get.value global.master")
                   process_cli("get.value global.sshkey")
                   return
                   break
              case 2:
                   thevalue := GetDbValue(cmd[1])
                   fmt.Printf("%s = %s\n",cmd[1],thevalue)
                   break
              case 3:
                   value := trimQuotes(cmd[2])
                   SetDbValue(cmd[1],value)
                   break
               }
              break
           case "quit":
              os.Stdin.Close()
              break
          }
    fmt.Printf("winoperator> ")
}

func main() {
    //f, err := os.OpenFile("/winoperator.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    //if err != nil {
    //   log.Fatalf("error opening file: %v", err)
    //   }
    //defer f.Close()
    //log.SetOutput(f)
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
    SetupDb()
    signal.Ignore(syscall.SIGHUP)
    signal.Ignore(syscall.SIGINT)
    //go docli();
    

    winmachineman_ip := GetMachineManIp(clientset)
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

