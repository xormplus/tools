package main

import (
	"log"
	"strings"

	"github.com/ChimeraCoder/gojson"
	//	"github.com/oskca/sciter"
)

const json = "json2struct"

func DoJson2Struct() {

	root, err := w.GetRootElement()
	if err != nil {
		PrintResult(json, err.Error())
		return
	}
	resultElement, err := root.SelectById("source")
	if err != nil {
		PrintResult(json, err.Error())
		return
	}
	jsonStr, err := resultElement.GetValue()
	if err != nil {
		PrintResult(json, err.Error())
		return
	}
	resultElement, err = root.SelectById("package")

	if err != nil {
		PrintResult(json, err.Error())
		return
	}
	packageStr, err := resultElement.GetValue()
	log.Println(packageStr)
	if err != nil {
		PrintResult(json, err.Error())
		return
	}
	resultElement, err = root.SelectById("structname")
	if err != nil {
		PrintResult(json, err.Error())
		return
	}
	structName, err := resultElement.GetValue()
	log.Println(structName)
	if err != nil {
		PrintResult(json, err.Error())
		return
	}
	j := strings.NewReader(jsonStr.String())
	b, err := json2struct.Generate(j, structName.String(), packageStr.String())

	if err != nil {
		PrintResult(json, err.Error())
		return
	}

	PrintResult(json, string(b))
}

func PrintResult(id, msg string) error {
	root, err := w.GetRootElement()
	if err != nil {
		return err
	}

	resultElement, err := root.SelectById(id)
	if err != nil {
		return err
	}

	err = resultElement.SetText(msg)
	if err != nil {
		return err
	}
	root.Update(true)
	return nil
}
