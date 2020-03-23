package semtools

// Namespace defines abbreviations for long uris and
// names within a graph context.
type Namespace struct {

	// ns contains the mapping between abbreviation
	// and value
	ns map[string]string
}

// NewNamespace creates a new namespace including the default
// values globally provided.
func NewNamespace() *Namespace {
	ns := NewEmptyNamespace()
	ns.Set("xsd", "http://www.w3.org/2001/XMLSchema")
	ns.Set("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns")
	ns.Set("rdfs", "http://www.w3.org/2000/01/rdf-schema")
	ns.Set("owl", "http://www.w3.org/2002/07/owl")
	ns.Set("sesame", "http://www.openrdf.org/schema/sesame")
	ns.Set("fn", "http://www.w3.org/2005/xpath-functions")
	ns.Set("foaf", "http://xmlns.com/foaf/0.1")
	ns.Set("dc", "http://purl.org/dc/elements/1.1")
	ns.Set("hint", "http://www.bigdata.com/queryHints")
	ns.Set("bd", "http://www.bigdata.com/rdf")
	ns.Set("bds", "http://www.bigdata.com/rdf/search")
	return ns
}

// NewEmptyNamespace creates a new namespace without any
// names.
func NewEmptyNamespace() *Namespace {
	return &Namespace{ns: map[string]string{}}
}

// ListKeys provides a list of the keys (abbreviations) defined in
// the namespace.
func (ns *Namespace) ListKeys() []string {
	keys := make([]string, len(ns.ns))
	idx := 0
	for k, _ := range ns.ns {
		keys[idx] = k
		idx += 1
	}
	return keys
}

// ListValues provides a list of the values (long versions) defined
// in the namespace
func (ns *Namespace) ListValues() []string {
	values := make([]string, len(ns.ns))
	idx := 0
	for _, v := range ns.ns {
		values[idx] = v
		idx += 1
	}
	return values
}

// Contains checks if a key exists within the namespace.
func (ns *Namespace) Contains(key string) bool {
	_, ok := ns.ns[key]
	return ok
}

// Get returns the value of the key. Check the ok flag
// in order to make sure that the key was also present.
func (ns *Namespace) Get(key string) (string, bool) {
	v, ok := ns.ns[key]
	return v, ok
}

// MustGet works like Get, just assumes the Key exists.
func (ns *Namespace) MustGet(key string) string {
	v, _ := ns.Get(key)
	return v
}

// Set adds/overwrites an abbreviation with a value.
func (ns *Namespace) Set(key string, value string) {
	ns.ns[key] = value
}

// Include incorporates another namespace into itself.
func (ns *Namespace) Include(other *Namespace) *Namespace {
	copy := ns.Copy()
	for _, k := range other.ListKeys() {
		copy.Set(k, other.MustGet(k))
	}
	return copy
}

// Copy creates a copy of the namespace.
func (ns *Namespace) Copy() *Namespace {
	copy := NewNamespace()
	for _, k := range ns.ListKeys() {
		copy.Set(k, ns.MustGet(k))
	}
	return copy
}
