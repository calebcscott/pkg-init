package template

import (
    "testing"
)





func Test_validateTemplate_BadInput(t *testing.T) {
    var emptyMap interface{}

    err := validateTemplate(emptyMap)

    if err == nil {
        t.Errorf("validateTemplate on empty map failed, expected error")
    }
}



func Test_newTemplateContent_BadInput(t *testing.T) {
    var emptyMap map[string]interface{}

    tc, err := newTemplateContent(emptyMap)

    if len(tc.contents) != 0 || err == nil {
        t.Errorf("newTemplateContent failed, expected -> %v, got -> %v", templateContent{}, tc)
    }
}
