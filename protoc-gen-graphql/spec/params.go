package spec

import (
	"errors"
	"regexp"
	"strings"
)

// Params spec have plugin parameters
type Params struct {
	QueryOut       string
	Excludes       []*regexp.Regexp
	Verbose        bool
	FieldCamelCase bool
	RelativePaths  bool
}

func NewParams(p string) (*Params, error) {
	params := &Params{
		Excludes: []*regexp.Regexp{},
	}
	if p == "" {
		return params, nil
	}

	for _, v := range strings.Split(p, ",") {
		kv := strings.SplitN(v, "=", 2)
		switch kv[0] {
		case "verbose":
			params.Verbose = true
		case "query":
			if len(kv) == 1 {
				return nil, errors.New("argument " + kv[0] + " must have value")
			}
			params.QueryOut = kv[1]
		case "exclude":
			if len(kv) == 1 {
				return nil, errors.New("argument " + kv[0] + " must have value")
			}
			regex, err := regexp.Compile(kv[1])
			if err != nil {
				return nil, errors.New("failed to compile regex for exclude argument " + kv[1])
			}
			params.Excludes = append(params.Excludes, regex)
		case "field_camel":
			params.FieldCamelCase = true
		case "paths":
			if len(kv) == 1 {
				return nil, errors.New("argument " + kv[0] + " must have value")
			}
			params.RelativePaths = kv[1] == "source_relative"
		default:
			return nil, errors.New("Unacceptable argument " + kv[0] + " provided")
		}
	}
	return params, nil
}

func (p *Params) IsExclude(pkg string) bool {
	for _, r := range p.Excludes {
		if r.MatchString(pkg) {
			return true
		}
	}
	return false
}
