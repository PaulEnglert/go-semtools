package semtools

import (
	"testing"
	"fmt"
	"github.com/sirupsen/logrus"
)


func init() {

	logrus.SetLevel(logrus.DebugLevel)

}


func TestTurtleParserOptions(t *testing.T) {

	opts := &TurtleParserOptions{
		Substitute: true,
		PrettyPrint: true,
	}

	customNs := NewEmptyNamespace()
	customNs.Set("test", "http://www.test.de/test")
	opts = &TurtleParserOptions{
		Namespace: customNs,
	}

	if _, ok := opts.Namespace.Get("test"); !ok {
		t.Errorf("TurtleParserOptions not keeping Namespace")
	}

}

func TestNewTurtleParser(t *testing.T) {

	p := NewTurtleParser(nil)
	if p == nil {
		t.Errorf("NewTurtleParser(nil) fails to return a new parser")
	}

	p = NewTurtleParser(&TurtleParserOptions{Substitute: true})
	if p == nil {
		t.Errorf("NewTurtleParser(opts) fails to return a new parser")
	}

}

func TestMarshalNode(t *testing.T) {

	nn1 := NewNamedNode("http://www.test.de/test")
	nn2 := NewNamedNode("http://www.w3.org/2001/XMLSchema#string")
	rdftype := NewNamedNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")

	ll1 := NewLocalizedLiteral("hy my name is \"Paul\"", "de")
	ll2 := NewLocalizedLiteral("ohh", "")

	tl1 := NewTypedLiteral("1", nn1)
	tl2 := NewTypedLiteral("2", nn2)

	p := NewTurtleParser(nil)
	ns := NewNamespace()


	// test
	s, err := p.marshalNode(nn1, false, ns)
	if err != nil || s != "<http://www.test.de/test>" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "<http://www.test.de/test>", s)
	}
	s, err = p.marshalNode(nn1, true, ns)
	if err != nil || s != "<http://www.test.de/test>" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "<http://www.test.de/test>", s)
	}
	s, err = p.marshalNode(nn2, false, ns)
	if err != nil || s != "<http://www.w3.org/2001/XMLSchema#string>" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "<http://www.w3.org/2001/XMLSchema#string>", s)
	}
	s, err = p.marshalNode(nn2, true, ns)
	if err != nil || s != "xsd:string" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "xsd:string", s)
	}
	s, err = p.marshalNode(ll1, false, ns)
	if err != nil || s != "\"hy my name is \\\"Paul\\\"\"@de" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "\"hy my name is \\\"Paul\\\"\"@de", s)
	}
	s, err = p.marshalNode(ll2, false, ns)
	if err != nil || s != "\"ohh\"@default" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "\"ohh\"@default", s)
	}
	s, err = p.marshalNode(tl1, true, ns)
	if err != nil || s != "\"1\"^^<http://www.test.de/test>" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "\"1\"^^<http://www.test.de/test>", s)
	}
	s, err = p.marshalNode(tl2, true, ns)
	if err != nil || s != "\"2\"^^xsd:string" {
		t.Errorf("marshalNode() expected '%v' but got '%v", "\"2\"^^xsd:string", s)
	}
	s, err = p.marshalNode(rdftype, true, ns)
	if err != nil || s != "a" {
		t.Errorf("marshalNode() expected '%v' but got '%v'", "a", s)
	}
	s, err = p.marshalNode(rdftype, false, ns)
	if err != nil || s != "<" + rdftype.Iri() + ">" {
		t.Errorf("marshalNode() expected '%v' but got '%v'", "<" + rdftype.Iri() + ">", s)
	}

}

func TestUnmarshalNode(t *testing.T) {

	p := NewTurtleParser(nil)
	ns := NewNamespace()


	// test
	n, err := p.unmarshalNode("<http://www.test.de/test>", ns)
	if err != nil || n.(NamedNode).Iri() != "http://www.test.de/test" {
		t.Errorf("unmarshalNode() go unexpected data")
	}
	n, err = p.unmarshalNode("xsd:string", ns)
	if err != nil || n.(NamedNode).Iri() != "http://www.w3.org/2001/XMLSchema#string" {
		t.Errorf("unmarshalNode() go unexpected data")
	}
	n, err = p.unmarshalNode("\"mys\\\"astring\"", ns)
	if err != nil || n.(LocalizedLiteral).Value() != "mys\\\"astring" || n.(LocalizedLiteral).Language() != "default" {
		t.Errorf("unmarshalNode() go unexpected data")
	}
	n, err = p.unmarshalNode("\"myst\\\"@ ,ring\"@de", ns)
	if err != nil || n.(LocalizedLiteral).Value() != "myst\\\"@ ,ring" || n.(LocalizedLiteral).Language() != "de" {
		fmt.Println(n.String())
		fmt.Println(n.(LocalizedLiteral).Language())
		t.Errorf("unmarshalNode() go unexpected data")
	}
	n, err = p.unmarshalNode("\"mystring\"^^<http://www.w3.org/2001/XMLSchema#string>", ns)
	if err != nil || n.(TypedLiteral).Value() != "mystring" || n.(TypedLiteral).Type().Iri() != "http://www.w3.org/2001/XMLSchema#string" {
		t.Errorf("unmarshalNode() go unexpected data")
	}
	n, err = p.unmarshalNode("\"mystring\"^^xsd:string", ns)
	if err != nil || n.(TypedLiteral).Value() != "mystring" || n.(TypedLiteral).Type().Iri() != "http://www.w3.org/2001/XMLSchema#string" {
		t.Errorf("unmarshalNode() go unexpected data")
	}
	n, err = p.unmarshalNode("\"http://myresources:80/resources/file.py\"^^<http://www.w3.org/2001/XMLSchema#anyURI>", ns)
	if err != nil || n.(TypedLiteral).Value() != "http://myresources:80/resources/file.py" || n.(TypedLiteral).Type().Iri() != "http://www.w3.org/2001/XMLSchema#anyURI" {
		t.Errorf("unmarshalNode() go unexpected data")
	}

}

func TestMarshal(t *testing.T) {

	xsdString := NewNamedNode("http://www.w3.org/2001/XMLSchema#string")
	g := NewNamedNode("http://www.test.de/test")
	u1 := NewNamedNode("http://www.test.de/test#User1")
	u2 := NewNamedNode("http://www.test.de/test#User2")
	u3 := NewNamedNode("http://www.test.de/test#Mustermann")
	u4 := NewNamedNode("http://www.test.de/test#Max")
	u5 := NewNamedNode("http://www.test.de/test#Dirk")
	hasLN := NewNamedNode("http://www.test.de/test#hasLastName")
	hasFN := NewNamedNode("http://www.test.de/test#hasFirstName")
	says := NewNamedNode("http://www.test.de/test#says")
	l1 := NewLocalizedLiteral("my tet \"aiaiai", "en")
	t1 := NewTypedLiteral("ui a string value", xsdString)

	stmts := []Statement{
		NewStatement(u1, hasLN, u3, g),
		NewStatement(u1, hasFN, u4, g),
		NewStatement(u1, hasFN, u5, g),
		NewStatement(u1, says, l1, g),
		NewStatement(u1, says, t1, g),
		NewStatement(u2, hasLN, u3, g),
	}

	expect := `@base <http://www.test.de/test> .
<http://www.test.de/test#User1> <http://www.test.de/test#hasFirstName> <http://www.test.de/test#Dirk> , <http://www.test.de/test#Max> ; <http://www.test.de/test#hasLastName> <http://www.test.de/test#Mustermann> ; <http://www.test.de/test#says> "my tet \"aiaiai"@en , "ui a string value"^^<http://www.w3.org/2001/XMLSchema#string> .
<http://www.test.de/test#User2> <http://www.test.de/test#hasLastName> <http://www.test.de/test#Mustermann> .
`

	parser := NewTurtleParser(nil)
	
	ttl, err := parser.Marshal(stmts)
	if err != nil {
		t.Errorf("Marshal() failed")
	}
	if ttl != expect {
		fmt.Printf(expect)
		fmt.Printf(ttl)
		t.Errorf("Marshal() returned unexpected data")
	}


}

func TestMarshalWithOptions(t *testing.T) {

	xsdString := NewNamedNode("http://www.w3.org/2001/XMLSchema#string")
	g := NewNamedNode("http://www.test.de/test")
	u1 := NewNamedNode("http://www.test.de/test#User1")
	u2 := NewNamedNode("http://www.test.de/test#User2")
	u3 := NewNamedNode("http://www.test.de/test#Mustermann")
	u4 := NewNamedNode("http://www.test.de/test#Max")
	u5 := NewNamedNode("http://www.test.de/test-unsub#Dirk")
	hasLN := NewNamedNode("http://www.test.de/test#hasLastName")
	hasFN := NewNamedNode("http://www.test.de/test#hasFirstName")
	says := NewNamedNode("http://www.test.de/test#says")
	l1 := NewLocalizedLiteral("my tet \"aiaiai", "en")
	t1 := NewTypedLiteral("ui a string value", xsdString)

	stmts := []Statement{
		NewStatement(u1, hasLN, u3, g),
		NewStatement(u1, hasFN, u4, g),
		NewStatement(u1, hasFN, u5, g),
		NewStatement(u1, says, l1, g),
		NewStatement(u1, says, t1, g),
		NewStatement(u2, hasLN, u3, g),
	}

	expect := `@base <http://www.test.de/test> .
@prefix : <http://www.test.de/test#> .
@prefix bd: <http://www.bigdata.com/rdf#> .
@prefix bds: <http://www.bigdata.com/rdf/search#> .
@prefix dc: <http://purl.org/dc/elements/1.1#> .
@prefix fn: <http://www.w3.org/2005/xpath-functions#> .
@prefix foaf: <http://xmlns.com/foaf/0.1#> .
@prefix hint: <http://www.bigdata.com/queryHints#> .
@prefix owl: <http://www.w3.org/2002/07/owl#> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix sesame: <http://www.openrdf.org/schema/sesame#> .
@prefix xsd: <http://www.w3.org/2001/XMLSchema#> .

# http://www.test.de/test#User1
:User1 
    :hasFirstName 
        :Max , 
        <http://www.test.de/test-unsub#Dirk> ; 
    :hasLastName 
        :Mustermann ; 
    :says 
        "my tet \"aiaiai"@en , 
        "ui a string value"^^xsd:string .

# http://www.test.de/test#User2
:User2 
    :hasLastName 
        :Mustermann .
`

	p := NewTurtleParser(&TurtleParserOptions{Substitute: true, PrettyPrint: true})
	
	ttl, err := p.Marshal(stmts)
	if err != nil {
		t.Errorf("Marshal() with options failed")
	}
	if ttl != expect {
		fmt.Printf(expect)
		fmt.Printf(ttl)
		t.Errorf("Marshal() returned unexpected data")
	}

}

func TestUnmarshal(t *testing.T) {
	ttl_str := `@base <http://www.test.de/test> .
				<http://www.test.de/test#User1>
					<http://www.test.de/test#hasFirstName> <http://www.test.de/test#Dirk> , <http://www.test.de/test#Max> ;
					<http://www.test.de/test#hasLastName> <http://www.test.de/test#Mustermann> ;
					<http://www.test.de/test#says> "my tet \"aiaiai"@en , "ui a string value"^^<http://www.w3.org/2001/XMLSchema#string> .
				<http://www.test.de/test#User2>
					<http://www.test.de/test#hasLastName> <http://www.test.de/test#Mustermann> ;<http://www.test.de/test#owns> "http://localhost:80/resoures/mine.py"^^<http://www.w3.org/2001/XMLSchema#anyURI> .
				:User1 a :Thing .`

	parser := NewTurtleParser(nil)
	
	stmts, err := parser.Unmarshal(ttl_str)
	if err != nil {
		t.Errorf("Unmarshal() failed: %v", err)
		return
	}
	if len(stmts) != 8 {
		t.Errorf("Unmarshal() returned unexpected number of statements")
	}

}

func TestExtractBaseIri(t *testing.T) {
	valid := `@prefix : <http://www.test.de/test#> .@base <http://www.test.de/test> .
			# http://www.test.de/test#User1
			:User1  :hasFirstName :Max .`
	nobasefull := `@prefix : <http://www.test.de/test#> .
			<http://www.test.de/test#User1>  <http://www.test.de/test#hasFirstName> <http://www.test.de/test#Max> .`
	nobasesub := `@prefix : <http://www.test.de/test#> .
			# http://www.test.de/test#User1
			:User1  :hasFirstName :Max .`
	invalid := `@prefix : <http://www.test.de/test#> .@base <http://www.test.de/test> .
			@base <http://www.test.de/test> .
			# http://www.test.de/test#User1
			:User1  :hasFirstName :Max .`
	ns := NewEmptyNamespace()
	ns.Set("", "http://www.test.de/test")

	p := NewTurtleParser(&TurtleParserOptions{RequireBaseIri: true})
	baseIri, err := p.extractBaseIri(valid, ns, false)
	if err != nil {
		t.Errorf("extractBaseIri() returned with errors")
	}
	if baseIri != "http://www.test.de/test" {
		t.Errorf("extractBaseIri() returned unexpected result")
	}

	baseIri, err = p.extractBaseIri(invalid, ns, false)
	if err == nil {
		t.Errorf("extractBaseIri() returned without errors, while it should fail")
	}

	baseIri, err = p.extractBaseIri(nobasefull, ns, false)
	if err == nil {
		t.Errorf("extractBaseIri() returned without errors, while it should fail")
	}

	baseIri, err = p.extractBaseIri(nobasefull, ns, true)
	if err != nil {
		t.Errorf("extractBaseIri() returned with errors")
	}
	if baseIri != "http://www.test.de/test#User1" {
		t.Errorf("extractBaseIri() returned unexpected data")
	}

	baseIri, err = p.extractBaseIri(nobasesub, ns, true)
	if err != nil {
		t.Errorf("extractBaseIri() returned with errors")
	}
	if baseIri != "http://www.test.de/test#User1" {
		t.Errorf("extractBaseIri() returned unexpected data")
	}

}

func TestExtractNamespace(t *testing.T) {
	valid := `@prefix : <http://www.test.de/test#> . @base <http://www.test.de/test> .
			@prefix abc: <http://www.test.de/abc#> .@prefix def: <http://www.test.de/def#> .
			# http://www.test.de/test#User1
			:User1  :hasFirstName :Max .`

	p := NewTurtleParser(nil)
	ns, err := p.extractNamespace(valid)
	if err != nil {
		t.Errorf("extractNamespace() returned with errors")
	}
	if len(ns.ListKeys()) != 3 {
		t.Errorf("extractNamespace() returned unexpected number of entries")
	}
	if v, ok := ns.Get(""); !ok || v != "http://www.test.de/test" {
		t.Errorf("extractNamespace() failed extracting namespace correctly")
	}
	if v, ok := ns.Get("abc"); !ok || v != "http://www.test.de/abc" {
		t.Errorf("extractNamespace() failed extracting namespace correctly")
	}
	if v, ok := ns.Get("def"); !ok || v != "http://www.test.de/def" {
		t.Errorf("extractNamespace() failed extracting namespace correctly")
	}
}

func TestStripComments(t *testing.T) {
	valid := `@prefix : <http://www.test.de/test#> . @base <http://www.test.de/test> .
			# http://www.test.de/test#User1
			:User1  :hasFirstName :Max .
			# http://www.test.de/test#User1
			:User1  :hasFirstName :Max .`
	expect := `@prefix : <http://www.test.de/test#> . @base <http://www.test.de/test> .
			:User1  :hasFirstName :Max .
			:User1  :hasFirstName :Max .`

	p := NewTurtleParser(nil)
	cleaned := p.stripComments(valid)
	if cleaned != expect {
		t.Errorf("stripComments() returned unexpected data")
	}
}

func TestExtractStatements(t *testing.T) {
	valid := `:User1 :hasFirstName :Max ; :hasLastName :Peter, :Dirk . <http://www.test.de/test#> rdfs:label "ta . da.s\""@de . :User2 :hasFirstName :Other .`
	expect1 := ":User1 :hasFirstName :Max ; :hasLastName :Peter, :Dirk"
	expect2 := `<http://www.test.de/test#> rdfs:label "ta . da.s\""@de`
	expect3 := ":User2 :hasFirstName :Other"

	p := NewTurtleParser(nil)
	statements := p.extractStatements(valid)
	if len(statements) != 3 {
		t.Errorf("extractStatements() returned unexpected number of statements")
	}
	if statements[0] != expect1 || statements[1] != expect2 || statements[2] != expect3 {
		t.Errorf("extractStatements() returned unexpected data in statements")
	}
}

func TestExtractSubStatements(t *testing.T) {
	valid := `:hasFirstName :Max ; :hasLastName :Peter, :Dirk`
	expect1 := ":hasFirstName :Max"
	expect2 := `:hasLastName :Peter, :Dirk`

	p := NewTurtleParser(nil)
	statements := p.extractSubStatements(valid)
	if len(statements) != 2 {
		t.Errorf("extractSubStatements() returned unexpected number of statements")
	}
	if statements[0] != expect1 || statements[1] != expect2 {
		t.Errorf("extractSubStatements() returned unexpected data in statements")
	}
}

func TestExtractObjects(t *testing.T) {
	valid := `<myuri#asds> , "wrds", "", ""@de ,"sdasds"@de,"sdasds"@de, <myuri#asds> , "sdasds"@de ,"sda\"@dasd , sds"@de , :Testx, sdas:Test, "sdasds"@de, "sdasds"^^xsd:string, "sdasds"^^<http://jdsj/asds#sads>,"http://localhost:80/resoures/mine.py"^^<http://www.w3.org/2001/XMLSchema#anyURI>`
	expect := []string{
		`"wrds"`,
		`""`,
		`""@de `,
		`"sdasds"@de`,
		`"sdasds"@de`,
		`"sdasds"@de `,
		`"sda\"@dasd , sds"@de `,
		`"sdasds"@de`,
		`"sdasds"^^xsd:string`,
		`"sdasds"^^<http://jdsj/asds#sads>`,
	 	`"http://localhost:80/resoures/mine.py"^^<http://www.w3.org/2001/XMLSchema#anyURI>`,
		`<myuri#asds> `,
		`<myuri#asds> `,
		`:Testx`,
		`sdas:Test`,
	}

	p := NewTurtleParser(nil)
	objects := p.extractObjects(valid)
	if len(objects) != len(expect) {
		t.Errorf("extractObjects() returned unexpected number of statements")
	}
	for idx, o := range objects {
		if o != expect[idx] {
			t.Errorf("extractObjects() returned unexpected data '%v' but expected '%v'", o, expect[idx])
		}
	}
	
}
