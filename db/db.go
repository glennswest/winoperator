package db

import (
    "regexp"
    "log"
    "os"
    "strings"
    "errors"
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
     SetDbValue("global.workerign",workerign)
     machinemanurl := os.Getenv("WINMACHINEMAN")
     SetDbValue("global.machinemanurl",machinemanurl)
     kubeconfigdata := os.Getenv("KUBECONFIGDATA")
     SetDbValue("global.kubeconfigdata",kubeconfigdata)

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




func Smartsplit(s string) []string {
    r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`) 
    arr := r.FindAllString(s, -1) 
    return(arr)
}
