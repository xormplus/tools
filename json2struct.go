package main

import (
	"log"
	"strings"

	"github.com/ChimeraCoder/gojson"
	//	"github.com/oskca/sciter"
)

func DoJson2Struct() {

	root, err := w.GetRootElement()
	if err != nil {
		PrintResult(err.Error())
		return
	}
	resultElement, err := root.SelectById("source")
	if err != nil {
		PrintResult(err.Error())
		return
	}
	jsonStr, err := resultElement.GetValue()
	if err != nil {
		PrintResult(err.Error())
		return
	}
	resultElement, err = root.SelectById("package")

	if err != nil {
		PrintResult(err.Error())
		return
	}
	packageStr, err := resultElement.GetValue()
	log.Println(packageStr)
	if err != nil {
		PrintResult(err.Error())
		return
	}
	resultElement, err = root.SelectById("structname")
	if err != nil {
		PrintResult(err.Error())
		return
	}
	structName, err := resultElement.GetValue()
	log.Println(structName)
	if err != nil {
		PrintResult(err.Error())
		return
	}
	json := strings.NewReader(jsonStr.String())
	b, err := json2struct.Generate(json, structName.String(), packageStr.String())

	if err != nil {
		PrintResult(err.Error())
		return
	}

	PrintResult(string(b))
}

func PrintResult(msg string) error {
	root, err := w.GetRootElement()
	if err != nil {
		return err
	}

	resultElement, err := root.SelectById("xml2struct")
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
