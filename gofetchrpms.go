// Dependencies: golang.org/x/net/html
// go get golang.org/x/net/html
// Program to download rpms from a given url link concurrently
// Usage go run fetch.go <url> -list <to list rpms>
// go run fetch.go <url> -download to download the rpms in the current
// Directory from where the file was executed
// Obviously more improvements are required.
package main

import (
       "fmt"
       "net/http"
       "golang.org/x/net/html"
       "strings"
       "os"
       "flag"
       "time"
       "io"
       )

var (
    url = flag.String("url", "", "Specify the url to download rpms")
    list = flag.Bool("list", false, "List rpms")
    download = flag.Bool("download", false, "Download rpms")
    )

func main() {
     
     flag.Parse()
     ch := make(chan string) 
     rpmList := []string{}
     listptr := &rpmList 
     getRpmList(url, listptr)
     start := time.Now()
     if *list {
        listRpms(listptr)
     }
     if *download {
        for _, filename := range rpmList {
             rpmlink := *url + "/" + filename 
             go fetch(filename, rpmlink, ch)
        }
        for index := range rpmList {
            fmt.Println(index, <-ch)
        }
        fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
    }
    
}     

func listRpms(listptr *[]string) {
     for _, rpm := range *listptr {
         fmt.Println(rpm)
      }
}
     
func fetch(filename string, rpmlink string, ch chan<- string)  {
     start := time.Now() 
     out, err := os.Create(filename)
     if err != nil {
        return
     }
     defer out.Close() 
     
     resp, err := http.Get(rpmlink)
     if err != nil {
        return 
     }
     defer resp.Body.Close()
     nbytes, err := io.Copy(out, resp.Body)
     if err != nil {
        ch <- fmt.Sprintf("While reading:%s: %", rpmlink, err)
        return
     }
     secs := time.Since(start).Seconds() 
     ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, filename)
}     
    


func getRpmList(url *string, listPtr *[]string) {
     resp, err := http.Get(*url)
     if err != nil {
        fmt.Fprintf(os.Stderr, "fetch: %v \n", err) 
        os.Exit(1)
     }
     defer resp.Body.Close() 
     Tokenizer := html.NewTokenizer(resp.Body)
     for {
         tt := Tokenizer.Next() 
         switch {
                case tt == html.ErrorToken:
                     return 
                case tt == html.StartTagToken:
                     tag := Tokenizer.Token() 
                     isAnchor := tag.Data == "a"
                     if isAnchor {
                        for _, a := range tag.Attr {
                            if a.Key == "href" {
                               if strings.Contains(a.Val, ".rpm") == true {
                                  fmt.Println(a.Val)
                                  *listPtr = append(*listPtr, a.Val) 
                               }
                               break 
                            }
                        }
                    }
               }
         }
}

           
