package main

import (
        "bytes"
        "fmt"
        "os"
        "os/exec"
        //"regexp"
        "strings"
        "time"
        "strconv"
        "net/http"
        "sync"
        "io"
)


func main() {
        var wg sync.WaitGroup
        var ipChina string
        var ipWg77k string
        var ipIfconfig string
        wg.Add(3)
        go myIP("https://w.g77k.com/ip.php", &ipChina, &wg)
        go myIP("https://f.g77k.com/ip.php", &ipWg77k, &wg)
        go myIP("https://ifconfig.me", &ipIfconfig, &wg)

        wg.Wait()
        fmt.Println("MyIP China:", ipChina)
        fmt.Println("MyIP w.g77k.com:", ipWg77k)
        fmt.Println("MyIP ifconfig:", ipIfconfig)

        for{
                time.Sleep(1 * time.Second)
                get_ping("114.114.114.114", ipChina)
        }
}


func get_ping(host string, myip string) {
        cmd := exec.Command("./get_ping.sh", host)
        cmdOutput := &bytes.Buffer{}
        cmd.Stdout = cmdOutput
        //printCommand(cmd)
        err := cmd.Run()
        printError(err)
        output := cmdOutput.Bytes()
        //printOutput(output)
        
        val, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
        write_val := fmt.Sprintf("%.2f", val)
        
        t := time.Now()
        formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())

        line := formatted + " FROM "+ myip+ " PING " + host + " " + string(write_val) +"ms"
        fmt.Println(line)


}

func printCommand(cmd *exec.Cmd) {
        fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
        if err != nil {
                os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
        }
}

func printOutput(outs []byte) {
        if len(outs) > 0 {
                fmt.Printf("==> Output: %s\n", string(outs))
        }
}


func myIP(ipServerURL string, ret *string, wg *sync.WaitGroup) string {
        defer wg.Done()
        resp, err := http.Get(ipServerURL)
        if err != nil {
                *ret = ipServerURL + " " + err.Error()
                return err.Error()
        }
        defer resp.Body.Close()
        if resp.StatusCode == http.StatusOK {
                bodyBytes, err := io.ReadAll(resp.Body)
                if err != nil {
                        *ret = ipServerURL + " " + err.Error()
                        return err.Error()
                }
                bodyString := string(bodyBytes)
                *ret =  bodyString
        } else {
                *ret =  "ERROR with read content"
        }
        return *ret
}
