package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//CurrentUpdateCenter is the data wrapper.
type CurrentUpdateCenter struct {
	Core    Core
	Plugins Plugins
}

//Core is the built core data carrier.
type Core struct {
	BuildDate string
	Name      string
	Sha1      string `json:"sha1"`
	URL       string
	Version   string
}

//Plugins is list of available plugins.
type Plugins []Plugin

//UnmarshalJSON deserializes a nestled json data plugin.
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

//Plugin is a jenkins plugin data carrier.
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

//WithoutJSONP parses raw http://mirror.xmission.com/jenkins/updates/current/update-center.json request.
func WithoutJSONP(w http.ResponseWriter, r *http.Request) {
	//url := "http://localhost:3000/jenkins/updates/current/update-center.json"
	url := "http://mirror.xmission.com/jenkins/updates/current/update-center.json"

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

//Parse parses the CurrentUpdateCenter data struct.
func Parse() (*CurrentUpdateCenter, error) {
	url := "http://localhost:3000/ParseAwayJSONP"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	log.Println("Fetched data from:", url)
	decoder := json.NewDecoder(response.Body)

	var data CurrentUpdateCenter
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(v)
		}
		return nil, err
	}
	return &data, nil
}
