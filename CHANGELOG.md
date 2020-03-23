# go-semtools Changelog


## [1.0.1] - 2019-09-18
### Changed
- Fix in Equals() of Statement function when both graphs are nil


## [1.0.0] - 2019-09-17
### Refactored
- use KnowledgeBase model, simplified more similar to graph databases, rather then separated graph objects
- add advanced query functionality
- update parser to return statements, rather then full blown graphs


## [0.0.5] - 2019-09-12
### Changed
- fix literal extractor to handle empty literals ala "", or ""@de


## [0.0.4] - 2019-09-03
### Changed
- modified base directive regex to be more narrow

### Added
- more default namespaces
- add Include to Graph


## [0.0.3] - 2019-08-05
### Added
- intitial version of graph model containing Graph, Vertice, Node, NamdNode, LiteralNode, LanguageLiteralNode, TypedLiteralNode
- initial version of parser interface
- initial version of text/ttl parser (from string & to string)
- initial version of namespace structure
- initial version of rudimentary queries: Describe, DescribeMultiple, DescribeAll using NodeDescription as result container