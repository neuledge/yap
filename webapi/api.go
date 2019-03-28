package webapi

import (
	"yap/nlp/format/conll"
	"yap/nlp/types"
)

type Node struct {
	Token    int            `json:"token"`
	Form     string         `json:"form"`
	Lemma    string         `json:"lemma"`
	CPOS     string         `json:"CPOS"`
	POS      string         `json:"POS"`
	Features conll.Features `json:"features"`
	Head     int            `json:"head,omitempty"`
	DepRel   string         `json:"dep,omitempty"`
}

func GraphToNodes(graph types.MorphDependencyGraph) []Node {
	sent := make([]Node, graph.NumberOfNodes())
	arcIndex := make(map[int]types.LabeledDepArc, graph.NumberOfNodes())

	var (
		node   *types.EMorpheme
		arc    types.LabeledDepArc
		headID int
		depRel string
	)

	for _, arcID := range graph.GetEdges() {
		arc = graph.GetLabeledArc(arcID)
		if arc != nil {
			arcIndex[arc.GetModifier()] = arc
		}
	}

	for i, nodeID := range graph.GetVertices() {
		node = graph.GetMorpheme(nodeID)

		if node == nil {
			panic("Can't find node")
		}

		arc, exists := arcIndex[i]
		if exists {
			headID = arc.GetHead()
			depRel = string(arc.GetRelation())
			if depRel == types.ROOT_LABEL {
				headID = -1
			}
		} else {
			headID = -1
			depRel = "None"
		}

		row := Node{
			Token:    node.TokenID - 1,
			Form:     node.Form,
			Lemma:    node.Lemma,
			CPOS:     node.CPOS,
			POS:      node.POS,
			Features: node.Features,
			Head:     headID,
			DepRel:   depRel,
		}

		sent[i] = row
	}

	return sent
}

func TokensToNodes(tokens []types.EMorpheme) []Node {
	sent := make([]Node, len(tokens))

	for i, token := range tokens {
		row := Node{
			Token:    token.TokenID - 1,
			Form:     token.Form,
			Lemma:    token.Lemma,
			CPOS:     token.CPOS,
			POS:      token.POS,
			Features: token.Features,
		}

		sent[i] = row
	}

	return sent
}
