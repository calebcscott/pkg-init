package template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/calebcscott/pkg-init/pkg/config"
	"github.com/spf13/viper"
	//"gopkg.in/yaml.v3"
)


type template interface {
    parse( *config.PkgConfig, string ) error
}

type templateContent struct {
    contents map[string]interface{}
}


func validate_template(contentMap interface{}, tld string) error {
    switch  contentMap := contentMap.(type) {
    case string:
    case []interface{}:
    case map[string]interface{}:
        for dirName, contents := range contentMap {
            newTld := tld + "/" + dirName
            
            validate_template(contents, newTld)
        }
    default:
        return errors.New("malformed template")
    }

    return nil
}


func build_template(contentMap interface{}, tld string) {
    switch  contentMap := contentMap.(type) {
    case string:
        // create file and add contents
        // filename will be tld, contents will be contentMap
    case []interface{}:
        // empty dir, just need to create it
    case map[string]interface{}:
        // possible dir/subdir/file
        for dirName, contents := range contentMap {
            newTld := tld + "/" + dirName

            // if match on switch create directory
            switch contents.(type) {
            case map[string]interface{}:
                    // if subdirs or files, need to create this dir
                    fmt.Println("Got dir to make", newTld)
                
            }

            // recurse to pick up files/subdirs
            build_template(contents, newTld)
        }
    }


}

func (t templateContent) parse( config *config.PkgConfig, tld string ) error {
    fmt.Println("Parsing template" ) 

    // first parse to validate
    if err := validate_template(t.contents, tld); err != nil {
        return err
    }

    build_template(t.contents, tld)

    return nil
}



func readTeamplate(templateName string, config *config.PkgConfig) (template, error) {

    if res := strings.Contains(templateName, "yaml"); !res {
        // check config for template
        v, found := config.TemplateMap[templateName]

        if !found {
            return nil, errors.New("No template with name: "+templateName+"\n\tTry adding it.")
        }

        templateMap := viper.GetStringMap(v)

        t, found := templateMap["type"] 
        if !found {
            return nil, errors.New("malformed template, couldn't find type")
        }

        switch t {
        case "raw":
            if contents, found := templateMap["contents"]; found {
                return templateContent{ contents.(map[string]interface{}) }, nil
            } else {
                return nil, errors.New("malformed template content")
            }
        }
    }


    return nil, errors.New("No template with name: "+templateName)
}
