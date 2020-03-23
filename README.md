# go-semtools

[![go-mod-version](https://img.shields.io/github/go-mod/go-version/PaulEnglert/go-semtools)](https://github.com/PaulEnglert/go-semtools)
[![tag](https://img.shields.io/github/v/tag/PaulEnglert/go-semtools)](https://github.com/PaulEnglert/go-semtools)


A toolset for working with semantic data, such as ontologies, triples and general graph data.

## Knowledge Base

The semtools define a graph scheme that consists of `Nodes` and `Statements`. Nodes can be named (ie. referenced via an Iri) or unnamed and contain data (ie. literal nodes). A Statement describes the connection between node, using a subject, predicate, object-syntax. All of these elements together are placed within a `Graph` context that provides "grouping" of statements.

Usage

    // create a new knowledge base   
    kb := NewKnowledgeBase("mybase")

    // add some statements to the base
    graph := NewNamedNode("http://mygraph")
    kb.Insert([]Statement{
        NewStatemet(
            NewNamedNode(graph.Iri() + "#Max"), NewNamedNode(graph.Iri() + "#knows"),
            NewNamedNode(graph.Iri() + "#Mara"), graph),
        NewStatemet(
            NewNamedNode(graph.Iri() + "#Mara"), NewNamedNode(graph.Iri() + "#knows"),
            NewNamedNode(graph.Iri() + "#Max"), graph),
        NewStatemet(
            NewNamedNode(graph.Iri() + "#Mara"), NewNamedNode(graph.Iri() + "#says"),
            NewLocalizedLiteral("Hi Max", "en"), graph),
    })

    // now we could query that graph for some information
    statementsAboutMara := kb.Select().
        Graph(graph). // optionally filter to specific graph
        Subject(NewNamedNode(graph.Iri() + "#Mara")).
        Results()

    for _, s := range statementsAboutMara {
        fmt.Printf("Mara %v %v\n", s.Predicate(), s.Object())
    }

## Querying

A knowledge base can either be manually worked with using the `Statements()`, or one can use the `Select()` or custom `Query` objects to work on the underlaying data. For details and examples see [Query](./query.go).

## Parsing

Knowledge Base content can be complex and need to be communicated to and from other sources. The parser component permit doing exactly that, the available parsers are currently:

* [Turtle](https://en.wikipedia.org/wiki/Turtle_(syntax)): `TurtleParser`

