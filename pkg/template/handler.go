package template

import (
	"fmt"

	"github.com/calebcscott/pkg-init/pkg/config"
)


type TemplateHandler struct {
    Name            string
    language        string
    directory       string
    template        template
}


func (th TemplateHandler) Init(config *config.PkgConfig) error {
    fmt.Println("Initializing template", th.Name, "in directory", th.directory)

    th.template.parse(config, th.directory)

    return nil
}


func NewTemplate(name string, lang string, dir string, config *config.PkgConfig) (TemplateHandler, error) {
    template, error := readTeamplate(name, config)

    if error != nil {
        return TemplateHandler{}, error
    }

    fmt.Println(template)
    return TemplateHandler{
        Name: name,
        language: lang,
        directory: dir,
        template: template,
    }, nil
}


