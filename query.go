package semtools

// Query creates a query structure for matching a set
// of statements either by providing the statments
// directly or by binding to a knowledge base.
type Query interface {

	// Group opens a new group. This in combination with
	// Or() allows to create more complex logic in the otherwise
	// rather simple Query interface.
	Group() Query

	// EndGroup closes an opened group.
	EndGroup() Query

	// Or combines the immediately preceeding and immediately
	// following statement with a OR-operator. Using groups
	// before and after combines the result of the groups
	// as a whole.
	Or() Query

	// Graph matches a given graph node in the statements.
	Graph(node NamedNode) Query

	// Subject matches a given subject node in the statements.
	Subject(node NamedNode) Query

	// Predicate matches a given predicate node in the statements.
	Predicate(node NamedNode) Query

	// Object matches a given object node in the statements.
	Object(node Node) Query

	// Bind binds a knowledge base to the query and will allow
	// the use of Result() and ResultIndexes()
	Bind(base KnowledgeBase) Query

	// Evaluate executes the query on the given statement
	// and returns true/false depending on if it was matched.
	Evaluate(stmt Statement) bool

	// Results executes the query on the bound knowledge base
	// returning the matched statements.
	Results() []Statement

	// ResultsFrom works like Results() just that it takes
	// the set of statements as a parameter rather then using
	// the bound knowledge base.
	ResultsFrom(stmts []Statement) []Statement

	// ResultIndexes executes the query on the bound knowledge base
	// returning the indexes of the matched statements.
	ResultIndexes() []int

	// ResultIndexesFrom works like ResultIndexes() just that it takes
	// the set of statements as a parameter rather then using
	// the bound knowledge base.
	ResultIndexesFrom(stmts []Statement) []int

}

// NewQuery creates a new query object.
func NewQuery() Query {
	return &query{
		base: nil,
		parent: nil,
		query: []matcher{func(stmt Statement) bool {return true}},
		nextOp: "and",
	}
}


// NewQueryWithParent creates a new query object using the parent.
func NewQueryWithParent(parent Query) Query {
	return &query{
		base: nil,
		parent: parent,
		query: []matcher{func(stmt Statement) bool {return true}},
		nextOp: "and",
	}
}



type matcher = func (stmt Statement) bool

type query struct {
	base KnowledgeBase
	parent Query
	query []matcher
	nextOp string
}


func (q *query) Group() Query {
	sq := NewQueryWithParent(q)
	q.add(func(stmt Statement) bool {
		return sq.Evaluate(stmt)
	})
	return sq
}

func (q *query) EndGroup() Query {
    if q.parent == nil {
        return q
    }
    return q.parent;
}

func (q *query) Or() Query {
	q.nextOp = "or"
	return q
}

func (q *query) Graph(node NamedNode) Query {
	q.add(func(stmt Statement) bool {
		if stmt.Graph() == nil {
			return node == nil
		}
		return stmt.Graph().Equals(node)  
	})
	return q
}

func (q *query) Subject(node NamedNode) Query {
	q.add(func(stmt Statement) bool {
		return stmt.Subject().Equals(node)  
	})
	return q
}

func (q *query) Predicate(node NamedNode) Query {
	q.add(func(stmt Statement) bool {
		return stmt.Predicate().Equals(node)  
	})
	return q
}

func (q *query) Object(node Node) Query {
	q.add(func(stmt Statement) bool {
		return stmt.Object().Equals(node)  
	})
	return q
}

func (q *query) Bind(base KnowledgeBase) Query {
	q.base = base
	return q
}

func (q *query) Evaluate(stmt Statement) bool {
	for _, m := range q.query {
		if !m(stmt) {
			return false
		}
	}
	return true
}

func (q *query) Results() []Statement {
	stmts := []Statement{}
	if q.base != nil {
		stmts = q.base.Statements()
	}
	return q.ResultsFrom(stmts)
}

func (q *query) ResultsFrom(stmts []Statement) []Statement {
	res := []Statement{}
	for _, stmt := range stmts {
		if q.Evaluate(stmt) {
			res = append(res, stmt)
		}
	}
	return res
}

func (q *query) ResultIndexes() []int {
	stmts := []Statement{}
	if q.base != nil {
		stmts = q.base.Statements()
	}
	return q.ResultIndexesFrom(stmts)
}

func (q *query) ResultIndexesFrom(stmts []Statement) []int {
	res := []int{}
	for idx, stmt := range stmts {
		if q.Evaluate(stmt) {
			res = append(res, idx)
		}
	}
	return res
}

func (q *query) add(qmatcher matcher) {
	op := q.nextOp
	q.nextOp = "and"

	if op == "and" {

		q.query = append(q.query, qmatcher)

	} else if op == "or" {

		lidx := len(q.query) - 1
		prev := q.query[lidx]
		q.query[lidx] = func(stmt Statement) bool {
			return prev(stmt) || qmatcher(stmt)
		}

	}

}
