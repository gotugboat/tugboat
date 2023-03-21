package tmpl

import (
	"bytes"
	"strings"
	"text/template"
	"time"
	"tugboat/internal/pkg/flags"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Template holds data that can be applied to a template string
type Template struct {
	fields Fields
}

// Fields that will be available to the template engine
type Fields map[string]interface{}

const (
	imageName   = "ImageName"
	version     = "Version"
	tag         = "Tag"
	branch      = "Branch"
	commit      = "Commit"
	shortCommit = "ShortCommit"
	fullCommit  = "FullCommit"
)

func New(opts *flags.Options) *Template {
	return &Template{
		fields: Fields{
			imageName:   opts.Image.Name,
			version:     opts.Image.Version,
			tag:         opts.Global.Git.Tag,
			branch:      opts.Global.Git.Branch,
			commit:      opts.Global.Git.Commit,
			shortCommit: opts.Global.Git.ShortCommit,
			fullCommit:  opts.Global.Git.FullCommit,
		},
	}
}

// Apply applies the given string against the Fields stored in the template
func (t *Template) Apply(s string) (string, error) {
	var output bytes.Buffer
	tmpl, err := template.New("tmpl").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"replace": strings.ReplaceAll,
			"split":   strings.Split,
			"time": func(s string) string {
				return time.Now().UTC().Format(s)
			},
			"tolower":    strings.ToLower,
			"toupper":    strings.ToUpper,
			"trim":       strings.TrimSpace,
			"trimprefix": strings.TrimPrefix,
			"trimsuffix": strings.TrimSuffix,
			"title":      cases.Title(language.English).String,
		}).
		Parse(s)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&output, t.fields)
	return output.String(), err
}

// CompileStringSlice will apply a template against a slice of strings
func CompileStringSlice(input []string, opts *flags.Options) ([]string, error) {
	compiledItems := []string{}
	for _, item := range input {
		tmpl, err := New(opts).Apply(item)
		if err != nil {
			return nil, err
		}
		compiledItems = append(compiledItems, tmpl)
	}
	return compiledItems, nil
}

// CompileString will apply a template against a given string
func CompileString(input string, opts *flags.Options) (string, error) {
	tmpl, err := New(opts).Apply(input)
	if err != nil {
		return "", err
	}
	return tmpl, nil
}
