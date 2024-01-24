package models

import (
	"encoding/json"
	"errors"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

const dataDir string = "./templates/data/"

// in order of precedence
var dataExtensions = [4]string{"yml", "yaml", "toml", "json"}

type Link struct {
	Href, Icon template.URL
	Text       string
}

// TemplateData holds data sent from handlers to templates.
type TemplateData struct {
	StringMap map[string]string      `json:"StringMap" yaml:"StringMap" toml:"StringMap"`
	IntMap    map[string]int         `json:"IntMap" yaml:"IntMap" toml:"IntMap"`
	FloatMap  map[string]float32     `json:"FloatMap" yaml:"FloatMap" toml:"FloatMap"`
	LinkMap   map[string][]Link      `json:"LinkMap" yaml:"LinkMap" toml:"LinkMap"`
	Data      map[string]interface{} `json:"Data" yaml:"Data" toml:"Data"`
	CSRFToken string                 `json:"Csrftoken" yaml:"Csrftoken" toml:"Csrftoken"`
	Flash     string                 `json:"Flash" yaml:"Flash" toml:"Flash"`
	Warning   string                 `json:"Warning" yaml:"Warning" toml:"Warning"`
	Error     string                 `json:"Error" yaml:"Error" toml:"Error"`
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
	var data TemplateData
	output := dataDir + strings.ReplaceAll(page, "tmpl", "")

	for _, extension := range dataExtensions {
		if info, err := os.Stat(output + extension); err == nil && !info.IsDir() {
			data, err = loadFromFile(output, extension)
			if err == nil {
				// don't try anymore files
				return data, nil
			}
		}
	}

	// couldn't load anything from file
	return TemplateData{}, errors.New("Couldn't load data from file")
}

// loadFromFile loads TemplateData from the specified filetype (yaml, toml, or json)
func loadFromFile(filename, filetype string) (TemplateData, error) {
	var data TemplateData
	file, err := os.ReadFile(filename + filetype)
	if err != nil {
		return TemplateData{}, err
	}

	switch filetype {
	case "json":
		err = json.Unmarshal(file, &data)
		if err != nil {
			return TemplateData{}, err
		}
	case "toml":
		err = toml.Unmarshal(file, &data)
		if err != nil {
			return TemplateData{}, err
		}
	case "yaml":
		err = yaml.Unmarshal(file, &data)
		if err != nil {
			return TemplateData{}, err
		}
	}
	return data, nil
}
