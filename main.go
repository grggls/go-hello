// our package is called 'main' -- so go.mod for more info
package main

// import pulls in an external module
import (
	"encoding/json"
	"net/http"
	"strings"
)

// main() reserved func name for executable packages -- entrypoint
func main() {
	http.HandleFunc("/hello", hello)

	// define this handler func inline
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		// strings.SplitN takes everything in the path after '/weather/' and puts it in 'city' ``
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		// if there's an error calling query, propogate that error vi http.Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// query was successful. tell the client we're returning json data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// use json.NewEncoder to JSON-encode the weatherData directly
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

// declare an http.HandlerFunc, (has a specific type signature, or implements the interface of ...)
// so can be passed as an argument to HandleFunc
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

/*
 * define a new type with "type" keyword
 * declare it a struct
 * each field gets a name (Name, Main, Kelvin)
 * and a type (string, float64, an inline struct [called Main])
 * the `json:"foo"` bits are called 'tags' and they're metadata or attributes
 *   they allow us to use the encoding/json package to unmarshall the API's
 *   responses, giving us the benefits of type safety when using a 3rd party API response
 */
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

// takes a string representing the city, and returns a weatherData struct and an error
func query(city string) (weatherData, error) {
	// fetch weather data from openweathermap using our 'city' string and the api key we requested
	Apikey := "28d2fc80b71bd20d670acf2326ad0b84"
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + Apikey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	// resource  mgmt - if the http.Get has succeeded, defer a call to close the response Body
	defer resp.Body.Close()

	// create a weatherData struct
	var d weatherData

	// use json.NewDecoder to unmarshall the API response into a wweatherData object
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	// return the weatherData to the caller, with a nil error to indicate success.
	return d, nil
}
