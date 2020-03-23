package semtools

import (
	"fmt"
	"sort"
)


// Interfaces to work with Knowledge
//
//
//
//

// KnowledgeBase is a named storage that
// manages a collection of statements.
// Management includes additions, removals
// and querying.
type KnowledgeBase interface {

	// Name returns the name of the knowledge base.
	Name() string

	// Statements returns the complete list of statements in
	// the knowledge base.
	Statements() []Statement

	// Insert adds the given statements into the base,
	// ignoring duplicates and already existing ones.
	Insert(stmts []Statement)

	// Delete removes the given statements if they exist.
	Delete(stmts []Statement)

	// Select returns a query operator bound to the knowledge
	// base.
	Select() Query

}

// NewKnowledgeBase will create a new basic knowledge
// base with the given name.
func NewKnowledgeBase(name string) KnowledgeBase {

	if name == "" {
		name = "kb"
	}

	return &knowledgeBase{
		name: name,
		statements: []Statement{},
		defaultGraph: NewNamedNode("default-graph"),
	}

}


// Node is a generic container that makes up a part
// of a statement.
type Node interface {

	// String returns a stringified version of the
	// node. This is not supposed to be used as "value"
	// or something but only for display/logging.
	String() string

	// Equals compares the node to something else and
	// returns if they're equal or not.
	Equals(other interface{}) bool

}

// NamedNode is a generic container of named nodes
// ie. nodes that are referecened by a Iri.
type NamedNode interface {

	Node

	// Iri returns the identifier of the node.
	Iri() string

}

// NewNamedNode creates a new named node with the given
// iri.
func NewNamedNode(iri string) NamedNode {
	return &namedNode{
		iri: iri,
	}
}



// LiteralNode is a generic container for literal
// data.
type LiteralNode interface {

	Node

	// Value returns the value of the literal node.
	Value() interface{}

}


// TypedLiteral is a literal container that has a typed
// value.
type TypedLiteral interface {

	LiteralNode

	// Type returns the type identification of the value.
	Type() NamedNode

}

// NewTypedLiteral creates a new typed literal with the value and type.
func NewTypedLiteral(value interface{}, typeNode NamedNode) TypedLiteral {
	return &typedLiteral{
		value: value,
		typeNode: typeNode,
	}
}


// LocalizedLiteral is a literal container whose value
// is localized, ie. has a language.
type LocalizedLiteral interface {

	LiteralNode

	// Language contains the language of the literals value.
	Language() string

}

// NewLocalizedLiteral creates a new literal from the value and langauge
func NewLocalizedLiteral(value string, language string) LocalizedLiteral {
	if language == "" {
		language = "default"
	}
	return &localizedLiteral{
		value: value,
		language: language,
	}
}


// Statement is the generic container of a statement in the knowledgebase.
// it contains subject, predicate and object information as well as a
// (optional) graph reference.
type Statement interface {

	// Subject returns the subject of the statement.
	Subject() NamedNode

	// Predicate returns the predicate of the statement.
	Predicate() NamedNode

	// Object returns the object of the statement.
	Object() Node

	// Graph returns the graph of the statement, mind that this might be nil.
	Graph() NamedNode

	// Equals compares the statement to something else and
	// returns if they're equal or not.
	Equals(other interface{}) bool

	// String returns a stringified version of the
	// statement. This is not supposed to be used as "value"
	// or something but only for display/logging.
	String() string

}

// NewStatement creates a new statemetn from the given subject, predicate, object and graph
func NewStatement(subject NamedNode, predicate NamedNode, object Node, graph NamedNode) Statement {
	return &statement{
		subject: subject,
		predicate: predicate,
		object: object,
		graph: graph,
	}
}


// Implementations of base KnowledgeBase structures
//
//
//
//



type knowledgeBase struct {
	name string
	statements []Statement
	defaultGraph NamedNode
}

func (kb *knowledgeBase) Name() string {
	return kb.name
}

func (kb *knowledgeBase) Statements() []Statement {
	return kb.statements
}

func (kb *knowledgeBase) Insert(stmts []Statement) {
	for _, stmt := range stmts {

		// make copy of statement and
		// ensure graph is set
		subj := stmt.Subject()
		pred := stmt.Predicate()
		obj := stmt.Object()
		graph := stmt.Graph()
		if graph == nil {
			graph = kb.defaultGraph
		}

		// check if already exists
		matches := kb.Select().
			Graph(graph).
			Subject(subj).
			Predicate(pred).
			Object(obj).
			Results();

		// insert only if no matches
		if len(matches) == 0 {
			kb.statements = append(kb.statements, NewStatement(subj, pred, obj, graph))
		}

	}
}

func (kb *knowledgeBase) Delete(stmts []Statement) {
	for _, stmt := range stmts {

		// query matches
		q := kb.Select()
		if stmt.Graph() != nil {
			q = q.Graph(stmt.Graph())
		}
		matches := q.
			Subject(stmt.Subject()).
			Predicate(stmt.Predicate()).
			Object(stmt.Object()).
			ResultIndexes();

		sort.Sort(sort.Reverse(sort.IntSlice(matches)))

		// insert only if no matches
		for _, idx := range matches {
			kb.statements = append(kb.statements[:idx], kb.statements[idx+1:]...)
		}

	}	
}

func (kb *knowledgeBase) Select() Query {
	return NewQuery().Bind(kb)
}



type namedNode struct {
	iri string
}

func (nn *namedNode) Iri() string {
	return nn.iri
}

func (nn *namedNode) Equals(other interface{}) bool {
	if v, ok := other.(NamedNode); ok {
		return nn.Iri() == v.Iri()
	}
	return false
}

func (nn *namedNode) String() string {
	return nn.Iri()
}



type localizedLiteral struct {
	value string
	language string
}

func (ln *localizedLiteral) Value() interface{} {
	return ln.value
}

func (ln *localizedLiteral) Language() string {
	return ln.language
}

func (ln *localizedLiteral) Equals(other interface{}) bool {
	if v, ok := other.(LocalizedLiteral); ok {
		return ln.Language() == v.Language() && ln.Value() == v.Value()
	}
	return false
}

func (ln *localizedLiteral) String() string {
	return fmt.Sprintf("%v", ln.Value())
}



type typedLiteral struct {
	value interface{}
	typeNode NamedNode
}

func (tn *typedLiteral) Value() interface{} {
	return tn.value
}

func (tn *typedLiteral) Type() NamedNode {
	return tn.typeNode
}

func (tn *typedLiteral) Equals(other interface{}) bool {
	if v, ok := other.(TypedLiteral); ok {
		return tn.Type().Equals(v.Type()) && tn.Value() == v.Value()
	}
	return false
}

func (tn *typedLiteral) String() string {
	return fmt.Sprintf("%v", tn.Value())
}



type statement struct {
	subject NamedNode
	predicate NamedNode
	object Node
	graph NamedNode
}

func (s *statement) Subject() NamedNode {
	return s.subject
}

func (s *statement) Predicate() NamedNode {
	return s.predicate
}

func (s *statement) Object() Node {
	return s.object
}

func (s *statement) Graph() NamedNode {
	return s.graph
}

func (s *statement) Equals(other interface{}) bool {
	if v, ok := other.(Statement); ok {
		if !s.Subject().Equals(v.Subject()) {
			return false
		}
		if !s.Predicate().Equals(v.Predicate()) {
			return false
		}
		if !s.Object().Equals(v.Object()) {
			return false
		}
		if s.Graph() == nil && v.Graph() != nil {
			return false
		}
		if s.Graph() != nil && v.Graph() == nil {
			return false
		}
		if s.Graph() != nil && !s.Graph().Equals(v.Graph()) {
			return false
		}
		return true
	}
	return false
}

func (s *statement) String() string {
	return fmt.Sprintf("%v - %v - %v (g: %v)", s.Subject(), s.Predicate(), s.Object(), s.Graph())
}


