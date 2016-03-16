package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	cache "jenkins-plugins/cache"
	"jenkins-plugins/parse"
)

// Handler handles http calls.
type Handler struct {
	Cache *cache.Cache
}

//NewHandler handles initzialisation of Handler.
func NewHandler(c *cache.Cache) Handler {
	return Handler{Cache: c}
}

//Index returns available plugins
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	plugins := h.getAllPlugins()
	title := "Jenkins_plugins GO!"
	resources := ""
	for _, p := range *plugins {
		resources = resources + fmt.Sprintf(`<div class="container" style="border: 1px solid black">
																					<div> <b>%s</b> <a href="%s">%s</a> </div>
																					<div>%s</div>
																				</div>`, p.Name, p.Title, p.URL, p.Wiki)
	}
	headernote := `<b>See also json at</b> <a href="/Plugins">/Plugins</a>`
	content := fmt.Sprintf(`<h1>%s</h1><div>%s</div>`, title, resources)
	page := "<html><head>" + headernote + "</head>" + "<body>" + content + "</body></html>"
	fmt.Fprintf(w, "%s", page)

}

//Plugins returns available plugins as json.
func (h *Handler) Plugins(w http.ResponseWriter, r *http.Request) {
	plugins := h.getAllPlugins()
	json.NewEncoder(w).Encode(plugins)
}

//Categories returns available Categories available plugins.
func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	labels := h.getAllPluginLabels()
	json.NewEncoder(w).Encode(labels)
}

func (h *Handler) getAllPlugins() *parse.Plugins {
	var plugins parse.Plugins
	keys := h.Cache.GetAll()
	fmt.Printf("Loaded all %v items keys from cache. \n", len(keys))
	for _, key := range keys {
		//fmt.Printf("Loading item with key %s from cache. \n", key)
		var plugin parse.Plugin
		err := json.Unmarshal(h.Cache.Get(key), &plugin)
		if err != nil {
			continue
		} else {
			plugins = append(plugins, plugin)
		}
	}
	return &plugins
}

func (h *Handler) getAllPluginLabels() *[]string {
	var labels []string
	keys := h.Cache.GetAll()
	fmt.Printf("Loaded all %v items keys from cache. \n", len(keys))
	for _, key := range keys {
		//fmt.Printf("Loading item with key %s from cache. \n", key)
		var plugin parse.Plugin
		err := json.Unmarshal(h.Cache.Get(key), &plugin)
		if err != nil {
			continue
		} else {
			for _, l := range plugin.Labels {
				if !contains(l, labels) {
					labels = append(labels, l)
				}
			}

		}
	}
	return &labels
}

//Checks if item exist in list.
func contains(i string, l []string) bool {
	for _, a := range l {
		if i == a {
			return true
		}
	}
	return false
}

//UpdateCache updates the handler cache with information from external source.
func (h *Handler) UpdateCache(w http.ResponseWriter, r *http.Request) {
	h.PerformUpdateCache()
}

//PerformUpdateCache executes a cache pull from external source into cache.
func (h *Handler) PerformUpdateCache() {
	currentUpdateCenter, err := parse.Parse()
	core, err := json.Marshal(currentUpdateCenter.Core)
	if err != nil {
		log.Println("Could not update cache, failed with error:", err)
		return
	}
	log.Println("Updated cache with core:", currentUpdateCenter.Core)
	h.Cache.Set("core", core)

	for _, plugin := range currentUpdateCenter.Plugins {
		bp, err := json.Marshal(plugin)
		if err != nil {
			log.Println("Could not update cache, failed with error:", err)
			return
		}
		h.Cache.Set(plugin.Name, bp)
	}
	log.Printf("Updated cache with %v plugins.\n", len(currentUpdateCenter.Plugins))
}

//DummyUpdateCenterFile returns dummy file.
func (h *Handler) DummyUpdateCenterFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

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

	log.Printf(
		"%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		time.Since(start),
	)
}
