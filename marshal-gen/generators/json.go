package generators

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// CustomArgs is used tby the go2idl framework to pass args specific to this
// generator.
type CustomArgs struct {
	ExtraPeerDirs []string // Always consider these as last-ditch possibilities for conversions.
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public":  namer.NewPublicNamer(0),
		"private": namer.NewPrivateNamer(0),
		"raw":     namer.NewRawNamer("", nil),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

// Packages makes the sets package definition.
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("// +build !%s\n\n", arguments.GeneratedBuildTag)))

	// We are generating defaults only for packages that are explicitly
	// passed as InputDir.
	for _, i := range context.Inputs {
		klog.V(5).Infof("considering pkg %q", i)

		pkg := context.Universe[i]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		typesPkg := pkg

		path := pkg.Path
		// if the source path is within a /vendor/ directory (for example,
		// k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/apis/meta/v1), allow
		// generation to output to the proper relative path (under vendor).
		// Otherwise, the generator will create the file in the wrong location
		// in the output directory.
		// TODO: build a more fundamental concept in gengo for dealing with modifications
		// to vendored packages.
		if strings.HasPrefix(pkg.SourcePath, arguments.OutputBase) {
			expandedPath := strings.TrimPrefix(pkg.SourcePath, arguments.OutputBase)
			if strings.Contains(expandedPath, "/vendor/") {
				path = expandedPath
			}
		}

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: filepath.Base(pkg.Path),
				PackagePath: path,
				HeaderText:  header,
				// GeneratorFunc returns a list of generators. Each generator makes a
				// single file.
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					generators = make([]generator.Generator, 0)
					// Since we want a file per type that we generate a set for, we
					// have to provide a function for this.
					for _, t := range c.Order {
						generators = append(generators, &marshalGen{
							DefaultGen: generator.DefaultGen{
								// Use the privatized version of the
								// type name as the file name.
								OptionalName: ToSnake(ToSnake(c.Namers["private"].Name(t) + "Marshal")),
							},
							targetPackage: pkg.Path,
							typeToMatch:   t,
							imports:       generator.NewImportTracker(),
						})
					}
					return generators
				},
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					if t.Name.Package != typesPkg.Path {
						return false
					}

					switch t.Kind {
					case types.Interface:
						return false
					default:
						return true
					}
				},
			})
	}

	return packages
}

type marshalGen struct {
	generator.DefaultGen
	targetPackage string
	typeToMatch   *types.Type
	imports       namer.ImportTracker
}

// Filter ignores all but one type because we're making a single file per type.
func (g *marshalGen) Filter(c *generator.Context, t *types.Type) bool { return t == g.typeToMatch }

func (g *marshalGen) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.targetPackage, g.imports),
	}
}

func (g *marshalGen) Imports(c *generator.Context) (imports []string) {
	importLines := []string{"encoding/json"}
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}

	return importLines
}

func (g *marshalGen) isOtherPackage(pkg string) bool {
	if pkg == g.targetPackage {
		return false
	}
	if strings.HasSuffix(pkg, "\""+g.targetPackage+"\"") {
		return false
	}
	return true
}

// args constructs arguments for templates. Usage:
// g.args(t, "key1", value1, "key2", value2, ...)
//
// 't' is loaded with the key 'type'.
//
// We could use t directly as the argument, but doing it this way makes it easy
// to mix in additional parameters. This feature is not used in this set
// generator, but is present as an example.
func (g *marshalGen) args(t *types.Type, kv ...interface{}) interface{} {
	m := map[interface{}]interface{}{"type": t}
	for i := 0; i < len(kv)/2; i++ {
		m[kv[i*2]] = kv[i*2+1]
	}
	return m
}

// GenerateType makes the body of a file implementing a set for type t.
func (g *marshalGen) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	sw.Do(templateCode, g.args(t))
	return sw.Error()
}

var templateCode = `
// MarshalJSON can marshal themselves into valid JSON.
func (obj *$.type|raw$) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj)
}

// UnmarshalJSON that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
func (obj *$.type|raw$) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	return nil
}

// String is used to print values passed as an operand
// to any format that accepts a string or to an unformatted printer
// such as Print.
func (obj *$.type|raw$) String() string {
	bs, _ := obj.MarshalJSON()
	return string(bs)
}
`

// ToSnake converts a string to snake_case
func ToSnake(s string) string {
	return ToDelimited(s, '_')
}

// ToDelimited converts a string to delimited.snake.case (in this case `del = '.'`)
func ToDelimited(s string, del uint8) string {
	return ToScreamingDelimited(s, del, false)
}

// ToScreamingDelimited converts a string to SCREAMING.DELIMITED.SNAKE.CASE (in this case `del = '.'; screaming = true`) or delimited.snake.case (in this case `del = '.'; screaming = false`)
func ToScreamingDelimited(s string, del uint8, screaming bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			if (v >= 'A' && v <= 'Z' && next >= 'a' && next <= 'z') || (v >= 'a' && v <= 'z' && next >= 'A' && next <= 'Z') {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && n[len(n)-1] != del && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += string(del) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(del)
			}
		} else if v == ' ' || v == '_' || v == '-' {
			// replace spaces/underscores with delimiters
			n += string(del)
		} else {
			n = n + string(v)
		}
	}

	if screaming {
		n = strings.ToUpper(n)
	} else {
		n = strings.ToLower(n)
	}
	return n
}

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}
