package fhirongo

import (
	"encoding/json"
	"encoding/base64"
	"io/ioutil"
	"fmt"
	"os"
)

type Image struct {
	ResourceType	string 		`json:"resourceType" bson:"resource_type"`
	Id 				string 		`json:"id"`
	Meta			MetaData	`json:"meta"`
	ContentType 	string 		`json:"contentType" bson:"content_type`
	Content 		string 		`json:"content"`
}

func (c *Connection) GetImage(id string) (*Image, error) {
	baseUrl := c.BaseURL
	image := &Image{}
	url := fmt.Sprintf("%sBinary/XR-%s",baseUrl,id)
	fmt.Printf("url : url")
	bytes, err := c.GetFhir(url)
	if err != nil {
		return nil, fmt.Errorf("error getting a pdfImage: %s", err.Error())
	}
	
	err = json.Unmarshal(bytes, image)
	// if err := os.WriteFile("./debbie.pdf", bytes, 0666); err != nil {
	// 	log.Fatal(err)
	// }
	//fmt.Printf("Wrote the pdf file\n")
	return image, err
}

func (c *Connection) GetPDF(id string) (string, error) {
	image, err := c.GetImage(id)
	if err != nil {
		return "", err
	}
	pdf, err := base64.StdEncoding.DecodeString(image.Content)
	if err != nil {
		return "", err
	}
	return string(pdf), nil
}

func (c *Connection) GetPDFb64(id string) (string, error) {
	baseUrl := c.BaseURL
	image := &Image{}
	url := fmt.Sprintf("%sBinary/XR-%s", baseUrl, id)
	fmt.Printf("url : url")
	bytes, err := c.GetFhir(url)
	if err != nil {
		return "", fmt.Errorf("error getting a pdfImage: %s", err.Error())
	}
	
	err = json.Unmarshal(bytes, image)
	// if err := os.WriteFile("./debbie.pdf", bytes, 0666); err != nil {
	// 	log.Fatal(err)
	// }
	//fmt.Printf("Wrote the pdf file\n")
	return image.Content, err
}

func Decode(fileName string, b64 string )  {
	//b64, _ := ioutil.ReadFile("../../sample.pdf")
	//b64, _ := ioutil.ReadFile("../model/sample.bas64")
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

func WriteImage(fileName string, pdf string )  {
	//b64, _ := ioutil.ReadFile("../../sample.pdf")
	//b64, _ := ioutil.ReadFile("../model/sample.bas64")
	// dec, err := base64.StdEncoding.DecodeString(b64)
	// if err != nil {
	// 	panic(err)
	// }

	// f, err := os.Create(fileName)
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	err := ioutil.WriteFile(fileName, []byte(pdf), 0666)
	if err != nil {
		fmt.Printf("Error writing pdf file: %s\n", err.Error())
	}


}