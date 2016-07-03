package main

import (
	"github.com/wicast/xj2s"
)

const xml = "xml2struct"

func DoXml2Struct() {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			PrintResult(xml, err.Error())
		}
	}()
	root, err := w.GetRootElement()
	if err != nil {
		PrintResult(xml, err.Error())
		return
	}
	resultElement, err := root.SelectById("xmlsource")
	if err != nil {
		PrintResult(xml, err.Error())
		return
	}
	xmlStr, err := resultElement.GetValue()
	if err != nil {
		PrintResult(xml, err.Error())
		return
	}

	str := xj2s.Xml2Struct([]byte(xmlStr.String()), false)

	if err != nil {
		PrintResult(xml, err.Error())
		return
	}

	PrintResult(xml, str)
}
