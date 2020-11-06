package main

import (
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9091", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")
	result := Result{Status: "declined"}

	resultCoupon := makeHttpCall("http://localhost:9092", coupon)

	if resultCoupon.Status == "valid" {
		if ccNumber == "1" {
			result.Status = "approved"
		}
	}

	if resultCoupon.Status == "invalid" {
		result.Status = "invalid coupon"
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}


func makeHttpCall(urlMicroservice string, coupon string) Result {

	values := url.Values{}
	values.Add("coupon", coupon)

	if coupon == "" {
		result := Result{Status:"Coupon not informed"}
		return result
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result

}
