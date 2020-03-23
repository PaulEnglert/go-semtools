package semtools


import (
	"testing"
	"github.com/sirupsen/logrus"
)


func init() {

	logrus.SetLevel(logrus.DebugLevel)

}


func TestQuery(t *testing.T) {

	kb := NewKnowledgeBase("kb")
	stmts := []Statement{
		NewStatement(NewNamedNode("max"), NewNamedNode("knows"), NewNamedNode("mara"), NewNamedNode("friends")),
		NewStatement(NewNamedNode("mara"), NewNamedNode("knows"), NewNamedNode("max"), NewNamedNode("friends")),
		NewStatement(NewNamedNode("bill"), NewNamedNode("knows"), NewNamedNode("max"), NewNamedNode("friends")),
		NewStatement(NewNamedNode("mara"), NewNamedNode("owns"), NewNamedNode("bill"), NewNamedNode("friends")),
	}
	kb.Insert(stmts)

	res := kb.Select().
		Graph(NewNamedNode("something")).Results()
	if len(res) != 0 {
		t.Errorf("Expected results to contain no statements")
	}

	res = kb.Select().
		Graph(NewNamedNode("friends")).Results()
	if len(res) != 4 {
		t.Errorf("Expected results to contain 4 statements")
	}

	res = kb.Select().
		Subject(NewNamedNode("max")).Results()
	if len(res) != 1 {
		t.Errorf("Expected results to contain 1 statements")
	}

	res = kb.Select().
		Object(NewNamedNode("max")).Results()
	if len(res) != 2 {
		t.Errorf("Expected results to contain 2 statements")
	}

	res = kb.Select().
		Predicate(NewNamedNode("knows")).Results()
	if len(res) != 3 {
		t.Errorf("Expected results to contain 3 statements")
	}

	res = kb.Select().
		Subject(NewNamedNode("max")).
		Or().
		Object(NewNamedNode("max")).Results()
	if len(res) != 3 {
		t.Errorf("Expected results to contain 3 statements")
	}

	res = kb.Select().
		Group().
			Subject(NewNamedNode("bill")).
			Object(NewNamedNode("max")).
		EndGroup().
		Or().
		Group().
			Subject(NewNamedNode("max")).
			Object(NewNamedNode("mara")).
		EndGroup().Results()
	if len(res) != 2 {
		t.Errorf("Expected results to contain 2 statements")
	}

}