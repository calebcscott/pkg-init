package template

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/calebcscott/pkg-init/pkg/config"
	"github.com/spf13/viper"
	//"gopkg.in/yaml.v3"
)

// main interface for interacting with template
type template interface {
    build( *config.PkgConfig, string ) error
}


type cmd struct {
    cmds []exec.Cmd
}

type templateContent struct {
    contents map[string]interface{}
    commands cmd
}

type templatePath struct {
    path string
    commands cmd
}

type templateGit struct {
    repo string
    commands cmd
}


func newCommand(cmds string) (cmd, error) {

    var commands []exec.Cmd
    
    for _, line := range strings.Split(cmds, "\n") {
        if line == "" {
            continue
        }
        cmdString := strings.Split(line, " ")
        commands = append(commands, *exec.Command(cmdString[0], cmdString[1:]...))

    }

    return cmd{ commands }, nil
}
func (c cmd) run (dir string) error {

    if len(c.cmds) == 0 {
        return nil
    }

    for _, command := range c.cmds {
        command.Dir = dir
        var out strings.Builder
        command.Stdout = &out
        err := command.Run()

        if err != nil {
            log.Fatal(err)
        }

        fmt.Print(out.String())
    }

    return nil
}

func newTemplateContent(temp map[string]interface{}) (templateContent, error) {

    // pull out contents
    contents, found := temp["contents"]

    if !found {
        return templateContent{}, errors.New("malformed template, no contents found")
    }
    if err:= validate_template(contents); err != nil {
        return templateContent{}, err
    }

    // pull out commands, optional field
    commands := temp["commands"]

    if commands == nil {
        commands = ""
    }
    cmd, err := newCommand(commands.(string))
    if err != nil {
        return templateContent{}, err
    }


    return templateContent{contents.(map[string]interface{}), cmd }, nil
}


/*
    validate_template

    Validate whether supplied template matches expected contents
    Do not need to pass a Top-Level-Directory / Prefix since we
        do not do any path validation, only template validation
*/
func validate_template(contentMap interface{}) error {
    switch  contentMap := contentMap.(type) {
    // empty case for string []interface{} since we expect those
    case string:
    case []interface{}:
    case map[string]interface{}:
        for _, contents := range contentMap {
            validate_template(contents)
        }
    default:
        return errors.New("malformed template")
    }

    return nil
}

/*
    build_template

    Still need to return possible errors with directory or file names/creation
*/
func build_template(contentMap interface{}, tld string) error {
    var err error = nil
    switch  contentMap := contentMap.(type) {
    case string:
        // create file and add contents
        // filename will be tld, contents will be contentMap
        err = os.WriteFile(tld, []byte(contentMap), 0644 )
        if err != nil {
            fmt.Println("Could not create file:", tld)
        }
    case []interface{}:
        // empty dir, just need to create it
        err = os.MkdirAll(tld, 0755)
    case map[string]interface{}:
        // possible dir/subdir/file
        for dirName, contents := range contentMap {
            // appends next path name with OS specific separator
            newTld := filepath.Join(tld, dirName)
            var dirErr error = nil
            // if match on switch create directory
            switch contents.(type) {
            case map[string]interface{}:
                // if subdirs or files, need to create this dir
                fmt.Println("Making directory:", newTld)
                dirErr = os.MkdirAll(newTld, 0755)
                if err == nil {
                    err = dirErr
                }
            }

            // recurse to pick up files/subdirs if no error
            // cannot recurse if error'd on dir creation
            if dirErr == nil {
                subErr := build_template(contents, newTld)

                if err == nil {
                    err = subErr
                }
            } 
        }
    }

    return err
}

func (t templateContent) build( config *config.PkgConfig, tld string ) error {
    // attempt to build template
    if err := build_template(t.contents, tld); err != nil {
        return err
    }


    if err := t.commands.run(tld); err != nil {
        return err
    }

    return nil
}

func findLangTemplate(lang string, config *config.PkgConfig) (template, error) {
    if lang == "" {
        return nil, errors.New("no language provided")
    }
    
    templateName, found := config.LanguageSupport[lang]

    if !found {
        return nil, errors.New("language ("+lang+") not supported")
    }

    return readTeamplate(templateName, config)
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
            return newTemplateContent(templateMap)
        }
    }


    return nil, errors.New("No template with name: "+templateName)
}
