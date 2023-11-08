package template

import (
	"fmt"
    "io"
    "os"
	"github.com/calebcscott/pkg-init/pkg/config"
)


type TemplateHandler struct {
    Name            string
    language        string
    directory       string
    template        template
}


func IsEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }

    defer f.Close()

    _, err = f.Readdirnames(1)
    if err == io.EOF {
        return true, nil
    }

    return false, err
}

func getChoiceB(msg string) bool {

    for {
        fmt.Print(msg)

        var choice string
        fmt.Scan(&choice)
        switch  choice {
        case "Y", "y":
            return true
        case "N", "n":
            return false
        default:
            msg = "Must select yes or not [Y/n]: "
        }
    }
}

func (th *TemplateHandler) Init(config *config.PkgConfig) error {
    fmt.Println("Initializing template", th.Name, "in directory", th.directory)

    empty, err := IsEmpty(th.directory)

    if err == nil && !empty {
        // need to check if user is Ok with starting in non-empty directory
        msg := fmt.Sprintf("Directory %s is not empty, do you want to continue [Y/n]: ", th.directory)
        if !getChoiceB(msg) {
            return nil
        }
    }

    if err := th.template.build(config, th.directory); err != nil {
        return err
    }

    return nil
}


func NewTemplate(name string, lang string, dir string, config *config.PkgConfig) (TemplateHandler, error) {
    var template template
    var err error
    if lang != "" {
        template, err = findLangTemplate(lang, config)
        if err != nil {
            fmt.Println(err)
            template, err = readTeamplate(name, config)
        } 
    } else {
        template, err = readTeamplate(name, config)
    }

    

    if err != nil {
        return TemplateHandler{}, err
    }

    return TemplateHandler{
        Name: name,
        language: lang,
        directory: dir,
        template: template,
    }, nil
}


