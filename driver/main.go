package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bluesNbrews/ParseWebPage/link"
	"github.com/bluesNbrews/ParseWebPage/web"
)

func main() {

	//Retrieve URL string from input
	var enteredurl = string(os.Args[1])

	//Channel for status code
	var c = make(chan int, 100)

	//Call the URL string and retrieve HTML content as io.Reader
	htmlcontent := web.Gethtml(enteredurl)

	//Parse htmlcontent and return an array of <a> tags found in it
	links, err := link.Parse(htmlcontent)
	if err != nil {
		panic(err)
	}

	//Read links array, add domain name to hrefs where missing, then place in new array
	newlinks := link.Fixlinks(links, enteredurl)

	//Create hash table to count the occurance of the various HTTP return codes
	var codestable map[string]int
	codestable = make(map[string]int)

	//Process each link concurrently to get http status code and assign to link and print
	//The sleep is used temporarily to prevent too many http requests at one time
	for i := 0; i < len(newlinks); i++ {

		go web.GetUrlStatus(newlinks[i], c)
		go web.UpdateAndPrint(newlinks[i], c, codestable)
		time.Sleep(125 * time.Millisecond)

	}

	//Print the content of the HTTP return codes hash table
	fmt.Println("")
	for key, value := range codestable {
		fmt.Println("Return code:", key, "Number of occurances:", value)
	}

}
