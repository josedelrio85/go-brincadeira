package utilities

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// Fileconfig is a struct that represents the structure of the yaml file from the config is obtained.
type Fileconfig struct {
	Source struct {
		Tipo      string `yaml:"type"`
		Path      string `yaml:"path"`
		Extension string `yaml:"extension"`
	} `yaml:"source"`
	Proccess struct {
		Tipo        string `yaml:"type"`
		Name        string `yaml:"name"`
		Packagename string `yaml:"packagename"`
		Estructura  []struct {
			Field      string       `yaml:"field"`
			TipoCampo  string       `yaml:"type"`
			Structname string       `yaml:"structname"`
			Structure  []properties `yaml:"structure"`
		} `yaml:"structure"`
	} `yaml:"proccess"`
	Destination struct {
		Tipo       string `yaml:"type"`
		Tabla      string `yaml:"table"`
		Estructura []struct {
			Field     string `yaml:"field"`
			Order     int    `yaml:"order"`
			TipoCampo string `yaml:"type"`
		} `yaml:"structure"`
	} `yaml:"destination"`
}

// strstructconfig is a struct that represents the structure of a Struct Type.
//		- Packagename is the name of the package where the struct will be stored.
//		- Namestruct is the name of the Struct.
//		- Property is an array of properties struct.
type structconfig struct {
	Packagename string
	Namestruct  string
	Property    []properties
}

// properties is a struct that represents a property of a Struct Type.
// Contains a Name (string) and a Type(string) to set the name and the type of the struct property.
type properties struct {
	Structname string
	Type       string
	Field      string
	Comillas   string
	Structure  []properties
}

//Parse function populates a Fileconfig struct with values from a Yaml configuration file path provided as param.
func (f Fileconfig) Parse(path string) (Fileconfig, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return f, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error: %v", err)
		return f, err
	}

	err = yaml.Unmarshal([]byte(data), &f)
	if err != nil {
		log.Fatalf("error marshaling: %v", err)
		return f, err
	}

	fmt.Println("--------- Parse result -----------")
	fmt.Println(f)
	fmt.Println("--------- ------------ -----------")
	return f, nil
}

func checkStructure(p *properties, array []properties) {
	if len(array) > 0 {
		for _, row := range array {
			subp := properties{
				Structname: row.Structname,
				Type:       row.Type,
				Field:      row.Field,
				Comillas:   "a",
			}

			if len(row.Structure) > 0 {
				b := subcheckStructure(row.Structure)
				subp.Structure = b
			}
			p.Structure = append(p.Structure, subp)
		}
	}
}

func subcheckStructure(array []properties) []properties {
	list := make([]properties, 0)
	for _, r := range array {
		subp := properties{
			Structname: r.Structname,
			Type:       r.Type,
			Field:      r.Field,
			Comillas:   "`",
		}
		if len(r.Structure) > 0 {
			a := subcheckStructure(r.Structure)
			subp.Structure = a
		}
		list = append(list, subp)
	}
	return list
}

// GenerateStruct is a public function that generates a Struct type dinamically.
// Uses a Fileconfig struct populated with parse method and creates a .go file with the properties and names provided in the config file.
func (f Fileconfig) GenerateStruct(pathconfig string) {

	fc, err := f.Parse(pathconfig)
	if err != nil {
		log.Println(err)
	}

	sc := structconfig{
		Packagename: fc.Proccess.Packagename,
		Namestruct:  fc.Proccess.Name,
		Property:    make([]properties, 0),
	}

	for _, jsonRow := range fc.Proccess.Estructura {
		p := properties{
			Structname: jsonRow.Structname,
			Type:       jsonRow.TipoCampo,
			Field:      jsonRow.Field,
			Comillas:   "`",
		}
		checkStructure(&p, jsonRow.Structure)
		sc.Property = append(sc.Property, p)
	}
	fmt.Println("--------- Structconfig result -----------")
	fmt.Println(sc)
	fmt.Println("--------- ------------ -----------")

	funcMap := template.FuncMap{
		"testFunc": func(array structconfig) string {
			a := ""
			a += "package " + array.Packagename
			a += "\n"
			a += "type " + array.Namestruct + " struct { "
			a += "\n"

			for _, k := range array.Property {

				if len(k.Structure) > 0 {
					for _, row := range k.Structure {
						a += row.Structname + " " + row.Type + "{"

						if len(row.Structure) > 0 {
							for _, r := range row.Structure {
								a += r.Structname + " " + r.Type + " " + r.Comillas + "json:\"" + r.Field + "\"" + r.Comillas + "\n"
							}
						}
						a += "}`json:\"" + row.Field + "\"`\n"
					}
				} else {
					a += k.Structname + " " + k.Type + " " + k.Comillas + "json:\"" + k.Field + "\"" + k.Comillas
					a += "\n"
				}
			}
			a += " }"
			return a
		},
	}

	var testTemplate = `
		{{ testFunc . }}
	`
	templ, ferr := template.New("test").Funcs(funcMap).Parse(testTemplate)
	err = templ.Execute(os.Stdout, sc)

	if _, err := os.Stat("./generatedStructs/"); !os.IsNotExist(err) {
		// path does not exist
		os.Mkdir("./generatedStructs/", 0755)
	}

	file, ferr := os.Create("./generatedStructs/" + fc.Proccess.Name + ".go")
	if ferr != nil {
		log.Println(ferr)
		return
	}
	defer file.Close()

	perr := templ.Execute(file, sc)
	if perr != nil {
		log.Println(perr)
	}
}

// GenerateStruct is a public function that generates a Struct type dinamically.
// Uses a Fileconfig struct populated with parse method and creates a .go file with the properties and names provided in the config file.
// In this case text templates are used to print the final values.
// Not use
func (f Fileconfig) withTemplates(pathconfig string) {

	fc, err := f.Parse(pathconfig)
	if err != nil {
		log.Println(err)
	}

	sc := structconfig{
		Packagename: fc.Proccess.Packagename,
		Namestruct:  fc.Proccess.Name,
		Property:    make([]properties, 0),
	}

	for _, jsonRow := range fc.Proccess.Estructura {
		p := properties{
			Structname: jsonRow.Structname,
			Type:       jsonRow.TipoCampo,
			Field:      jsonRow.Field,
			Comillas:   "`",
		}
		checkStructure(&p, jsonRow.Structure)
		sc.Property = append(sc.Property, p)
	}
	fmt.Println("--------- Structconfig result -----------")
	fmt.Println(sc)
	fmt.Println("--------- ------------ -----------")

	var structTemplate = `
		package {{.Packagename}}

		type {{.Namestruct}} struct {
	`
	var keyvalueTemplate = `
		{{range $y, $x := .Property }}
			{{if .Structure}}
				{{ $x.Structname }}  {{ $x.Type }} {
				{{range $yy, $xx := .Structure }}
					{{ $xx.Structname }}  {{ $xx.Type }} {{ $xx.Comillas }}json:"{{ $xx.Field }}"{{ $xx.Comillas }}
				{{end}}
				}  {{ $x.Comillas }}json:"{{ $x.Field }}"{{ $x.Comillas }}
			{{else}}
				{{ $x.Structname }}  {{ $x.Type }}  {{ $x.Comillas }}json:"{{ $x.Field }}"{{ $x.Comillas }}
			{{end}}
		{{end}}
	`

	structTemp := template.Must(template.New("structTemplate").Parse(structTemplate))
	keyvalueTemp := template.Must(template.New("keyvalueTemplate").Parse(keyvalueTemplate))

	parent := template.Must(template.New("finalTemplate").Parse(`{{ template "structTemplate" .}}{{ template "keyvalueTemplate" .}} }`))

	addChildTemplate(parent, structTemp)
	addChildTemplate(parent, keyvalueTemp)

	perr := parent.Execute(os.Stdout, sc)
	if perr != nil {
		log.Println(perr)
	}
}

// addChildTemplate adds child templates to a parent template
func addChildTemplate(parent *template.Template, child *template.Template) (*template.Template, error) {
	return parent.AddParseTree(child.Name(), child.Tree)
}
