package main


import (
    "fmt"
    "flag"
    "log"
    "bufio"
    "os"
    "github.com/glennswest/winoperator/db"
)

func trimQuotes(s string) string {
    if len(s) >= 2 {
        if s[0] == '"' && s[len(s)-1] == '"' {
            return s[1 : len(s)-1]
        }
    }
    return s
}


func init() {
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

func process_cli(text string){
       if (len(text) == 0){
          return
          }
        cmd := db.Smartsplit(text)
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
              db.SetDbValue(cmd[1],value)
              break
           case "get.value":
              if (len(cmd) < 2){
                 fmt.Printf("More arguments needed: get.value variable\n")
                 return
                 }
              thevalue := db.GetDbValue(cmd[1])
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
                   thevalue := db.GetDbValue(cmd[1])
                   fmt.Printf("%s = %s\n",cmd[1],thevalue)
                   break
              case 3:
                   value := trimQuotes(cmd[2])
                   db.SetDbValue(cmd[1],value)
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
    
    Cptr := flag.String("c", "", "set values")

    flag.Parse();
    fmt.Println("coption:", *Cptr)
    docli();
}


