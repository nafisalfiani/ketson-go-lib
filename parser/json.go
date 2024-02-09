package parser

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/nafisalfiani/ketson-go-lib/codes"
	"github.com/nafisalfiani/ketson-go-lib/errors"
	"github.com/nafisalfiani/ketson-go-lib/log"
	"github.com/xeipuuv/gojsonschema"
)

type jsonConfig string

const (
	vanillaCompatible jsonConfig = "standard"
	defaultConfig     jsonConfig = "default"
	fastestConfig     jsonConfig = "fastest"
	customConfig      jsonConfig = "custom"
	errJSON                      = `JSON parser error %s`
)

type JSONOptions struct {
	Config                        jsonConfig
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJSONRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
	// Schema contains schema definitions with key as schema name and value as source path
	// schema sources can be file or URL. Schema definition will be iniatialized during
	// JSON parser object initialization.
	Schema map[string]string
}

type JSONInterface interface {
	// Marshal go structs into bytes
	Marshal(orig interface{}) ([]byte, error)
	// Unmarshal bytes into go structs
	Unmarshal(blob []byte, dest interface{}) error
}

type jsonParser struct {
	API    jsoniter.API
	schema map[string]*gojsonschema.Schema
	logger log.Interface
}

func initJSON(opt JSONOptions, log log.Interface) JSONInterface {
	var jsonAPI jsoniter.API
	switch opt.Config {
	case defaultConfig:
		jsonAPI = jsoniter.ConfigDefault
	case fastestConfig:
		jsonAPI = jsoniter.ConfigFastest
	case customConfig:
		jsonAPI = jsoniter.Config{
			IndentionStep:                 opt.IndentionStep,
			MarshalFloatWith6Digits:       opt.MarshalFloatWith6Digits,
			EscapeHTML:                    opt.EscapeHTML,
			SortMapKeys:                   opt.SortMapKeys,
			UseNumber:                     opt.UseNumber,
			DisallowUnknownFields:         opt.DisallowUnknownFields,
			TagKey:                        opt.TagKey,
			OnlyTaggedField:               opt.OnlyTaggedField,
			ValidateJsonRawMessage:        opt.ValidateJSONRawMessage,
			ObjectFieldMustBeSimpleString: opt.ObjectFieldMustBeSimpleString,
			CaseSensitive:                 opt.CaseSensitive,
		}.Froze()
	default:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}

	p := &jsonParser{
		API:    jsonAPI,
		schema: make(map[string]*gojsonschema.Schema),
		logger: log,
	}

	// init defined schema
	p.initSchema(opt.Schema)

	return p
}

// initSchema initialize all defined schema from sources
func (p *jsonParser) initSchema(sources map[string]string) {
	for sch, src := range sources {
		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(src))
		if err != nil {
			p.logger.Fatal(context.Background(), errors.NewWithCode(codes.CodeJSONSchemaInvalid, fmt.Sprintf("error on load : %s", sch), err))
			return
		}
		p.schema[sch] = schema
	}
}

// Marshal marshal input blobs into go structs.
func (p *jsonParser) Marshal(orig interface{}) ([]byte, error) {
	stream := p.API.BorrowStream(nil)
	defer p.API.ReturnStream(stream)
	stream.WriteVal(orig)
	result := make([]byte, stream.Buffered())
	if stream.Error != nil {
		return nil, errors.NewWithCode(codes.CodeJSONMarshalError, stream.Error.Error())
	}
	copy(result, stream.Buffer())
	return result, nil
}

// Unmarshal unmarshal input blobs into go structs.
func (p *jsonParser) Unmarshal(blob []byte, dest interface{}) error {
	iter := p.API.BorrowIterator(blob)
	defer p.API.ReturnIterator(iter)
	iter.ReadVal(dest)
	if iter.Error != nil {
		return errors.NewWithCode(codes.CodeJSONUnmarshalError, iter.Error.Error())
	}
	return nil
}
