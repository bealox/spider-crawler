package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("http://www.dogsvictoria.org.au/DogsPuppies/BuyingAPuppy/Breedersdirectory.aspx")

	fmt.Println("http tranport error is ", err)

	body, err := ioutil.ReadAll((resp.Body))

	fmt.Println("read error is ", err)

	fmt.Println(string(body))
}
