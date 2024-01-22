package templates

import (
	"fmt"
	"html/template"
	"os"
)

var Temp *template.Template

func InitTemplate() {
	temp, errTemp := template.ParseGlob("./templates/*home.gohtml")
	if errTemp != nil {
		fmt.Printf("Oupss une erreur lié au Templates : %v", errTemp.Error())
		os.Exit(1)
	}
	Temp = temp
}
