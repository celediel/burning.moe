package models

import (
	"encoding/json"
	"html/template"
	"os"
	"strings"
	"time"
)

const dataDir string = "./templates/data/"

type Link struct {
	Href template.URL
	Text string
	Icon template.URL
}

// TemplateData holds data sent from handlers to templates.
type TemplateData struct {
	StringMap map[string]string      `json:"stringMap"`
	IntMap    map[string]int         `json:"intMap"`
	FloatMap  map[string]float32     `json:"floatMap"`
	LinkMap   map[string][]Link      `json:"linkMap"`
	Data      map[string]interface{} `json:"data"`
	CSRFToken string                 `json:"csrftoken"`
	Flash     string                 `json:"flash"`
	Warning   string                 `json:"warning"`
	Error     string                 `json:"error"`
}

// makeBasicTemplateData creates a blank TemplateData containing only the
// time the related template was generated
func MakeBasicTemplateData(when time.Time) TemplateData {
	strMap := map[string]string{
		"GeneratedAt": when.Format(time.UnixDate),
	}

	templateData := TemplateData{
		StringMap: strMap,
	}
	return templateData
}

// LoadTemplateData loads template data from file. If that
// fails, it returns an empty TemplateData and an error
func LoadTemplateData(page string) (TemplateData, error) {
	// TODO: load from something other than JSON
	var data *TemplateData
	output := dataDir + strings.ReplaceAll(page, "tmpl", "json")

	file, err := os.ReadFile(output)
	if err != nil {
		return TemplateData{}, err
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return TemplateData{}, err
	}

	return *data, nil
}
