package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type CurrentUpdateCenter struct {
	Core    Core
	Plugins Plugins
}

type Core struct {
	BuildDate string
	Name      string
	Sha1      string `json:"sha1"`
	URL       string
	Version   string
}

type Plugins []Plugin

func (p *Plugins) UnmarshalJSON(b []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	for k, _ := range m {
		var plugin Plugin
		err = json.Unmarshal(*m[k], &plugin)
		if err != nil {
			return err
		}
		*p = append(*p, plugin)
	}

	return nil
}

type Plugin struct {
	Name             string
	ReleaseTimeStamp string
	RequiredCore     string
	SCM              string
	Title            string
	URL              string
	Version          string
	Wiki             string
}

func handler(w http.ResponseWriter, r *http.Request) {
	d, err := os.Open("update-center.json")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer d.Close()

	_, err = io.Copy(w, d)
	if err != nil {
		http.Error(w, "Internal error while writing response: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlerParse(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:3000/ParseAwayJSONP"
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var data CurrentUpdateCenter
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(v)
		}
	}
	md, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal error while writing response: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", md)
}
func main() {
	port := 3000
	mux := http.NewServeMux()

	mux.HandleFunc("/jenkins/updates/current/update-center.json", handler)
	mux.HandleFunc("/ParseAwayJSONP", ParseAwayJSONP)
	mux.HandleFunc("/parse", handlerParse)

	log.Println("Server to listen on a port: ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))

}

//ParseAwayJSONP parses raw http://mirror.xmission.com/jenkins/updates/current/update-center.json request.
func ParseAwayJSONP(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:3000/jenkins/updates/current/update-center.json"
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	jsonp := string(b)
	rmFirst := "updateCenter.post(\n"
	rmLast := ");\n"
	json := jsonp[len(rmFirst) : len(jsonp)-len(rmLast)]
	fmt.Fprintf(w, "%s", json)

}
