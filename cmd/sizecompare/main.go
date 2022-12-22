package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"google.golang.org/protobuf/proto"
	"moviehub.com/gen"
	"moviehub.com/metadata/pkg/model"
)

var metadata = &model.Metadata{
	ID: "123",
	Title: "The Movie 2",
	Director: "Foo Bars",
	Description: "Sequel of the legendary movie",
}

var genMetadata = &gen.Metadata{
	Id: "123",
	Title: "The Movie 2",
	Director: "Foo Bars",
	Description: "Sequel of the legendary movie",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))
	fmt.Printf("proto size:\t%dB\n", len(protoBytes))
}

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}