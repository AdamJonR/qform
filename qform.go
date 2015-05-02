package qform

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/adamjonr/dialects"
)

// DSL struct provides the definition of the dialect
type DSL struct{}

type Attribute struct {
	Name  string
	Value string
}

type Field struct {
	Name       string
	Label      string
	InputType  string
	ID         string
	Attributes map[string]string
	Options    map[string]string
}

type Model struct {
	Attributes []Attribute
	Fields     []Field
}

const Indent = "  "

// NewDialect returns the FastFormsDialect struct for parsing of input
func (*DSL) NewDialect() (dialect *dialects.Dialect) {
	dialect = &dialects.Dialect{
		Title:       "Fast Forms",
		Description: "The Fast Forms DSL speeds the creation of HTML5 forms, often cutting the number of characters required in half.",
		RootName:    "form",
		Version:     1.0,
		Examples: map[string]string{
			"Basic Contact Form": `- method post

text
- name name
- maxlength 30
- required

email
- name email

textarea
- name message

submit
- value Send message`,
			"Option Fields": `radio
- name preference
- [
  call Call me back
  email Email me a message
  mail Send me a letter
]

checkbox
- name permission
- [
  yes I give my permission to contact me
]

select
- name department
- [
  sales
  tech Tech Support
  receivables
]`,
		},
		PartDefinitions: map[string]dialects.PartDefinition{
			"form": {
				Description:  "Composed of zero-or-more attributes and zero-or-more fields.",
				Constituents: [][]string{{"form attribute*", "form field*"}},
			},
			"form attribute": {
				Description:  "Defines attribute of the form tag.",
				Constituents: [][]string{{"hyphen", "name", "value?", "newline"}},
				Handler: func(part *dialects.Part, any interface{}) (ok bool) {
					model, ok := any.(*Model)
					if !ok {
						return false
					}
					name := part.Constituents[0].Value
					value := name

					if len(part.Constituents) == 2 {
						value = part.Constituents[1].Value
					}

					model.Attributes = append(model.Attributes, Attribute{Name: name, Value: value})
					return true
				},
			},
			"form field": {
				Description:  " Composed of optional new-line, field type, and zero-or-more field attributes.",
				Constituents: [][]string{{"newline?", "field type", "field attribute*"}},
				Handler: func(part *dialects.Part, any interface{}) (ok bool) {
					model, ok := any.(*Model)
					if !ok {
						return false
					}
					// use the first constituent of the first constituent to set the inputType
					field := Field{
						InputType:  part.Constituents[0].Constituents[0].Value,
						Attributes: make(map[string]string),
						Options:    make(map[string]string),
					}
					// cycle through the field attributes
					for i, length := 1, len(part.Constituents); i < length; i = i + 1 {
						// store constituent
						constituent := part.Constituents[i]
						// check attribute type
						if constituent.Constituents[0].Name == "name" {
							// handle standard attribute
							// store field attribute name
							name := constituent.Constituents[0].Value
							// store field attribute name as value
							value := name
							// store value if available
							if len(constituent.Constituents) > 1 {
								value = constituent.Constituents[1].Value
							}
							// store as label if named label
							if name == "label" {
								field.Label = value
								continue
							}
							// store id
							if name == "id" {
								field.ID = value
							}
							// store name
							if name == "name" {
								field.Name = value
							}
							// store attribute
							field.Attributes[name] = value
						} else {
							// handle array
							options := constituent.Constituents[0].Constituents
							// cycle through options
							for _, option := range options {
								// save option name and default value
								name := option.Constituents[0].Value
								value := strings.ToUpper(name[0:1]) + name[1:]
								// use option value if available
								if len(option.Constituents) > 1 {
									value = option.Constituents[1].Value
								}
								// add attribute
								field.Options[name] = value
							}
						}
					}
					// add field to model
					model.Fields = append(model.Fields, field)
					// return model
					return true
				},
			},
			"field type": {
				Constituents: [][]string{{"field name", "newline"}},
			},
			"field name": {
				Regex: `^[a-zA-Z][a-zA-Z0-9_-]+`, // grab up to newline
			},
			"field attribute": {
				Constituents: [][]string{{"hyphen", "name", "value?", "newline?"}, {"hyphen", "array", "newline?"}},
			},
			"array": {
				Constituents: [][]string{{"array open", "newline", "option*", "array close"}},
			},
			"array open": {
				Regex:  `^\[`,
				Ignore: true,
			},
			"option": {
				Constituents: [][]string{{"indent", "name", "value?", "newline"}},
			},
			"array close": {
				Regex:  `^\]`,
				Ignore: true,
			},
			"name": {
				Regex: `^([a-zA-Z0-9_\.-]+)([ ])?`, // grab up to and including first space
				FormatMatch: func(matches []string) string {
					return matches[1]
				},
			},
			"value": {
				Regex: `^[^\n]+`,
			},
			"hyphen": {
				Regex:  `^[-][ ]`,
				Ignore: true,
			},
			"indent": {
				Regex:  `^[ ][ ]`,
				Ignore: true,
			},
			"newline": {
				Regex:  `^[\n]`,
				Ignore: true,
			},
		},
	}
	return dialect
}

func (*DSL) NewModel() interface{} {
	model := &Model{}
	return model
}

func (*DSL) GenerateOutput(any interface{}) (string, error) {
	// ensure we have a Model
	model, ok := any.(*Model)
	if !ok {
		return "", errors.New("fastForms error: the appropriate model was not passed into the GenerateOutput function")
	}
	// create output pointer
	output := new(bytes.Buffer)
	// start form output
	output.WriteString("<form")
	// render form attributes
	renderFormAttributes(model.Attributes, output)
	// render close of form tag
	output.WriteString(">\n")
	// render fields
	renderFormFields(model.Fields, output)
	// render closing form tag
	output.WriteString("</form>\n")
	// return the final string
	return output.String(), nil
}

func capitalize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

func renderFieldLabel(field Field, output *bytes.Buffer) {
	if field.Label != "" {
		output.WriteString(Indent + Indent + "<label for=\"" + field.ID + "\">" + field.Label + "</label>\n")
	} else {
		if field.InputType != "submit" {
			output.WriteString(Indent + Indent + "<label for=\"" + field.ID + "\">" + capitalize(field.Name) + "</label>\n")
		}
	}
}

func renderFormAttributes(attributes []Attribute, output *bytes.Buffer) {
	for _, attribute := range attributes {
		output.WriteString(" " + attribute.Name + "=\"" + attribute.Value + "\"")
	}
}

func renderFormFields(fields []Field, output *bytes.Buffer) {
	for i, field := range fields {
		// output opening tag
		output.WriteString(Indent + "<div class=\"form-group\">\n")
		// ensure each field has id and name
		if field.Name == "" {
			field.Name = "field" + strconv.Itoa(i+1)
			field.Attributes["name"] = field.Name
		}
		if field.ID == "" {
			field.ID = field.Name
			field.Attributes["id"] = field.Name
		}
		// handle field-type-specific output
		switch field.InputType {
		case "textarea":
			renderFieldLabel(field, output)
			renderTextarea(field, output)
		case "select":
			renderFieldLabel(field, output)
			renderSelect(field, output)
		case "radio":
			renderRadio(field, output)
		case "checkbox":
			renderCheckbox(field, output)
		default:
			renderFieldLabel(field, output)
			renderInput(field, output)
		}
		// output closing tag
		output.WriteString(Indent + "</div>\n")
	}
}

func renderInput(field Field, output *bytes.Buffer) {
	output.WriteString(Indent + Indent + "<input type=\"" + field.InputType + "\"")

	for name, value := range field.Attributes {
		output.WriteString(" " + name + "=\"" + value + "\"")
	}

	output.WriteString(" />\n")
}
func renderTextarea(field Field, output *bytes.Buffer) {
	output.WriteString(Indent + Indent + "<textarea")

	for name, value := range field.Attributes {
		output.WriteString(" " + name + "=\"" + value + "\"")
	}

	output.WriteString("></textarea>\n")
}
func renderSelect(field Field, output *bytes.Buffer) {
	output.WriteString(Indent + Indent + "<select")

	for name, value := range field.Attributes {
		output.WriteString(" " + name + "=\"" + value + "\"")
	}

	output.WriteString(">\n")

	for name, value := range field.Options {
		output.WriteString(Indent + Indent + Indent + "<option value=\"" + name + "\">" + value + "</option>\n")
	}

	output.WriteString(Indent + Indent + "</select>\n")
}
func renderRadio(field Field, output *bytes.Buffer) {
	for name, value := range field.Options {
		output.WriteString(Indent + Indent + "<label><input type=\"radio\"")
		count := 0
		for attrName, attrValue := range field.Attributes {
			count = count + 1
			if attrName == "id" && count > 1 {
				continue
			}

			output.WriteString(" " + attrName + "=\"" + attrValue + "\"")
		}

		output.WriteString(" value=\"" + name + "\"/>" + value + "</label>\n")
	}
}
func renderCheckbox(field Field, output *bytes.Buffer) {
	for name, value := range field.Options {
		output.WriteString(Indent + Indent + "<label><input type=\"checkbox\"")
		count := 0
		for attrName, attrValue := range field.Attributes {
			count = count + 1
			if attrName == "id" && count > 1 {
				continue
			}

			output.WriteString(" " + attrName + "=\"" + attrValue + "\"")
		}
		output.WriteString(" value=\"" + name + "\"/>" + value + "</label>\n")
	}
}
