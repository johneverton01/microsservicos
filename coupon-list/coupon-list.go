package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)



type Coupons struct {
	Coupons []string
}

func Check(coupons []string , code string) string {
	for i := range coupons {
		 if code == coupons[i] {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var coupons []string

func main() {

	coupon := []string{
		"abc",
		"aaa",
		"bbb",
		"ccc",
		}

	coupons = append(coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid := Check(coupons, coupon)

	result := Result{Status: valid}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}