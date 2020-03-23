package semtools

import (
	"fmt"
	"regexp"
	"strings"
	"sort"
)

// TurtleParserDefaultNamespace defines a couple
// of defaults for namespaces and substitution during
// parsing. These can e.g. also be missing in the
// ttl that should be Unmarshaled and the re-substitution
// will not fail.
var TurtleParserDefaultNamespace = NewNamespace()

// TurtleParserOptions are options that configure
// the parser with some settings
type TurtleParserOptions struct {

	// Substitute will be used during Marshal of
	// statements into string format. If this is true,
	// the ttl will be substituted with the data
	// from the given namespace and a potential base iri.
	Substitute bool

	// Namespace are additional namespace definitions
	// to use during substitution. The TurtleParserDefaultNamespace
	// will always be used. This namespace, after Unmarshal
	// will as well be updated to all namespaces found within
	// the turtle definition.
	Namespace *Namespace

	// PrettyPrint will configure the Marshal function
	// to pretty print the output with nice indents etc.
	PrettyPrint bool

	// RequireBaseIri will configure the Unmarshal function
	// to require a base iri to be present in the ttl.
	RequireBaseIri bool

	// FallbackToFirstSubjectForBaseIri will allow the unmarshal
	// to fallback on the first subject if no base tag is found.
	FallbackToFirstSubjectForBaseIri bool

}

// TurtleParser is a entity compatible with Parser
// that works with text/ttl content. It can parse
// Graphs to and from the turtle representation.
type TurtleParser struct {

	// options contains the runtime options to
	// apply during parsing
	options *TurtleParserOptions

}

// NewTurtleParser creates a new turtle parser with the given
// options
func NewTurtleParser(opts *TurtleParserOptions) *TurtleParser {
	// make sure things are initialized
	if opts == nil {
		opts = &TurtleParserOptions{}
	}
	if opts.Namespace == nil {
		opts.Namespace = NewEmptyNamespace()
	}
	// create new parser
	return &TurtleParser{
		options: opts,
	}
}

// Marshal creates a text/ttl representation from
// the provided statements.
func (p *TurtleParser) Marshal(stmts []Statement) (string, error) {
	// get full namespace
	ns := TurtleParserDefaultNamespace.Include(p.options.Namespace)

	// setup result
	ttl := ""

	// check if statements contain more than one
	// graph - if not, we can use that as the
	// base iri in the ttl.
	baseIri := ""
	unique := true
	for _, stmt := range stmts {
		if stmt.Graph() != nil {
			if baseIri != "" && baseIri != stmt.Graph().Iri() {
				unique = false
			}
			baseIri = stmt.Graph().Iri()
		}
	}
	if unique && baseIri != "" {
		// add into namespace
		// and set base directive in ttl
		ns.Set("", baseIri)
		ttl = "@base " + p.marshalIri(baseIri, false, nil) + " .\n"
	}

	// add @prefix directives if we're substituting
	if p.options.Substitute {
		keys := ns.ListKeys()
		sort.Strings(keys)
		for _, k := range keys {
			ttl = ttl + "@prefix " + k + ": " + p.marshalIri(ns.MustGet(k) + "#", false, nil) + " .\n"
		}
	}

	// we'll only use the vertices as that's the information
	// we store in ttl format.
	// therefore we first order the vertices by subject
	// and predicate
	grouped := make(map[string]map[string][]Node)
	for _, v := range stmts {
		if _, ok := grouped[v.Subject().Iri()]; !ok {
			grouped[v.Subject().Iri()] = make(map[string][]Node)
		}
		if _, ok := grouped[v.Subject().Iri()][v.Predicate().Iri()]; !ok {
			grouped[v.Subject().Iri()][v.Predicate().Iri()] = make([]Node, 0)
		}
		grouped[v.Subject().Iri()][v.Predicate().Iri()] = append(grouped[v.Subject().Iri()][v.Predicate().Iri()], v.Object())
	}

	// now we can produce ttl grouped by subject and predicate
	// for this we'll have to access the subjects in the grouped
	// map in a sorted manner
	subjectIris := make([]string, 0)
	for k, _ := range grouped {
		subjectIris = append(subjectIris, k)
	}
	sort.Strings(subjectIris)
	for _, subjectIri := range subjectIris {

		// add potential pretty print comment
		if p.options.PrettyPrint {
			ttl = ttl + "\n# " + subjectIri + "\n"
		}

		// add marshaled subject
		ms, err := p.marshalNode(NewNamedNode(subjectIri), p.options.Substitute, ns)
		if err != nil {
			return "", err
		}
		ttl = ttl + ms + " "

		// go through predicates (sorted) and add the data
		// for the current subject
		predicateIris := make([]string, 0)
		for k, _ := range grouped[subjectIri] {
			predicateIris = append(predicateIris, k)
		}
		sort.Strings(predicateIris)
		for pIdx, predicateIri := range predicateIris {

			// make sure to separate previous predicate
			if pIdx > 0 {
				ttl = ttl + "; "
			}

			// add potential pretty print spacing
			if p.options.PrettyPrint {
				ttl = ttl + "\n    "
			}

			// add marshaled predicate
			ps, _ := p.marshalNode(NewNamedNode(predicateIri), p.options.Substitute, ns)
			if err != nil {
				return "", err
			}
			ttl = ttl + ps + " "

			// go through objects (sorted) and add the targets
			// for the current subject - predicate combination
			objs := grouped[subjectIri][predicateIri]
			sort.Slice(objs, func(i, j int) bool {
			  return objs[i].String() < objs[j].String()
			})
			for idx, obj := range objs {

				// make sure to separate previous value
				if idx > 0 {
					ttl = ttl + ", "
				}

				// add potential pretty print spacing
				if p.options.PrettyPrint {
					ttl = ttl + "\n        "
				}

				// add marshaled object
				vs, _ := p.marshalNode(obj, p.options.Substitute, ns)
				if err != nil {
					return "", err
				}
				ttl = ttl + vs + " "

			}

		}

		// close the subjects statement
		ttl = ttl + ".\n"

	}

	// return with not error
	return ttl, nil
}

// Unmarshal creates statements from the given
// text/ttl data.
func (p *TurtleParser) Unmarshal(str string) ([]Statement, error) {

	// clean ttl string
	str = p.stripComments(str)

	// extract body from ttl
	ttlDirectiveMatcher := regexp.MustCompile(`\s*@(prefix|base)\s+.+?>\s*\.`)
	ttl_body := strings.ReplaceAll(ttlDirectiveMatcher.ReplaceAllString(str, ""), "\n", "")

	// prepare namespace
	ns := TurtleParserDefaultNamespace.Copy()
	extractedNs, err := p.extractNamespace(str)
	if err != nil {
		return nil, err
	}
	ns = ns.Include(extractedNs)

	// get base uri (if possible)
	baseIri, err := p.extractBaseIri(str, ns, p.options.FallbackToFirstSubjectForBaseIri)
	if err != nil {
		return nil, err
	}
	var graph NamedNode
	if baseIri != "" {
		// set graph to be used in statements
		graph = NewNamedNode(baseIri)
		// add base iri as well to prefixes
		ns.Set("", graph.Iri())
	}

	// initialize result slice
	result := []Statement{}

	// extract statements to then parse into nodes and vertices
	leadingIriMatcher := regexp.MustCompile(`^(a|<.*?>|\S*?:\S+?)\s+`)
	for _, statement := range p.extractStatements(ttl_body) {

		// extract subject of statement
		sm := leadingIriMatcher.FindStringSubmatch(statement)
		if len(sm) != 2 {
			return nil, fmt.Errorf("Unable to extract subject from statement: %v", statement)
		}

		// create named node for subject if not exist
		subjIri := p.unmarshalIri(sm[1], ns)
		subject := NewNamedNode(subjIri)

		// extract the substatements for the current subject
		statement = strings.TrimSpace(leadingIriMatcher.ReplaceAllString(statement, ""))
		for _, subStatement := range p.extractSubStatements(statement) {

			// extract predicate of sub statement
			pm := leadingIriMatcher.FindStringSubmatch(subStatement)
			if len(pm) != 2 {
				return nil, fmt.Errorf("Unable to extract predicate from substatement: %v", subStatement)
			}

			// create named node for predicate if not exist
			predIri := p.unmarshalIri(pm[1], ns)
			predicate := NewNamedNode(predIri)

			// finally, let's parse the objects
			// - first remove the leading predicate definition and trim whitespace again
			subStatement = strings.TrimSpace(leadingIriMatcher.ReplaceAllString(subStatement, ""))
			// - now go through all object strings, unmarshal them and add as vertice
			for _, objectStr := range p.extractObjects(subStatement) {

				// unmarshal the string into a node
				object, err := p.unmarshalNode(objectStr, ns)
				if err != nil {
					return nil, fmt.Errorf("Unable to unmarshal object string: %v", objectStr)
				}

				// check if we have already the same statemetn in the results
				new := true
				s := NewStatement(subject, predicate, object, graph)
				for _, r := range result {
					if r.Equals(s) {
						new = false
						break;
					}
				}
				if new {
					result = append(result, s)
				}

			}

		}

	}


	return result, nil
}

// marshalNode creates the string representation of a single node
// applying potential substitution. For substitution to be used
// set the flag in the parameters, and make sure to provide a
// namespace in the ns-parameter
func (p *TurtleParser) marshalNode(n Node, substitute bool, ns *Namespace) (string, error) {


	// switch by type and marshal
	switch n.(type) {

	case NamedNode:

		// in case of named nodes, we have to substitute the
		// iri in case it's requird, otherwise just put it between
		// <...> brackets, unless substituted that is...
		iri := p.marshalIri(n.(NamedNode).Iri(), substitute, ns)
		return iri, nil

	case LocalizedLiteral:

		// for language nodes we have to parse it into a literal string
		// using the ""@lang format. Any "-characters within the value
		// must be escaped with a forward slash.
		ln := n.(LocalizedLiteral)
		v := p.marshalString(ln.Value().(string))
		l := ln.Language()
		if l == "" {
			l = "default"
		}
		return v + "@" + l, nil

	case TypedLiteral:

		// for language nodes we have to parse it into a literal string
		// using the ""@lang format. Any "-characters within the value
		// must be escaped with a forward slash.
		tn := n.(TypedLiteral)
		v := p.marshalString(tn.String())
		iri := p.marshalIri(tn.Type().Iri(), substitute, ns)
		return v + "^^" + iri, nil

	default:
		// we've come across a type that is unknown ...
		return "", fmt.Errorf("Unable to marshal Node '%t'", n)
	}

}

// marshalIri returns a ttl string version of the iri, potentially substituted by
// the namespace if set via the parameter.
func (p *TurtleParser) marshalIri(iri string, substitute bool, ns *Namespace) string {

	// try substitution
	substituted := false
	if substitute {
		for _, k := range ns.ListKeys() {
			if strings.Contains(iri, ns.MustGet(k) + "#") {
				substituted = true
				iri = strings.Replace(iri, ns.MustGet(k) + "#", k + ":", 1)
				break;
			}
		}
	}

	// add brackets if not substituted
	if !substituted {
		iri = "<" + iri + ">"
	}

	// return special case - if we have
	// a substituted iri, we can replace
	// "rdf:type" with "a"
	if iri == "rdf:type" {
		return "a"
	}
	// return marshaled iri
	return iri

}

// marshalString returns an escaped ttl string including the wrapping
// apostrophes.
var ttlStringEscapeRe = regexp.MustCompile(`([^\\])(")`)  // cache compilation of regex
func (p *TurtleParser) marshalString(str string) string {

	// replace all " with a \" unless they already have
	// a \ infront of the "
	return "\"" + ttlStringEscapeRe.ReplaceAllString(str, "$1\\\"") + "\""

}

// unmarshalNode creates the node from the given string.
// The string can be either a literal, or another iri,
// the function will return the respective Node type.
func (p *TurtleParser) unmarshalNode(str string, ns *Namespace) (Node, error) {

	// make sure to remove all whitespace
	str = strings.TrimSpace(str)

	// switch type
	if str[0] == '"' {

		// it's an object literal
		if str[len(str) - 1] == '"' {

			// plain string -> defaults to language literal
			return NewLocalizedLiteral(str[1:len(str) - 1], "default"), nil

		} else if strings.Contains(str, "@") {

			// interpret as language literal
			idx := strings.LastIndex(str, "@")
			return NewLocalizedLiteral(str[1:idx - 1], str[idx + 1:len(str)]), nil

		} else if strings.Contains(str, "^^") {

			// interpret as typed literal
			idx := strings.LastIndex(str, "^^")
			tpIri := p.unmarshalIri(str[idx + 2:len(str)], ns)
			return NewTypedLiteral(str[1:idx - 1], NewNamedNode(tpIri)), nil

		}

	} else {

		// it's an iri so we return a named node
		iri := p.unmarshalIri(str, ns)
		return NewNamedNode(iri), nil

	}

	// getting here means it wasn't successful
	return nil, fmt.Errorf("Failed to unmarshal node from ttl string: %v", str)

}

// unmarshalIri produces a fully qualified iri from the string
// using the namespace.
func (p *TurtleParser) unmarshalIri(str string, ns *Namespace) string {

	// trim whitespace to be sure
	str = strings.TrimSpace(str)

	// handle special cases
	if str == "a" {
		str = "<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>"
	}

	if str[0] == '<' {

		// fully qualified already
		str = str[1:len(str) - 1]

	} else {

		// undo substitution
		sSplits := strings.Split(str, ":")
		if full, ok := ns.Get(sSplits[0]); ok && len(sSplits) == 2 {
			str = full + "#" + sSplits[1]
		}

	}

	// return the fully qualified iri
	return str
}

// extractBaseIri tries to extract a '@base <...> .' directive from the string,
// if it is not found and fallbackToFirstSubject is true, the first found subject
// will be used as the base uri.
// If no base uri could be extracted (e.g. also no actual statement in the str) the
// function returns an error.
var ttlBaseDirective = regexp.MustCompile(`(^|\s|\.)\s*@base\s*?<(.+?)>\s*?\.`)  // cache compilation of regex
var ttlFirstSubjectMatcher = regexp.MustCompile(`(^|\n|\.)\s*?(<.*?>|\S*?:\S+?)\s+`)  // cache compilation of regex
func (p *TurtleParser) extractBaseIri(str string, ns *Namespace, fallbackToFirstSubject bool) (string, error) {

	// try to find base uri via directive
	res := ttlBaseDirective.FindAllStringSubmatch(str, -1)
	if len(res) > 1 {
		return "", fmt.Errorf("More than one '@base <...> .' directive found in ttl data.")
	} else if len(res) == 1 {
		return res[0][2], nil
	}

	if fallbackToFirstSubject {

		// let's find the first mentioned subject
		fs := ttlFirstSubjectMatcher.FindStringSubmatch(str)
		if len(fs) == 3 {

			// make sure it's unmarshaled
			return p.unmarshalIri(fs[2], ns), nil

		}

	}

	// if we get here, we've got to error out in case it was required
	if p.options.RequireBaseIri {
		if !fallbackToFirstSubject {
			return "", fmt.Errorf("No '@base <...> .' directive found in ttl data.")
		}
		return "", fmt.Errorf("Not one single subject found in ttl data to use as base uri.")
	}
	return "", nil

}

// extractNamespace tries to extract all '@prefix ...: <...> .' directives from the string,
// and return the respective namespace map.
var ttlPrefixDirective = regexp.MustCompile(`(^|\s|\.)\s*?@prefix\s+?(\S*?):\s*?<(.+?)(#*)>\s*?.`)  // cache compilation of regex
func (p *TurtleParser) extractNamespace(str string) (*Namespace, error) {

	// try to find prefixes uri via directive
	res := ttlPrefixDirective.FindAllStringSubmatch(str, -1)

	// merge results into prefix map
	ns := NewEmptyNamespace()
	for _, m := range res {
		ns.Set(m[2], m[3])
	}

	// done
	return ns, nil

}

// extractStatements tries to extract ttl statements from the string. A
// statement is one subject, with one or more predicates with one or more objects,
// so a statement is an ttl-aggregate of multiple "triples", starting with an 
// (potentially substituted) Iri and ending with a ".". Note that the
// statements returned, will NOT have a trailing "." and thus are not valid
// ttl anymore! This function generally expects a single-lined ttl body
// without any prefix or other directive whatsoever.
var ttlStatementMatcher = regexp.MustCompile(`(^|\n|\.)\s*?(<.*?>|[A-Za-z0-9_-]*?:\S+?)\s+`)  // cache compilation of regex
func (p *TurtleParser) extractStatements(str string) []string {
	// expects: `:Me a :Thing , :Otherthing ; :knowns :Friend . :You a :Thing .`

	// try to find leading iris via regex
	res := ttlStatementMatcher.FindAllStringSubmatchIndex(str, -1)

	// make results
	statements := make([]string, len(res))
	for iterIdx, r := range res {
		// start index will be the second group
		// within the match
		sIdx := r[2]
		// if the first character at the start index
		// is a '.', let's move ahead one to be clean
		if str[sIdx] == '.' {
			sIdx += 1
		}
		// end index of the statement will be the
		// next start index in the matches, or the
		// end of the ttl
		eIdx := len(str) - 1
		if iterIdx < len(res) - 1 {
			eIdx = res[iterIdx + 1][2]
		}
		// make sure to ignore trailing '.'
		if str[eIdx] == '.' {
			eIdx -= 1
		}
		// let's store the current statement
		statements[iterIdx] = strings.TrimSpace(str[sIdx:eIdx + 1])
	}

	// done
	return statements

}

// extractSubStatements works no a statement (as recieved from the extractStatements,
// WITHOUT the leading subject in the statement, so this has to be removed before!)
// and returns a list of sub statemetns leadig with their predicate and followed
// by the (list of) object(s). Note that the result will be not valid ttl anymore
// so this can't really be used anywhere except for in this context.
var ttlSubStatementMatcher = regexp.MustCompile(`(^|\n|;)\s*?(a|<.*?>|\S*?:\S+?)\s+`)  // cache compilation of regex
func (p *TurtleParser) extractSubStatements(str string) []string {
	// expects: `a :Thing , :Otherthing ; :knowns :Friend`
	
	// try to find leading iris via regex
	res := ttlSubStatementMatcher.FindAllStringSubmatchIndex(str, -1)

	// make results
	subStatements := make([]string, len(res))
	for iterIdx, r := range res {
		// start index will be the second group
		// within the match
		sIdx := r[2]
		// if the first character at the start index
		// is a ';', let's move ahead one to be clean
		if str[sIdx] == ';' {
			sIdx += 1
		}
		// end index of the statement will be the
		// next start index in the matches, or the
		// end of the ttl
		eIdx := len(str) - 1
		if iterIdx < len(res) - 1 {
			eIdx = res[iterIdx + 1][2]
		}
		// make sure to ignore trailing ';'
		if str[eIdx] == ';' {
			eIdx -= 1
		}
		// let's store the current statement
		subStatements[iterIdx] = strings.TrimSpace(str[sIdx:eIdx + 1])
	}

	// done
	return subStatements

}

// extractObjects works no a statement (as recieved from the extractSubStatements,
// WITHOUT the leading predicate in the substatement, so this has to be removed before!)
// and returns a list of object strings. Each object string may be a literal, or another
// iri.
var ttlLiteralMatcher = regexp.MustCompile(`\s*?"([\S\s]*?[^\\])*?"(((\^\^|@)(.*?))?)\s*?(,|$)`)  // cache compilation of regex
var ttlIriMatcher = regexp.MustCompile(`\s*?(<.*?>|\S*?:\S+?)\s*?(,|$)`)
func (p *TurtleParser) extractObjects(str string) []string {
	// expects: `:Thing , :Otherthing , "sdasds"@de`

	// initialize objects list
	objects := []string{}

	// extract literals
	for _, match := range ttlLiteralMatcher.FindAllStringSubmatch(str, -1) {
		match[0] = strings.TrimSpace(match[0])
		if match[0][len(match[0]) - 1] == ',' {
			match[0] = match[0][:len(match[0]) - 1]
		}
		objects = append(objects, match[0])
	}

	// extract iris (first remove the previously matched literals)
	str = ttlLiteralMatcher.ReplaceAllString(str, "")
	for _, match := range ttlIriMatcher.FindAllStringSubmatch(str, -1) {
		match[0] = strings.TrimSpace(match[0])
		if match[0][len(match[0]) - 1] == ',' {
			match[0] = match[0][:len(match[0]) - 1]
		}
		objects = append(objects, match[0])
	}

	return objects

}

// stripComments will remove all lines starting with # from the ttl content
var ttlCommentMatcher = regexp.MustCompile(`(^|\n)\s*?#.*`)  // cache compilation of regex
func (p *TurtleParser) stripComments(str string) string {

	// replace all occurences of comment matcher
	return ttlCommentMatcher.ReplaceAllString(str, "")

}
