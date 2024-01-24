package td

import (
	"errors"
	"os"
	"strings"
	"time"

	"git.burning.moe/celediel/burning.moe/internal/models"

	"github.com/ilyakaznacheev/cleanenv"
)

const dataDir string = "./templates/data/"

// in order of precedence
var dataExtensions = [4]string{"yml", "yaml", "toml", "json"}

// makeBasicTemplateData creates a blank TemplateData
// containing only the time the template was generated
func MakeBasicTemplateData(when time.Time) models.TemplateData {
	strMap := map[string]string{
		"GeneratedAt": when.Format(time.UnixDate),
	}

	templateData := models.TemplateData{
		StringMap: strMap,
	}
	return templateData
}

// LoadTemplateData loads template data from file. If that
// fails, it returns an empty TemplateData and an error
func LoadTemplateData(page string) (models.TemplateData, error) {
	var data models.TemplateData
	output := dataDir + strings.ReplaceAll(page, ".tmpl", "")

	for _, extension := range dataExtensions {
		if info, err := os.Stat(output + "." + extension); err == nil && !info.IsDir() {
			err = cleanenv.ReadConfig(output+"."+extension, &data)
			if err == nil {
				// don't try anymore files, but do setup an empty StringMap
				if len(data.StringMap) == 0 {
					data.StringMap = make(map[string]string)
				}
				return data, nil
			}
		}
	}

	// couldn't load anything from file
	return models.TemplateData{}, errors.New("Couldn't load data from file")
}
