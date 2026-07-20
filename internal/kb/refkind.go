package kb

// Reference kinds carried by Ref.Kind. Content links (RefLink/RefTransclude) and
// typed relations (RefRelates/RefBlocks/RefDependsOn/RefDuplicate) share one
// forward-edge table. Directed kinds read from → to (RefBlocks: from blocks to;
// RefDependsOn: from depends on to). Cross-workspace edges are allowed. Tags are
// NOT refs — they are a separate node↔tag membership (see the node_tags table).
const (
	RefLink       = "link"       // free reference to another node
	RefTransclude = "transclude" // live embed of another node's content
	RefRelates    = "relates"    // undirected association
	RefBlocks     = "blocks"     // from blocks to (directed)
	RefDependsOn  = "depends_on" // from depends on to (directed)
	RefDuplicate  = "duplicate"  // from duplicates to
)

// refKinds is the set of valid Ref.Kind values.
var refKinds = map[string]bool{
	RefLink:       true,
	RefTransclude: true,
	RefRelates:    true,
	RefBlocks:     true,
	RefDependsOn:  true,
	RefDuplicate:  true,
}

// ValidRefKind reports whether k is a known reference kind.
func ValidRefKind(k string) bool {
	return refKinds[k]
}
