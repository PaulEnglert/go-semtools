package semtools


// Parser is a entity that can parse
// string data to a graph representation
// or the other way around
type Parser interface {

	// Marshal creates a string representation from
	// the provided statements.
	Marshal(stmts []Statement) (string, error)

	// Unmarshal creates a list of statements from the given
	// string data.
	Unmarshal(str string) ([]Statement, error)

}