package utilities

import (
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
		Structname  string `yaml:"structname"`
		Packagename string `yaml:"packagename"`
	} `yaml:"proccess"`
	Destination struct {
		Tipo           string   `yaml:"type"`
		Tabla          string   `yaml:"table"`
		Estructura     []string `yaml:"structure"`
		EstructuraJSON []struct {
			Field     string `yaml:"field"`
			Order     int    `yaml:"order"`
			TipoCampo string `yaml:"type"`
		} `yaml:"structurejson"`
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
	Name string
	Type string
}

//parse function populates a Fileconfig struct with values from a Yaml configuration file provided as param.
func (f Fileconfig) parse(path string) (Fileconfig, error) {

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return f, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	err = yaml.Unmarshal([]byte(data), &f)
	if err != nil {
		log.Fatalf("error: %v", err)
		return f, err
	}
	return f, nil
}

// GenerateStruct is a public function that generates a Struct type dinamically.
// Uses a Fileconfig struct populated with parse method and creates a .go file with the properties and names provided in the config file.
func (f Fileconfig) GenerateStruct(pathconfig string) {

	fc, err := f.parse(pathconfig)
	if err != nil {
		log.Println(err)
	}

	sc := structconfig{
		Packagename: fc.Proccess.Packagename,
		Namestruct:  fc.Proccess.Structname,
		Property:    make([]properties, 0),
	}

	for _, jsonRow := range fc.Destination.EstructuraJSON {
		sc.Property = append(sc.Property, properties{Name: jsonRow.Field, Type: jsonRow.TipoCampo})
	}

	var structTemplate = `
		package {{.Packagename}}

		type {{.Namestruct}} struct {
	`
	var keyvalueTemplate = `
		{{range $y, $x := .Property }}
		{{ $x.Name }}  {{ $x.Type }}
		{{end}}
	`

	structTemp := template.Must(template.New("structTemplate").Parse(structTemplate))
	keyvalueTemp := template.Must(template.New("keyvalueTemplate").Parse(keyvalueTemplate))

	parent := template.Must(template.New("finalTemplate").Parse(`{{ template "structTemplate" .}}{{ template "keyvalueTemplate" .}} }`))

	addChildTemplate(parent, structTemp)
	addChildTemplate(parent, keyvalueTemp)

	// file, ferr := os.Create("./generatedStructs/" + f.Proccess.Structname + ".go")
	file, ferr := os.Create(fc.Proccess.Structname + ".go")
	if ferr != nil {
		log.Println(ferr)
		return
	}
	defer file.Close()

	perr := parent.Execute(file, sc)
	if perr != nil {
		log.Println(perr)
	}
}

// addChildTemplate adds child templates to a parent template
func addChildTemplate(parent *template.Template, child *template.Template) (*template.Template, error) {
	return parent.AddParseTree(child.Name(), child.Tree)
}
