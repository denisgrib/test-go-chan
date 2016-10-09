package main

import (
	"fmt"
    "os"
    "io/ioutil"
    "strings"
    "net/http"
)

func main() {

    var (
        maxConcurrentGrt int = 5
        searchString string = "Go"
        cntAllFound int
    )

    // urls, count all urls
    arrUrls, cntUrls := parseStdin()

    // chan for start next concurrent goroutine
	startNextGrt := make(chan struct{}, maxConcurrentGrt)

    // chan for wait one url
    done := make(chan bool)

    // chan for wait all urls
    waitAllUrl := make(chan bool)

	// start limit count goroutines
	for i := 0; i < maxConcurrentGrt; i++ {
		startNextGrt <- struct{}{}
	}

	// goroutine for controll done and run next goroutine
	go func() {
		for i := 0; i < cntUrls; i++ {
			<-done
			// run next
			startNextGrt <- struct{}{}
		}

        // all done
		waitAllUrl <- true
	}()

	// main loop
	for i := 0; i < cntUrls; i++ {

		<-startNextGrt

		go func(i int) {

            url := arrUrls[i]
            resp, err := http.Get(url)
            Panic(err)

            defer resp.Body.Close()

            body, err := ioutil.ReadAll(resp.Body)
            Panic(err)

            cntFound := strings.Count(string(body), searchString)
            cntAllFound += cntFound
            fmt.Printf("Count for %s: %d\n", url, cntFound)

			done <- true
		}(i)
	}

	// all url finish
	<-waitAllUrl
    fmt.Printf("Total: %d\n", cntAllFound)
}

func Panic(err error) {
    if err != nil {
        panic(err)
    }
}

func parseStdin() ([]string, int) {

    var (
        strUrls []byte
        arrUrls []string
    )

    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) != 0 {
        fmt.Println("no data")
        os.Exit(0)
    }

    strUrls, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        panic(err)
    }

    if len(strUrls) != 0 {
        arrUrls = strings.Split(string(strUrls), "\n")
    }

    return arrUrls, len(arrUrls) - 1
}