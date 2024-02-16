package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Scene struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type storyHandler struct {
	story map[string]Scene
}

func (s *storyHandler) renderScene(scene Scene, w http.ResponseWriter) {
	tmpl := template.Must(template.New("tpl").Parse(`{{.Title}}`))
	tmpl.Execute(w, scene)
}

func (s *storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arc := r.URL.Query().Get("arc")
	// Or if arc not in s.story
	if arc == "" {
		arc = "intro"
	}

	s.renderScene(s.story[arc], w)
}

func readJson() (map[string]Scene, error) {
	fileName := "gopher.json"
	file, err := os.ReadFile(fileName)
	res := make(map[string]Scene)

	if err != nil {
		return res, err
	}

	err = json.Unmarshal(file, &res)

	return res, err
}

func startServer(story map[string]Scene) {
	mux := http.NewServeMux()
	addr := ":8080"
	mux.Handle("/story", &storyHandler{
		story: story,
	})
	fmt.Println("Server listening on port", addr)
	http.ListenAndServe(addr, mux)
}

func main() {
	story, err := readJson()

	if err != nil {
		log.Fatalln("Error reading JSON", err.Error())
	}

	startServer(story)
}
