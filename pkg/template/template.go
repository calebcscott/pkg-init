package template

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/calebcscott/pkg-init/pkg/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	//"gopkg.in/yaml.v3"
)

// main interface for interacting with template
type template interface {
    build( *config.PkgConfig, string ) error
}

type File struct {
    name    string
    content string
}
type Dir struct {
    name    string
    files []File
}
type TemplateContents struct {
    dirs []Dir
    files []File
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


/*
    TODO: maybe figure out how to create custom unmarshaler for a template contents object
        could be useful(?)
*/
func (tc *TemplateContents) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var details map[string]interface{}

    if err := unmarshal(&details); err != nil {
        return err
    }

    for entry, content := range details {
        switch content := content.(type) {

        case string:
            file := File{ entry, content }
            tc.files = append(tc.files, file)

        case []interface{}:
            var dir Dir
            dir.name = entry
            tc.dirs = append(tc.dirs, dir)

        case map[string]interface{}:
            //for subEntry, subContents := range content {
            //}
        }
    }


    return nil
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

func (c *cmd) run (dir string) error {

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


func newTemplatePath(temp map[string]interface{}) (*templatePath, error) {

    possiblePath, found := temp["path"]

    if !found {
        return nil, errors.New("malformed template, no path found")
    }

    path, ok := possiblePath.(string)
    if !ok {
        return nil, errors.New("malformed value for path, expected string")
    }

    _, err := os.Stat(path)
    if err != nil {
        return nil, errors.New("cannot find provided path in template")
    } 

    // pull out commands, optional field
    commands := temp["commands"]

    if commands == nil {
        commands = ""
    }
    cmd, err := newCommand(commands.(string))
    if err != nil {
        return nil, err
    }

    return &templatePath{ path, cmd }, nil
}

func newTemplateContent(temp map[string]interface{}) (*templateContent, error) {

    // pull out contents
    contents, found := temp["contents"]

    if !found {
        return nil, errors.New("malformed template, no contents found")
    }
    if err:= validateTemplateContents(contents); err != nil {
        return nil, err
    }

    // pull out commands, optional field
    commands := temp["commands"]

    if commands == nil {
        commands = ""
    }
    cmd, err := newCommand(commands.(string))
    if err != nil {
        return nil, err
    }


    return &templateContent{contents.(map[string]interface{}), cmd }, nil
}


/*
    validateTemplateContents

    Validate whether supplied template matches expected contents
    Do not need to pass a Top-Level-Directory / Prefix since we
        do not do any path validation, only template validation
*/
func validateTemplateContents(contentMap interface{}) error {
    switch  contentMap := contentMap.(type) {
    // empty case for string []interface{} since we expect those
    case string:
    case []interface{}:
    case map[string]interface{}:
        for _, contents := range contentMap {
            validateTemplateContents(contents)
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
func buildTemplateContents(contentMap interface{}, tld string) error {
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
                subErr := buildTemplateContents(contents, newTld)

                if err == nil {
                    err = subErr
                }
            } 
        }
    }

    return err
}


func (t *templatePath) build( config *config.PkgConfig, tld string ) error {
    // attempt to create dir/file in provided path as check(?)


    return nil
}

func (t *templateContent) build( config *config.PkgConfig, tld string ) error {
    // attempt to build template
    if err := buildTemplateContents(t.contents, tld); err != nil {
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

/*
    TODO: @ccs change signature or add procedure to provide template name
        may not want template name == filename
*/
func readYamlFile(fileName string) (map[string]interface{}, error) {
    fd, err := os.Open(fileName)

    if err != nil {
        return nil, errors.New("Could not find/open tempalte file: " + fileName)
    }

    stat, err := fd.Stat()
    if err != nil {
        return nil, errors.New("Could not read tempalte file: " + fileName)
    }


    yamlFile := make([]byte, stat.Size())
    _, err = bufio.NewReader(fd).Read(yamlFile)

    if err != nil && err != io.EOF {
        return nil, errors.New("Could not read tempalte file: " + fileName)
    }
    templateMap := make(map[string]interface{})

    err = yaml.Unmarshal(yamlFile, templateMap)

    if err != nil {
        return nil, errors.New("Could not parse yaml file: " + fileName)
    }

    _, templateName := filepath.Split(fileName)
    templateName = strings.Split(templateName, ".")[0]

    var found bool
    templateMap, found = templateMap[templateName].(map[string]interface{})

    if !found {
        return nil, errors.New("malformed template, could not entry: " + templateName)
    }

    return templateMap, nil
}

func readTeamplate(templateName string, config *config.PkgConfig) (template, error) {

    var templateMap map[string]interface{}
    var err error

    // cant enter this block if we were not handed a ref to the config
    // if config is nil and templateName is a valid path to valid template yaml we are fine
    if res := strings.Contains(templateName, "yaml"); !res && config != nil {
        // check config for template
        v, found := config.TemplateMap[templateName]

        if !found {
            return nil, errors.New("No template with name: "+templateName+"\n\tTry adding it.")
        }


        if res := strings.Contains(v, "yaml"); !res {
            // read template from config
            templateMap = viper.GetStringMap(v)
        } else {
            // read template from yaml file directly
            templateMap, err = readYamlFile(v) 

            if err != nil {
                return nil, err
            }
        }
        
    } else {
        // read yaml file directly for templateMap
        templateMap, err = readYamlFile(templateName) 

        if err != nil {
            return nil, err
        }

    }

    // switch on type to perform necessary steps for template struct building
    t, found := templateMap["type"] 
    if !found {
        return nil, errors.New("malformed template, couldn't find type")
    }

    switch t {
    case "raw":
        return newTemplateContent(templateMap)
    // TODO: @ccs add additional template types 
    default:
        return nil, errors.New("Template type("+t.(string)+") not implemented.")
    }

}
