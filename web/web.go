package web

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/bluesNbrews/ParseWebPage/link"
)

var counter int = 0

//Gethtml makes a GET call to a URL and returns the HTML body
func Gethtml(enteredurl string) io.Reader {

	resp, err := http.Get(enteredurl)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	htmlcontent := bytes.NewReader(body)

	return htmlcontent

}

//Get http respone code for each link passed in. Use channel to pass back response code,
//which will be concurrently assigned and printed via UpdateAndPrint
func GetUrlStatus(newlinks link.Link, c chan int) {

	//Get the response or error from GET request. They are both mutually exclusive (returns one or the other)
	resp, err := http.Get(newlinks.Href)

	//If there is an error, the error will be logged and the program will exit
	//Lastly, the error will be displayed after the program closes
	//There may be other functions to replace log.Fatal, but I had a hard time finding and implementing one
	if err != nil {

		log.Fatal(err)

	}
	defer resp.Body.Close()

	//Pass the status code via channel
	c <- resp.StatusCode

}

//UpdateAndPrint updates the status code for the new links and prints the output
func UpdateAndPrint(newlinks link.Link, c chan int, codestable map[string]int) {

	//Receive and assign status code via channel (waits for send)
	newlinks.Code = <-c

	//Print table rows
	fmt.Printf("| %-3d | %-120s | %-60s | %-3d |\n", counter, newlinks.Href, newlinks.Text, newlinks.Code)
	counter++

	//Increment the HTTP return code counter by 1
	codestable[strconv.Itoa(newlinks.Code)]++

}
