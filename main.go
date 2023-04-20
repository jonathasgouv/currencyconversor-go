package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Conversor struct {
    Amount float64
    Base string
    Date string
    Rates map[string]float64
    Initialized bool
}

func (c *Conversor) initialize() {
    c.Rates[c.Base] = c.Amount
    c.Initialized = true
}

func (c Conversor) convert(originCurrency string, targetCurrency string, amount float64) (float64, error) {
    if (!c.Initialized) {
        return 0, errors.New("you must initialize the conversor first")
    }

    originCurrencyRate := c.Rates[originCurrency]
    targetCurrencyRate := c.Rates[targetCurrency]

    result := (originCurrencyRate / targetCurrencyRate) * amount

    return result, nil
}

func main() {
    // Set up the HTTP client
    client := &http.Client{}

    // Make the request to the API
    req, err := http.NewRequest("GET", "https://api.frankfurter.app/latest", nil)
	
    if err != nil {
        fmt.Println(err)
        return
    }

    resp, err := client.Do(req)
	
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    // Parse the response body
    var conversor Conversor
    if err := json.NewDecoder(resp.Body).Decode(&conversor); err != nil {
        fmt.Println(err)
        return
    }

    conversor.initialize()

    fmt.Println("Hello, welcome to the JG's currency conversor.\n\nAll rates are base on the Frankfurter API (https://api.frankfurter.app)")

    for {
        var originCurrency int
        var targetCurrency int
        var amount float64
        var shouldBreak string

        fmt.Println("These are the available currencies:")

        availableOptions := make([]string, len(conversor.Rates))
        i:= 0

        for k := range conversor.Rates {
            fmt.Printf("%v - %s\n", i + 1, k)

            availableOptions[i] = k

            i++
        }

        fmt.Println("Enter the origin currency: ")

        for {
            fmt.Scan(&originCurrency)

            if originCurrency >= 1 && originCurrency < len(availableOptions) {
                break
            }

            fmt.Println("\nInvalid option! Enter the origin currency:")
        }

        fmt.Println("Enter the yarget currency (Can't be the same as origin): ")


        for {
            fmt.Scan(&targetCurrency)

            if targetCurrency >= 1 && targetCurrency < len(availableOptions) && targetCurrency != originCurrency {
                break
            }

            fmt.Println("\nInvalid option! Enter the target currency:")
        }

        fmt.Println("Enter the amount: ")
        fmt.Scan(&amount)



        result, err := conversor.convert(availableOptions[originCurrency-1], availableOptions[targetCurrency-1], amount)

        if err != nil {
            fmt.Println("There was some error, restarting process")

            continue
        }


        fmt.Printf("\n%v %s is equal to %v %s", amount, availableOptions[originCurrency-1], result, availableOptions[targetCurrency-1])

        fmt.Println("\n\nPress X to end the program and anything else to coninue")
        fmt.Scan(&shouldBreak)

        if strings.ToUpper(shouldBreak) == "X" {
            break
        }
    }
}
