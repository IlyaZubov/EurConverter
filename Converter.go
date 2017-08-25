package main

import (
	"fmt"
	"net/http"
	"strconv"
	"flag"
	"io/ioutil"
	"encoding/xml"
)

const DEFAULT_CURRENCY = "unset"
const DEFAULT_AMOUNT = "1"
//EUR exchange rate here
const URL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type Envelope struct {
    XMLName     xml.Name `xml:"Envelope"`
    Cube        Cube `xml:"Cube"`
}

type Cube struct {
    XMLName     xml.Name `xml:"Cube"`
    TimeCube    TimeCube `xml:"Cube"`
}

type TimeCube struct {
    XMLName     xml.Name `xml:"Cube"`
    CurrencyCube    []CurrencyCube `xml:"Cube"`
}

type CurrencyCube struct {
    XMLName     xml.Name `xml:"Cube"`
    Currency    string `xml:"currency,attr"`
    Rate    	string `xml:"rate,attr"`
}

func strToFloat(input_str string) float64 {
	f, err := strconv.ParseFloat(input_str, 64)
	if (err == nil) {
		return f
	}
	fmt.Println("!!!Invalid amount")
	return 0
}

func floatToStr(inputNum float64) string {
    return strconv.FormatFloat(inputNum, 'g', 8, 64)
}

func getCurrencyCubes() []CurrencyCube {
	response, err := http.Get(URL)
	if err != nil {
		fmt.Errorf("GET error: %v", err)
		return nil
	} else {
		xmlResponse, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Errorf("Read body: %v", err)
			return nil
		}
		envelope := Envelope{}
		err = xml.Unmarshal(xmlResponse, &envelope)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		return envelope.Cube.TimeCube.CurrencyCube
	}
}

func getCurrencyRate(currency string) float64{
	Cubes := getCurrencyCubes()
	for _, v := range Cubes {
		if v.Currency == currency {
			return strToFloat(v.Rate)
		}
	}
	return 0
}

func convert(currency string, amount float64) string{
	rate := getCurrencyRate(currency)
	if rate == 0 {
		return "Not found currency. Can't convert to"
	}
	return floatToStr(amount/rate)
}

func parseFlags() (string, float64) {
	fmt.Println("available flags: \n	-currency=<currency>\n\r	-amount=<amount>\n")
	
	currency := flag.String("currency", DEFAULT_CURRENCY, "initial currency")
	amount := flag.String("amount", DEFAULT_AMOUNT, "amount to convert")
	
	flag.Parse()
	
	return *currency, strToFloat(*amount)
}

func main() {
	currency, amount := parseFlags()

	if currency == "unset" {
		fmt.Println("Choose currency by adding -currency=<CURRENCY> flag")	
	} else {
		fmt.Printf("Chosen currency: %v\n", currency)
		fmt.Printf("Chosen amount: %.3f\n", amount)
		fmt.Println("=================")
		fmt.Printf("%v %v = %v EUR\n", amount, currency, convert(currency, amount))
	}
}

