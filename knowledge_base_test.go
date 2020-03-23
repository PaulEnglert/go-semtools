package semtools

import (
	"testing"
	"github.com/sirupsen/logrus"
)


func init() {

	logrus.SetLevel(logrus.DebugLevel)

}


func TestNamedNode(t *testing.T) {

	n := NewNamedNode("my-uri")
	if n.Iri() != "my-uri" {
		t.Errorf("NewNamedNode() fails to set proper Iri")
	}

}


func TestLocalizedLiteral(t *testing.T) {

	n := NewLocalizedLiteral("my-text", "lang")
	if n.Value().(string) != "my-text" {
		t.Errorf("NewLocalizedLiteral() fails to set proper Value")
	}
	if n.Language() != "lang" {
		t.Errorf("NewLocalizedLiteral() fails to set proper Language")
	}

	n = NewLocalizedLiteral("my-text", "")
	if n.Value().(string) != "my-text" {
		t.Errorf("NewLocalizedLiteral() fails to set proper Value")
	}
	if n.Language() != "default" {
		t.Errorf("NewLocalizedLiteral() fails to set proper Language")
	}

}


func TestTypedLiteral(t *testing.T) {

	n := NewTypedLiteral("my-text", NewNamedNode("my-type"))
	if n.Value().(string) != "my-text" {
		t.Errorf("NewTypedLiteral() fails to set proper Value")
	}
	if !n.Type().Equals(NewNamedNode("my-type")) {
		t.Errorf("NewTypedLiteral() fails to set proper Type")
	}

}


func TestStatement(t *testing.T) {

	a := NewNamedNode("a")
	b := NewNamedNode("b")
	c := NewNamedNode("c")
	g := NewNamedNode("g")

	s := NewStatement(a, b, c, g)
	if !s.Subject().Equals(a) {
		t.Errorf("NewStatement() fails to set proper Subject")
	}
	if !s.Predicate().Equals(b) {
		t.Errorf("NewStatement() fails to set proper Predicate")
	}
	if !s.Object().Equals(c) {
		t.Errorf("NewStatement() fails to set proper Object")
	}
	if !s.Graph().Equals(g) {
		t.Errorf("NewStatement() fails to set proper Graph")
	}

	s = NewStatement(a, b, c, nil)
	if s.Graph() != nil {
		t.Errorf("NewStatement() fails to set nil Graph")
	}

}


func TestKnowledgeBase(t *testing.T) {

	kb := NewKnowledgeBase("kb")

	stmts := []Statement{
		NewStatement(NewNamedNode("a"), NewNamedNode("b"), NewNamedNode("c"), NewNamedNode("g")),
		NewStatement(NewNamedNode("a"), NewNamedNode("b"), NewNamedNode("d"), NewNamedNode("g")),
		NewStatement(NewNamedNode("a"), NewNamedNode("b"), NewNamedNode("d"), NewNamedNode("g")),
		NewStatement(NewNamedNode("a"), NewNamedNode("b"), NewNamedNode("e"), nil),
	}


	kb.Insert(stmts)
	if len(kb.Statements()) != 3 {
		t.Errorf("Insert() fails to set statements")
	}
	if kb.Statements()[2].Graph() == nil {
		t.Errorf("Insert() fails to set default graph on statements with nil graph")
	}


	kb.Delete([]Statement{stmts[2]})
	if len(kb.Statements()) != 2 {
		t.Errorf("Delete() fails to remove statements")
	}


}
