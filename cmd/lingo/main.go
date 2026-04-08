// lingo is a code generator for go-i18n template data structs.
//
// It reads a package-level variable of type Underliers from the target
// repository's locale package and generates three files:
//
//   - messages-cobra-auto.go
//   - messages-errors-auto.go
//   - messages-general-auto.go
//
// # Installation
//
//	go install github.com/snivilised/li18ngo/cmd/lingo@latest
//
// # Usage
//
// Add the following directive to underlying-templ-data.go in your locale
// package:
//
//	//go:generate lingo
//
// Then run:
//
//	go generate ./...
//
// # Flags
//
//	--locale <path>   Path to the locale directory, relative to the repo
//	                  root (detected via go.mod). If omitted, the generator
//	                  tries ./locale then ./src/locale, then searches the
//	                  whole repo.
//	--dry-run         Validate the Underliers map and report all errors
//	                  without writing any files.
//
//nolint:all
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/snivilised/li18ngo/locale/enums"
)

// ---------------------------------------------------------------------------
// UnderlyingType reverse lookup
// ---------------------------------------------------------------------------

// underlyingTypeByName maps the stringer-trimmed constant names
// (e.g. "StaticCobra") back to their enum values.  The trimmed names come
// from identOrSel, which returns the selector or ident text as written in
// the source — e.g. enums.UnderlyingTypeStaticCobra → "UnderlyingTypeStaticCobra".
// We therefore need to strip the "UnderlyingType" prefix ourselves here.
var underlyingTypeByName = func() map[string]enums.UnderlyingType {
	all := []enums.UnderlyingType{
		enums.UnderlyingTypeUndefined,
		enums.UnderlyingTypeStaticCobra,
		enums.UnderlyingTypeDynamicCobra,
		enums.UnderlyingTypeStaticGeneral,
		enums.UnderlyingTypeDynamicGeneral,
		enums.UnderlyingTypeStaticError,
		enums.UnderlyingTypeSentinelError,
		enums.UnderlyingTypeStaticErrorWrapper,
		enums.UnderlyingTypeDynamicError,
		enums.UnderlyingTypeDynamicErrorWrapper,
	}
	m := make(map[string]enums.UnderlyingType, len(all))
	for _, v := range all {
		// String() returns the trimprefix'd name, e.g. "StaticCobra".
		m[v.String()] = v
		// Also register the full constant name so that source written as
		// UnderlyingTypeStaticCobra (without a package selector) resolves.
		m["UnderlyingType"+v.String()] = v
	}
	return m
}()

// parseUnderlyingType converts the string produced by identOrSel back to an
// enum value.  identOrSel returns the selector name from expressions like
// enums.UnderlyingTypeStaticCobra → "UnderlyingTypeStaticCobra", or just
// the ident for an unqualified reference → "UnderlyingTypeStaticCobra".
func parseUnderlyingType(s string) (enums.UnderlyingType, error) {
	if v, ok := underlyingTypeByName[s]; ok {
		return v, nil
	}
	return enums.UnderlyingTypeUndefined,
		fmt.Errorf("unknown UnderlyingType constant %q", s)
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "lingo: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	localeFlagVal := flag.String("locale", "", "path to locale dir relative to repo root")
	dryRun := flag.Bool("dry-run", false, "validate only, do not write files")
	verbose := flag.Bool("verbose", false, "print per-message diagnostic info during validation")
	flag.Parse()

	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	localeDir, err := resolveLocaleDir(repoRoot, *localeFlagVal)
	if err != nil {
		return err
	}

	underliers, baseStruct, pkgName, err := parseUnderliers(localeDir)
	if err != nil {
		return err
	}

	if err := validate(underliers, *verbose); err != nil {
		return err
	}

	if *dryRun {
		fmt.Printf("lingo: dry-run OK — no errors found, no files written (locale dir: '%s')\n", localeDir)
		return nil
	}

	return generate(localeDir, pkgName, baseStruct, underliers)
}

// ---------------------------------------------------------------------------
// Repo root detection
// ---------------------------------------------------------------------------

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not find go.mod — are you inside a Go module?")
		}
		dir = parent
	}
}

// ---------------------------------------------------------------------------
// Locale directory resolution
// ---------------------------------------------------------------------------

func resolveLocaleDir(repoRoot, flagVal string) (string, error) {
	if flagVal != "" {
		p := filepath.Join(repoRoot, flagVal)
		if err := checkLocaleDir(p); err != nil {
			return "", fmt.Errorf("--locale %q: %w", flagVal, err)
		}
		return p, nil
	}

	// Try the two conventional defaults.
	for _, rel := range []string{"locale", filepath.Join("src", "locale")} {
		p := filepath.Join(repoRoot, rel)
		if err := checkLocaleDir(p); err == nil {
			return p, nil
		}
	}

	// Fall back to full repo search.
	found, err := searchForLocaleDir(repoRoot)
	if err != nil {
		return "", err
	}
	return found, nil
}

func checkLocaleDir(p string) error {
	info, err := os.Stat(p)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", p)
	}
	// Must contain at least one .go file with an Underliers declaration.
	has, err := dirContainsUnderliers(p)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("%q contains no Underliers declaration", p)
	}
	return nil
}

func dirContainsUnderliers(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		if bytes.Contains(src, []byte("Underliers")) {
			return true, nil
		}
	}
	return false, nil
}

func searchForLocaleDir(repoRoot string) (string, error) {
	var found []string
	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		has, _ := dirContainsUnderliers(path)
		if has {
			found = append(found, path)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("searching repo for locale dir: %w", err)
	}
	switch len(found) {
	case 0:
		return "", errors.New("could not find a directory containing an Underliers declaration; use --locale to specify")
	case 1:
		return found[0], nil
	default:
		return "", fmt.Errorf("found multiple directories with Underliers declarations: %v; use --locale to specify", found)
	}
}

// ---------------------------------------------------------------------------
// AST parsing
// ---------------------------------------------------------------------------

// underlierEntry is the in-memory representation of one Underliers map entry.
type underlierEntry struct {
	MessageID   string
	Seed        string
	TypeName    enums.UnderlyingType
	Description string
	Story       string
	Other       string
	Fields      []fieldEntry
}

type fieldEntry struct {
	Note   string
	GoType string
	Tale   string
}

func parseUnderliers(dir string) ([]underlierEntry, string, string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)

	if err != nil {
		return nil, "", "", fmt.Errorf("parsing locale dir: %w", err)
	}

	var (
		entries    []underlierEntry
		baseStruct string
		pkgName    string
	)

	for name, pkg := range pkgs {
		pkgName = name
		for _, file := range pkg.Files {
			if bs := findBaseStruct(file); bs != "" && baseStruct == "" {
				baseStruct = bs
			}
			found, err := extractUnderliers(file)
			if err != nil {
				return nil, "", "", err
			}
			entries = append(entries, found...)
		}
	}

	if baseStruct == "" {
		return nil, "", "", errors.New("could not find base embed struct (expected a struct with a SourceID() string method)")
	}
	if len(entries) == 0 {
		return nil, "", "", errors.New("no Underliers entries found in locale dir")
	}

	return entries, baseStruct, pkgName, nil
}

// findBaseStruct locates the struct that has a SourceID() string method,
// which is the base embed struct for all TemplData types.
func findBaseStruct(file *ast.File) string {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "SourceID" || fn.Recv == nil {
			continue
		}
		if len(fn.Recv.List) != 1 {
			continue
		}
		// Check return type is string.
		if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
			continue
		}
		if id, ok := fn.Type.Results.List[0].Type.(*ast.Ident); ok && id.Name == "string" {
			// Receiver type name is the base struct.
			switch t := fn.Recv.List[0].Type.(type) {
			case *ast.Ident:
				return t.Name
			case *ast.StarExpr:
				if id, ok := t.X.(*ast.Ident); ok {
					return id.Name
				}
			}
		}
	}
	return ""
}

// isUnderliersType reports whether an AST expression refers to the
// Underliers type, handling both forms a client package may use:
//
//   - bare ident:      Underliers{...}         (local alias / dot-import)
//   - selector expr:   li18ngo.Underliers{...}  (normal qualified import)
//
// Only the type-name "Underliers" is checked; the package qualifier is
// accepted as-is so that clients are free to alias the import.
func isUnderliersType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name == "Underliers"
	case *ast.SelectorExpr:
		return t.Sel.Name == "Underliers"
	}
	return false
}

// extractUnderliers finds the var declaration of type Underliers and
// extracts all map entries from its composite literal.
//
// It recognises three declaration styles:
//
//  1. Explicit bare type:
//     var messages Underliers = Underliers{...}
//
//  2. Explicit qualified type:
//     var messages li18ngo.Underliers = li18ngo.Underliers{...}
//
//  3. Inferred type (most common):
//     var messages = li18ngo.Underliers{...}
func extractUnderliers(file *ast.File) ([]underlierEntry, error) {
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			// Explicit type annotation: var messages Underliers = ...
			// or var messages li18ngo.Underliers = ...
			if vs.Type != nil && !isUnderliersType(vs.Type) {
				continue
			}
			// Inferred type — inspect the composite literal directly.
			for _, val := range vs.Values {
				cl, ok := val.(*ast.CompositeLit)
				if !ok {
					continue
				}
				if !isUnderliersType(cl.Type) {
					continue
				}
				return extractMapEntries(cl)
			}
		}
	}
	return nil, nil
}

func extractMapEntries(cl *ast.CompositeLit) ([]underlierEntry, error) {
	var entries []underlierEntry
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		val, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			continue
		}
		entry, err := extractEntry(val)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func extractEntry(cl *ast.CompositeLit) (underlierEntry, error) {
	var e underlierEntry
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "MessageID":
			e.MessageID = stringLit(kv.Value)
		case "Seed":
			e.Seed = stringLit(kv.Value)
		case "TypeName":
			// identOrSel returns e.g. "UnderlyingTypeStaticCobra" (from a
			// selector expression enums.UnderlyingTypeStaticCobra) or the bare
			// ident name for an unqualified reference.
			raw := identOrSel(kv.Value)
			ut, err := parseUnderlyingType(raw)
			if err != nil {
				return e, fmt.Errorf("entry %q TypeName: %w", e.MessageID, err)
			}
			e.TypeName = ut
		case "Description":
			e.Description = stringLit(kv.Value)
		case "Story":
			e.Story = stringLit(kv.Value)
		case "Other":
			e.Other = stringLit(kv.Value)
		case "Fields":
			fields, err := extractFields(kv.Value)
			if err != nil {
				return e, err
			}
			e.Fields = fields
		}
	}
	return e, nil
}

func extractFields(node ast.Expr) ([]fieldEntry, error) {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return nil, nil
	}
	var fields []fieldEntry
	for _, elt := range cl.Elts {
		inner, ok := elt.(*ast.CompositeLit)
		if !ok {
			continue
		}
		var f fieldEntry
		for _, felt := range inner.Elts {
			kv, ok := felt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			switch key.Name {
			case "Note":
				f.Note = stringLit(kv.Value)
			case "GoType":
				f.GoType = stringLit(kv.Value)
			case "Tale":
				f.Tale = stringLit(kv.Value)
			}
		}
		fields = append(fields, f)
	}
	return fields, nil
}

// stringLit extracts the string value from an AST expression.
// It handles plain string literals ("foo" or `foo`) and binary
// concatenation expressions ("foo" + "bar"), including multi-line chains.
func stringLit(node ast.Expr) string {
	switch v := node.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			s := v.Value
			if len(s) >= 2 {
				if s[0] == '`' {
					// Raw string literal — strip surrounding backticks.
					return s[1 : len(s)-1]
				}
				// Interpreted string literal — use strconv to handle escapes.
				if unquoted, err := strconv.Unquote(s); err == nil {
					return unquoted
				}
				return s[1 : len(s)-1]
			}
		}
	case *ast.BinaryExpr:
		// Handle string concatenation: "a" + "b" + "c".
		// The AST is left-associative so we recurse on both sides.
		if v.Op == token.ADD {
			return stringLit(v.X) + stringLit(v.Y)
		}
	}
	return ""
}

func identOrSel(node ast.Expr) string {
	switch v := node.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return v.Sel.Name
	}
	return ""
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

type validationError struct {
	MessageID string
	Field     string
	Msg       string
}

func (e validationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("message %q field %q: %s", e.MessageID, e.Field, e.Msg)
	}
	return fmt.Sprintf("message %q: %s", e.MessageID, e.Msg)
}

func validate(entries []underlierEntry, verbose bool) error {
	var errs []error
	seen := map[string]bool{}

	for _, e := range entries {
		if verbose {
			fmt.Printf("  validating %s %s (%s)\n",
				emoji(e.TypeName), e.Seed, e.MessageID)
		}

		if seen[e.MessageID] {
			errs = append(errs, validationError{e.MessageID, "", "duplicate MessageID"})
		}
		seen[e.MessageID] = true

		ut := e.TypeName

		isStatic := ut == enums.UnderlyingTypeStaticGeneral ||
			ut == enums.UnderlyingTypeStaticError ||
			ut == enums.UnderlyingTypeSentinelError ||
			ut == enums.UnderlyingTypeStaticErrorWrapper ||
			ut == enums.UnderlyingTypeStaticCobra

		isDynamic := ut == enums.UnderlyingTypeDynamicGeneral ||
			ut == enums.UnderlyingTypeDynamicError ||
			ut == enums.UnderlyingTypeDynamicErrorWrapper ||
			ut == enums.UnderlyingTypeDynamicCobra

		isWrapper := ut == enums.UnderlyingTypeStaticErrorWrapper ||
			ut == enums.UnderlyingTypeDynamicErrorWrapper

		hasFields := len(e.Fields) > 0

		// Static types must have no fields (except StaticWrapper which has
		// no Fields at all — Wrapped is implicit).
		if ut == enums.UnderlyingTypeStaticErrorWrapper && hasFields {
			errs = append(errs, validationError{e.MessageID, "Fields",
				"UnderlyingTypeStaticErrorWrapper must not declare Fields (Wrapped is implicit)"})
		} else if isStatic && !isWrapper && hasFields {
			errs = append(errs, validationError{e.MessageID, "Fields",
				"static type must not have Fields"})
		}

		// Dynamic types must have fields.
		if isDynamic && !hasFields {
			errs = append(errs, validationError{e.MessageID, "Fields",
				"dynamic type must have at least one Fields entry"})
		}

		// Count error-typed fields.
		var errorFields []fieldEntry
		for _, f := range e.Fields {
			if f.GoType == "error" {
				errorFields = append(errorFields, f)
			}
		}
		if len(errorFields) > 1 {
			errs = append(errs, validationError{e.MessageID, "Fields",
				"at most one field may have GoType \"error\""})
		}
		if len(errorFields) == 1 {
			if errorFields[0].Note != "Wrapped" {
				errs = append(errs, validationError{e.MessageID, errorFields[0].Note,
					"the error-typed field must be named \"Wrapped\""})
			}
			if !isWrapper {
				errs = append(errs, validationError{e.MessageID, "Fields",
					"a field with GoType \"error\" is only permitted on wrapper types"})
			}
		}
		if isWrapper && ut != enums.UnderlyingTypeStaticErrorWrapper && len(errorFields) == 0 {
			errs = append(errs, validationError{e.MessageID, "Fields",
				"wrapper type must have a Fields entry {Note:\"Wrapped\", GoType:\"error\"}"})
		}

		// Validate {{.Token}} consistency.
		tokens := extractTemplateTokens(e.Other)
		fieldNames := map[string]bool{}
		for _, f := range e.Fields {
			if f.GoType != "error" {
				// error fields are represented as string in the template data.
				fieldNames[f.Note] = true
			} else {
				// Wrapped can appear as {{.Wrapped}} — it is valid.
				fieldNames[f.Note] = true
			}
		}
		for _, tok := range tokens {
			if !fieldNames[tok] {
				errs = append(errs, validationError{e.MessageID, tok,
					fmt.Sprintf("{{.%s}} in Other has no matching Fields entry", tok)})
			}
		}
		for name := range fieldNames {
			found := false
			for _, tok := range tokens {
				if tok == name {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, validationError{e.MessageID, name,
					fmt.Sprintf("Fields entry %q has no matching {{.%s}} in Other", name, name)})
			}
		}

		if e.Seed == "" {
			errs = append(errs, validationError{e.MessageID, "Seed", "Seed must not be empty"})
		}
		if e.MessageID == "" {
			errs = append(errs, validationError{e.MessageID, "MessageID", "MessageID must not be empty"})
		}
		if ut == enums.UnderlyingTypeUndefined {
			errs = append(errs, validationError{e.MessageID, "TypeName", "TypeName must not be UnderlyingTypeUndefined"})
		}
	}

	if len(errs) == 0 {
		return nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("lingo: %d validation error(s) found — no files written:\n", len(errs)))
	for _, err := range errs {
		sb.WriteString("  • ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return errors.New(sb.String())
}

var templateTokenRe = regexp.MustCompile(`\{\{\.([A-Za-z_][A-Za-z0-9_]*)\}\}`)

func extractTemplateTokens(s string) []string {
	matches := templateTokenRe.FindAllStringSubmatch(s, -1)
	var tokens []string
	for _, m := range matches {
		tokens = append(tokens, m[1])
	}
	return tokens
}

// ---------------------------------------------------------------------------
// Code generation
// ---------------------------------------------------------------------------

func generate(dir, pkgName, baseStruct string, entries []underlierEntry) error {
	cobra, general, errs := splitEntries(entries)

	cobraFile, err := generateCobra(pkgName, baseStruct, cobra)
	if err != nil {
		return fmt.Errorf("generating cobra file: %w", err)
	}
	generalFile, err := generateGeneral(pkgName, baseStruct, general)
	if err != nil {
		return fmt.Errorf("generating general file: %w", err)
	}
	errorsFile, err := generateErrors(pkgName, baseStruct, errs)
	if err != nil {
		return fmt.Errorf("generating errors file: %w", err)
	}

	for path, src := range map[string][]byte{
		filepath.Join(dir, "messages-cobra-auto.go"):   cobraFile,
		filepath.Join(dir, "messages-general-auto.go"): generalFile,
		filepath.Join(dir, "messages-errors-auto.go"):  errorsFile,
	} {
		if err := os.WriteFile(path, src, 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		fmt.Printf("lingo: wrote %s\n", path)
	}
	return nil
}

func splitEntries(entries []underlierEntry) (cobra, general, errs []underlierEntry) {
	for _, e := range entries {
		switch e.TypeName {
		case enums.UnderlyingTypeStaticCobra, enums.UnderlyingTypeDynamicCobra:
			cobra = append(cobra, e)
		case enums.UnderlyingTypeStaticGeneral, enums.UnderlyingTypeDynamicGeneral:
			general = append(general, e)
		default:
			errs = append(errs, e)
		}
	}
	sort.Slice(cobra, func(i, j int) bool { return cobra[i].Seed < cobra[j].Seed })
	sort.Slice(general, func(i, j int) bool { return general[i].Seed < general[j].Seed })
	sort.Slice(errs, func(i, j int) bool { return errs[i].Seed < errs[j].Seed })
	return
}

// ---------------------------------------------------------------------------
// Template data passed to every render function.
// ---------------------------------------------------------------------------

// templateData is the unified context passed to every Go template.
// The template accesses derived names (StructName, ErrorStruct etc.) as
// pre-computed fields so that the templates themselves contain no logic
// beyond ranging over Fields.
type templateData struct {
	Seed        string
	MessageID   string
	Description string
	Other       string
	Base        string
	Fields      []fieldEntry

	StructName  string
	ErrorTD     string
	ErrorStruct string
	Params      string

	// StructComment is the doc comment for the XxxTemplData struct used in
	// cobra, general, and dynamic error templates.
	StructComment string

	// ErrorTDComment is the doc comment for the XxxErrorTemplData struct
	// used in static, core, and static-wrapper error templates.
	ErrorTDComment string

	// ErrorComment is the doc comment for the XxxError struct, used in
	// all error templates.
	ErrorComment string
}

func newTemplateData(e underlierEntry, base string) templateData {
	nef := nonErrorFields(e.Fields)

	structName := e.Seed + "TemplData"
	errorTD := e.Seed + "ErrorTemplData"
	errorStruct := e.Seed + "Error"

	return templateData{
		Seed:           e.Seed,
		MessageID:      e.MessageID,
		Description:    e.Description,
		Other:          goStringLit(e.Other),
		Base:           base,
		Fields:         nef,
		StructName:     structName,
		ErrorTD:        errorTD,
		ErrorStruct:    errorStruct,
		Params:         renderParams(nef),
		StructComment:  structDocComment(structName, e.Description),
		ErrorTDComment: structDocComment(errorTD, e.Description),
		ErrorComment:   structDocComment(errorStruct, e.Description),
	}
}

// structDocComment returns a wrapped doc comment for a struct. When desc is
// non-empty the comment reads "// <typeName> <desc>.", wrapped at 80 chars.
// When desc is empty a 🔥 TODO reminder is returned instead.
func structDocComment(typeName, desc string) string {
	if desc == "" {
		return fmt.Sprintf("// 🔥 %s TODO: add struct Description", typeName)
	}
	return wrapComment(typeName+" "+desc+".", "// ", 80)
}

// ---------------------------------------------------------------------------
// Template definitions
// ---------------------------------------------------------------------------

// execTemplate parses text as a Go template, executes it with data, and
// returns the resulting string.  A non-nil error means the template itself
// is malformed — this is always a programmer error in lingo.
func execTemplate(name, text string, data any) (string, error) {
	t, err := template.New(name).Funcs(tmplFuncs).Parse(text)
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", name, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %q: %w", name, err)
	}
	return buf.String(), nil
}

// tmplFuncs are the custom functions available inside every template.
var tmplFuncs = template.FuncMap{
	// lower returns the first character of s in lower case.
	"lower": lowerFirst,
	// wrap word-wraps text at width characters, prefixing every line
	// with prefix. Used in templates to keep doc comments within 80 chars.
	"wrap": func(text, prefix string, width int) string {
		return wrapComment(text, prefix, width)
	},
}

// tmplCobra generates a cobra short/long message entry.
// Dynamic cobra messages (Fields non-empty) include a NewXxxTemplData
// constructor; static messages do not.
const tmplCobra = `{{.StructComment}}
type {{.StructName}} struct {
	{{.Base}}
{{- range .Fields}}
{{- if .Tale}}
	// {{.Note}} {{.Tale}}
{{- else}}
	// 🔥 {{.Note}} TODO: add field Tale
{{- end}}
	{{.Note}} {{.GoType}}
{{- end}}
}

{{wrap (printf "Message returns the i18n message for %s." .StructName) "// " 80}}
func (td {{.StructName}}) Message() *i18n.Message {
	return &i18n.Message{
		ID:          {{printf "%q" .MessageID}},
		Description: {{printf "%q" .Description}},
		Other:       {{.Other}},
	}
}
{{if .Fields}}
{{wrap (printf "New%s creates a new %s." .StructName .StructName) "// " 80}}
func New{{.StructName}}({{.Params}}) {{.StructName}} {
	return {{.StructName}}{
		{{.Base}}: {{.Base}}{},
{{- range .Fields}}
		{{.Note}}: {{lower .Note}},
{{- end}}
	}
}
{{end}}`

// tmplGeneral generates a general (non-error) message entry.
// Dynamic general messages include a NewXxxTemplData constructor.
const tmplGeneral = `{{.StructComment}}
type {{.StructName}} struct {
	{{.Base}}
{{- range .Fields}}
{{- if .Tale}}
	// {{.Note}} {{.Tale}}
{{- else}}
	// 🔥 {{.Note}} TODO: add field Tale
{{- end}}
	{{.Note}} {{.GoType}}
{{- end}}
}

{{wrap (printf "Message returns the i18n message for %s." .StructName) "// " 80}}
func (td {{.StructName}}) Message() *i18n.Message {
	return &i18n.Message{
		ID:          {{printf "%q" .MessageID}},
		Description: {{printf "%q" .Description}},
		Other:       {{.Other}},
	}
}
{{if .Fields}}
{{wrap (printf "New%s creates a new %s." .StructName .StructName) "// " 80}}
func New{{.StructName}}({{.Params}}) {{.StructName}} {
	return {{.StructName}}{
		{{.Base}}: {{.Base}}{},
{{- range .Fields}}
		{{.Note}}: {{lower .Note}},
{{- end}}
	}
}
{{end}}`

// tmplErrorStatic generates a static error:
//
//	XxxErrorTemplData  +  XxxError  +  ErrXxx sentinel
const tmplErrorStatic = `{{.ErrorTDComment}}
type {{.ErrorTD}} struct {
	{{.Base}}
}

// Message creates a new i18n message using the template data.
func (td {{.ErrorTD}}) Message() *i18n.Message {
	return &i18n.Message{
		ID:          {{printf "%q" .MessageID}},
		Description: {{printf "%q" .Description}},
		Other:       {{.Other}},
	}
}

{{.ErrorComment}}
type {{.ErrorStruct}} struct {
	li18ngo.LocalisableError
}

{{wrap (printf "Err%s is the exported sentinel error for %s." .Seed .ErrorStruct) "// " 80}}
var Err{{.Seed}} = {{.ErrorStruct}}{
	LocalisableError: li18ngo.LocalisableError{
		Data: {{.ErrorTD}}{},
	},
}
`

// tmplErrorCore generates a core sentinel error:
//
//	XxxErrorTemplData  +  XxxError  +  ErrXxx (exported)
//
// Callers use errors.Is(err, locale.ErrXxx) directly.
const tmplErrorCore = `{{.ErrorTDComment}}
type {{.ErrorTD}} struct {
	{{.Base}}
}

// Message creates a new i18n message using the template data.
func (td {{.ErrorTD}}) Message() *i18n.Message {
	return &i18n.Message{
		ID:          {{printf "%q" .MessageID}},
		Description: {{printf "%q" .Description}},
		Other:       {{.Other}},
	}
}

{{.ErrorComment}}
{{wrap (printf "Use errors.Is(err, locale.Err%s) to test for this error." .Seed) "// " 80}}
type {{.ErrorStruct}} struct {
	li18ngo.LocalisableError
}

{{wrap (printf "Err%s is the exported sentinel for %s." .Seed .ErrorStruct) "// " 80}}
var Err{{.Seed}} = {{.ErrorStruct}}{
	LocalisableError: li18ngo.LocalisableError{
		Data: {{.ErrorTD}}{},
	},
}
`

// tmplErrorStaticWrapper generates a static wrapping error:
//
//	XxxErrorTemplData  +  XxxError{Wrapped error}  +
//	Error()/Unwrap()   +  NewXxxError  +  AsXxxError
const tmplErrorStaticWrapper = `{{.ErrorTDComment}}
type {{.ErrorTD}} struct {
	{{.Base}}
}

// Message creates a new i18n message using the template data.
func (td {{.ErrorTD}}) Message() *i18n.Message {
	return &i18n.Message{
		ID:          {{printf "%q" .MessageID}},
		Description: {{printf "%q" .Description}},
		Other:       {{.Other}},
	}
}

{{.ErrorComment}}
type {{.ErrorStruct}} struct {
	li18ngo.LocalisableError
	Wrapped error
}

// Error returns the combined wrapped and localised error message.
func (e {{.ErrorStruct}}) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

// Unwrap returns the wrapped error.
func (e {{.ErrorStruct}}) Unwrap() error {
	return e.Wrapped
}

{{wrap (printf "New%s creates a new %s wrapping wrapped." .ErrorStruct .ErrorStruct) "// " 80}}
func New{{.ErrorStruct}}(wrapped error) error {
	return &{{.ErrorStruct}}{
		LocalisableError: li18ngo.LocalisableError{Data: {{.ErrorTD}}{}},
		Wrapped:          wrapped,
	}
}
`

// tmplErrorDynamic generates a dynamic error with no wrapping:
//
//	XxxTemplData  +  XxxError  +  NewXxxError  +  AsXxxError
const tmplErrorDynamic = `{{.StructComment}}
type {{.StructName}} struct {
	{{.Base}}
{{- range .Fields}}
{{- if .Tale}}
	// {{.Note}} {{.Tale}}
{{- else}}
	// 🔥 {{.Note}} TODO: add field Tale
{{- end}}
	{{.Note}} {{.GoType}}
{{- end}}
}

// Message creates a new i18n message using the template data.
func (td {{.StructName}}) Message() *i18n.Message {
	return &i18n.Message{
		ID:          {{printf "%q" .MessageID}},
		Description: {{printf "%q" .Description}},
		Other:       {{.Other}},
	}
}

{{.ErrorComment}}
type {{.ErrorStruct}} struct {
	li18ngo.LocalisableError
{{- range .Fields}}
	{{.Note}} {{.GoType}}
{{- end}}
}

{{wrap (printf "New%s creates a new %s." .ErrorStruct .ErrorStruct) "// " 80}}
func New{{.ErrorStruct}}({{.Params}}) error {
	return &{{.ErrorStruct}}{
		LocalisableError: li18ngo.LocalisableError{
			Data: {{.StructName}}{
{{- range .Fields}}
				{{.Note}}: {{lower .Note}},
{{- end}}
			},
		},
{{- range .Fields}}
		{{.Note}}: {{lower .Note}},
{{- end}}
	}
}
`

// tmplErrorDynamicWrapper generates a dynamic wrapping error:
//
//	XxxTemplData (Wrapped as string)  +  XxxError (Wrapped as error)  +
//	Error()/Unwrap()  +  NewXxxError(wrapped, fields...)  +  AsXxxError
//
// The Wrapped field is stored as string in the template data (for
// go-i18n interpolation via {{.Wrapped}}) and as error in the error
// struct (for unwrapping via errors.As / errors.Is).
const tmplErrorDynamicWrapper = `{{.StructComment}}
type {{.StructName}} struct {
	{{.Base}}
{{- range .Fields}}
{{- if .Tale}}
	// {{.Note}} {{.Tale}}
{{- else}}
	// 🔥 {{.Note}} TODO: add field Tale
{{- end}}
	{{.Note}} {{.GoType}}
{{- end}}
	// Wrapped is the string representation of the wrapped error,
	// used for go-i18n template interpolation via {{"{{"}} .Wrapped {{"}}"}}.
	Wrapped string
}

// Message creates a new i18n message using the template data.
func (td {{.StructName}}) Message() *i18n.Message {
	return &i18n.Message{
		ID:          {{printf "%q" .MessageID}},
		Description: {{printf "%q" .Description}},
		Other:       {{.Other}},
	}
}

{{.ErrorComment}}
type {{.ErrorStruct}} struct {
	li18ngo.LocalisableError
	Wrapped error
{{- range .Fields}}
	{{.Note}} {{.GoType}}
{{- end}}
}

// Error returns the combined wrapped and localised error message.
func (e {{.ErrorStruct}}) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

// Unwrap returns the wrapped error.
func (e {{.ErrorStruct}}) Unwrap() error {
	return e.Wrapped
}

{{wrap (printf "New%s creates a new %s wrapping wrapped." .ErrorStruct .ErrorStruct) "// " 80}}
func New{{.ErrorStruct}}(wrapped error, {{.Params}}) error {
	return &{{.ErrorStruct}}{
		LocalisableError: li18ngo.LocalisableError{
			Data: {{.StructName}}{
{{- range .Fields}}
				{{.Note}}: {{lower .Note}},
{{- end}}
				Wrapped: wrapped.Error(),
			},
		},
		Wrapped: wrapped,
{{- range .Fields}}
		{{.Note}}: {{lower .Note}},
{{- end}}
	}
}
`

// ---------------------------------------------------------------------------
// Per-file generators
// ---------------------------------------------------------------------------

func formatSource(src string) ([]byte, error) {
	b, err := format.Source([]byte(src))
	if err != nil {
		// Return the unformatted source alongside the error so the caller
		// can write it for debugging.
		return []byte(src), fmt.Errorf("formatting generated source: %w", err)
	}
	return b, nil
}

func renderHeader(pkg string, imports []string) string {
	var sb strings.Builder
	sb.WriteString("// Code generated by lingo. DO NOT EDIT.\n")
	sb.WriteString("// Re-generate by running: go generate\n\n")
	sb.WriteString(fmt.Sprintf("package %s\n\n", pkg))
	sb.WriteString("import (\n")
	for _, imp := range imports {
		sb.WriteString(fmt.Sprintf("\t%q\n", imp))
	}
	sb.WriteString(")\n")
	return sb.String()
}

// ---------------------------------------------------------------------------
// Banner helpers
// ---------------------------------------------------------------------------

func emoji(ut enums.UnderlyingType) string {
	switch ut {
	case enums.UnderlyingTypeStaticCobra, enums.UnderlyingTypeDynamicCobra:
		return "🧊"
	case enums.UnderlyingTypeStaticGeneral, enums.UnderlyingTypeDynamicGeneral:
		return "📨"
	default:
		return "❌"
	}
}

// wrapComment wraps text at maxWidth, prefixing every line with prefix.
func wrapComment(text, prefix string, maxWidth int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	var lines []string
	line := prefix
	for _, w := range words {
		candidate := line + w
		if len(line) > len(prefix) {
			candidate = line + " " + w
		}
		if len(candidate) > maxWidth && len(line) > len(prefix) {
			lines = append(lines, line)
			line = prefix + w
		} else {
			if len(line) > len(prefix) {
				line += " " + w
			} else {
				line += w
			}
		}
	}
	if len(line) > len(prefix) {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func banner(e underlierEntry) string {
	const rule = "// ============================================================================="
	em := emoji(e.TypeName)
	story := e.Story
	if story == "" {
		story = fmt.Sprintf("🔥 %s TODO: add message Story", e.Seed)
	}
	wrapped := wrapComment(story, "// ", 80)
	return fmt.Sprintf("%s\n// %s %s\n//\n%s\n%s", rule, em, e.Seed, wrapped, rule)
}

// ---------------------------------------------------------------------------
// Cobra generation
// ---------------------------------------------------------------------------

func generateCobra(pkg, base string, entries []underlierEntry) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString(renderHeader(pkg, []string{
		"github.com/nicksnyder/go-i18n/v2/i18n",
	}))
	for _, e := range entries {
		sb.WriteString("\n")
		sb.WriteString(banner(e))
		sb.WriteString("\n\n")
		out, err := execTemplate("cobra", tmplCobra, newTemplateData(e, base))
		if err != nil {
			return nil, fmt.Errorf("cobra entry %q: %w", e.Seed, err)
		}
		sb.WriteString(out)
	}
	return formatSource(sb.String())
}

// ---------------------------------------------------------------------------
// General generation
// ---------------------------------------------------------------------------

func generateGeneral(pkg, base string, entries []underlierEntry) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString(renderHeader(pkg, []string{
		"github.com/nicksnyder/go-i18n/v2/i18n",
	}))
	for _, e := range entries {
		sb.WriteString("\n")
		sb.WriteString(banner(e))
		sb.WriteString("\n\n")
		out, err := execTemplate("general", tmplGeneral, newTemplateData(e, base))
		if err != nil {
			return nil, fmt.Errorf("general entry %q: %w", e.Seed, err)
		}
		sb.WriteString(out)
	}
	return formatSource(sb.String())
}

// ---------------------------------------------------------------------------
// Error generation
// ---------------------------------------------------------------------------

func generateErrors(pkg, base string, entries []underlierEntry) ([]byte, error) {
	needFmt := false
	for _, e := range entries {
		switch e.TypeName {
		case enums.UnderlyingTypeStaticErrorWrapper,
			enums.UnderlyingTypeDynamicErrorWrapper:
			needFmt = true
		}
	}
	imports := []string{
		"errors",
		"github.com/nicksnyder/go-i18n/v2/i18n",
		"github.com/snivilised/li18ngo",
	}
	if needFmt {
		imports = append([]string{"fmt"}, imports...)
	}

	var sb strings.Builder
	sb.WriteString(renderHeader(pkg, imports))
	for _, e := range entries {
		sb.WriteString("\n")
		sb.WriteString(banner(e))
		sb.WriteString("\n\n")
		out, err := renderErrorEntry(e, base)
		if err != nil {
			return nil, err
		}
		sb.WriteString(out)
	}
	return formatSource(sb.String())
}

func renderErrorEntry(e underlierEntry, base string) (string, error) {
	td := newTemplateData(e, base)
	switch e.TypeName {
	case enums.UnderlyingTypeStaticError:
		return execTemplate("errorStatic", tmplErrorStatic, td)
	case enums.UnderlyingTypeSentinelError:
		return execTemplate("errorCore", tmplErrorCore, td)
	case enums.UnderlyingTypeStaticErrorWrapper:
		return execTemplate("errorStaticWrapper", tmplErrorStaticWrapper, td)
	case enums.UnderlyingTypeDynamicError:
		return execTemplate("errorDynamic", tmplErrorDynamic, td)
	case enums.UnderlyingTypeDynamicErrorWrapper:
		return execTemplate("errorDynamicWrapper", tmplErrorDynamicWrapper, td)
	default:
		return fmt.Sprintf("// unknown type %q for %q\n", e.TypeName, e.Seed), nil
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func nonErrorFields(fields []fieldEntry) []fieldEntry {
	var out []fieldEntry
	for _, f := range fields {
		if f.GoType != "error" {
			out = append(out, f)
		}
	}
	return out
}

// renderParams renders a parameter list from the given fields.
func renderParams(fields []fieldEntry) string {
	var parts []string
	for _, f := range fields {
		parts = append(parts, fmt.Sprintf("%s %s", lowerFirst(f.Note), f.GoType))
	}
	return strings.Join(parts, ", ")
}

// lowerFirst returns s with its first rune lowercased.
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// goStringLit returns a Go string literal for s, using backtick form when
// the string contains double quotes or newlines.
func goStringLit(s string) string {
	if strings.ContainsAny(s, "\"\n\t") {
		return "`" + s + "`"
	}
	return fmt.Sprintf("%q", s)
}
