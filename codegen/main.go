package main
import (
	"flag"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"io/ioutil"
	"strings"
)
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
	for k, v := range doc.Components.Schemas {
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
					fmt.Printf(" `json : "+`"`+"%v\n", v.Value.Description+`"`+"`")
				} else {
					fmt.Printf(" %+v\t", v.Value.Type)
					fmt.Printf(" `json : "+`"`+"%v\n", v.Value.Description+`"`+"`")
				}
			}
			
		}
		fmt.Println("}")
	}

}
