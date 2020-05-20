package main
import (
	"bytes"
	"flag"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"go/format"
	"io/ioutil"
	"strings"
)

var buffer  = bytes.NewBufferString("")

func main() {
	var specPath = flag.String("spec-path", "", "please use --spec-path /path/to/files/located")
	flag.Parse()
	specBytes, err := ioutil.ReadFile(*specPath)
	if err != nil {
		panic(err)
	}
	loader := openapi3.NewSwaggerLoader()
	doc, err := loader.LoadSwaggerFromData(specBytes)
	if err != nil {
		panic(err)
	}

	generatenewypes(doc)

	p, err := format.Source(buffer.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(p))

	//generateTeypes(doc)

}

func generatenewypes(doc *openapi3.Swagger) {
	for name, schema := range doc.Components.Schemas {
		generatetruct(name, schema)
	}
}

func generatetruct(name string, schema *openapi3.SchemaRef) {
	fmt.Fprintf(buffer,"\ntype %s struct {\n", name)
	if schema.Value == nil {
		return
	}

	for fieldName, field := range schema.Value.Properties{
		generatefield(fieldName, field)
	}

	fmt.Fprintf(buffer,"}\n")
}

func generatefield(name string, field *openapi3.SchemaRef) {
	gotype := resolvegotype(field.Value)
	name = strings.ToUpper(name[:1]) + name[1:]
	fmt.Fprintf(buffer,"%s %s %s\n", name, gotype, generatejsontag(name))
}

func resolvegotype(v *openapi3.Schema) string {
	switch v.Type {
	case "string":
		return "string"
	case "integer":
		return "int64"
	case "array":
		items := strings.Split(v.Items.Ref, "/")
		return fmt.Sprintf("[]%s", items[len(items) - 1])
	default:
		panic("unsupported type")
	}

}

func generatejsontag(name string) string {
	return fmt.Sprintf("`json:\"%s\"`", name)
}


func generateTeypes(doc *openapi3.Swagger) {


	for k, v := range doc.Components.Schemas {
		//if v.Ref != "" {
		//	fmt.Printf("ref: %v\n", v.Ref)
		//}
		fmt.Printf(" type %s struct {\n", k)
		if v.Value != nil {
			for k, v := range v.Value.Properties {
				fmt.Printf(" %s%s\t", strings.ToUpper(k[:1]), k[1:])
				if v.Value != nil {
				}
				if v.Value.Type == "integer" {
					fmt.Printf("int\t")
					fmt.Printf(" `json : "+`"`+"%v\n", v.Value.Description+`"`+"`")
				} else if v.Value.Type == "array" {
					fmt.Printf(" [] %s\t", k)
					fmt.Printf(" `json:"+`"`+"%v\n", v.Value.Description+`"`+"`")
				} else {
					fmt.Printf(" %+v\t", v.Value.Type)
					fmt.Printf(" `json:"+`"`+"%v\n", v.Value.Description+`"`+"`")
				}
			}

		}
		fmt.Println("}")
	}
}
