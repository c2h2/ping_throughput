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

        t := time.Now()
        formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
        filename := formatted + "_pings.txt"

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

        num:=0
        //var results []float64
        for{
                time.Sleep(500 * time.Millisecond)
                get_ping("114.114.114.114", ipChina, num, filename)
                num++
                /*if(res< 50000){
                        results.append(results, res)
                }*/
                //fmt.Println("num/min/max/avg = ",len(results),min(results),max(results),sum(results)/len(results))
                
        }
}


func get_ping(host string, myip string, num int, filename string) float64{
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

        line :=  strconv.Itoa(num) + " " + formatted + " FROM "+ myip+ " PING " + host + " " + string(write_val) +"ms"
        fmt.Println(line)

        text := formatted +","+myip+","+string(write_val)+"\n"
        
        f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        if err != nil {
            panic(err)
        }

        defer f.Close()

        if _, err = f.WriteString(text); err != nil {
            panic(err)
        }
        return val
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
