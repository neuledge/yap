package types

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"yap/alg"
	"yap/alg/graph"
	"yap/util"
)

type Morpheme struct {
	graph.BasicDirectedEdge
	Form       string
	Lemma      string
	CPOS       string
	POS        string
	Features   map[string]string
	TokenID    int
	FeatureStr string
}

type EMorpheme struct {
	Morpheme
	EForm, EFCPOS, EPOS int
	EFeatures           int
	EMHost, EMSuffix    int
}

var _ DepNode = &Morpheme{}
var _ DepNode = &EMorpheme{}

func NewRootMorpheme() *EMorpheme {
	return &EMorpheme{Morpheme: Morpheme{
		graph.BasicDirectedEdge{0, 0, 0},
		ROOT_TOKEN, ROOT_TOKEN, ROOT_TOKEN, ROOT_TOKEN,
		nil, 0, "",
	}}
}

func (m *Morpheme) ID() int {
	return m.BasicDirectedEdge.ID()
}

func (m *Morpheme) From() int {
	return m.BasicDirectedEdge.From()
}

func (m *Morpheme) To() int {
	return m.BasicDirectedEdge.To()
}

func (m *Morpheme) String() string {
	return fmt.Sprintf("%v-%v-%v-%s", m.Form, m.CPOS, m.POS, m.FeatureStr)
}

func (m *Morpheme) Copy() *Morpheme {
	newMorph := new(Morpheme)
	*newMorph = *m
	newMorph.Features = make(map[string]string)
	for k, v := range m.Features {
		newMorph.Features[k] = v
	}
	return newMorph
}

func (m *Morpheme) EMorpheme() *EMorpheme {
	newMorph := new(Morpheme)
	*newMorph = *m
	newMorph.Features = make(map[string]string)
	for k, v := range m.Features {
		newMorph.Features[k] = v
	}
	return &EMorpheme{Morpheme: *newMorph}
}

func (m *Morpheme) Equal(otherEq util.Equaler) bool {
	other := otherEq.(*Morpheme)
	featEq := (len(m.Features) == 0 && len(other.Features) == 0) || reflect.DeepEqual(m.Features, other.Features)
	return m.Form == other.Form &&
		m.CPOS == other.CPOS &&
		m.POS == other.POS &&
		featEq
}

func (m *EMorpheme) Equal(otherEq util.Equaler) bool {
	other := otherEq.(*EMorpheme)
	return m.EForm == other.EForm &&
		m.EPOS == other.EPOS &&
		m.EFCPOS == other.EFCPOS &&
		m.EFeatures == other.EFeatures
}

func (m *EMorpheme) Copy() *EMorpheme {
	newMorph := new(EMorpheme)
	*newMorph = *m
	newMorph.Features = make(map[string]string)
	for k, v := range m.Features {
		newMorph.Features[k] = v
	}
	return newMorph
}

var _ graph.DirectedEdge = &Morpheme{}
var _ graph.DirectedEdge = &EMorpheme{}

type Morphemes []*EMorpheme
type BasicMorphemes []*Morpheme

var _ alg.Index = make(Morphemes, 1)

func (m Morphemes) Len() int {
	return len(m)
}

func (m Morphemes) Less(i, j int) bool {
	return m[i].From() < m[j].From() ||
		(m[i].From() == m[j].From() && m[i].To() < m[j].To())
}

func (m Morphemes) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Morphemes) Index(index int) (int, bool) {
	if index >= len(m) {
		return 0, false
	}
	return len(m) - 1 - index, true
}

func (m *BasicMorphemes) Union(others BasicMorphemes) {
	if len(others) != 1 {
		panic("Can't Union with another morpheme set with size other than 1")
	}
	other := others[0]
	for _, cur := range *m {
		if cur.Equal(other) {
			return
		}
	}
	other.BasicDirectedEdge[0] = len(*m)
	*m = append(*m, other)
}

func (m Morphemes) Standalone() BasicMorphemes {
	// TODO: think of a better name - should mean 'retrieve the
	// raw morphemes, as if they appear by themselves'
	if len(m) != 1 {
		panic("Can't return standalone for morpheme set with size other than 1")
	}
	newMorph := new(Morpheme)
	*newMorph = m[0].Morpheme
	newMorph.BasicDirectedEdge = [3]int{0, 0, 1}
	newMorph.TokenID = 0
	return BasicMorphemes{newMorph}
}

type Spellout Morphemes

func (s Spellout) Compare(other Spellout, paramFuncName string) (TP, TN, FP, FN int) {
	// log.Println("Comparing", s.AsString(), other.AsString())
	// if s.Equal(other) {
	// 	log.Println("Are Equal")
	// }
	paramFunc, exists := MDParams[paramFuncName]
	if !exists {
		panic("Unsupported parameter function: " + paramFuncName)
	}
	curMorphs, otherMorphs := make(map[string]bool, len(s)), make(map[string]bool, len(other))
	for _, m := range s {
		curMorphs[paramFunc(m)] = true
	}
	for _, m := range other {
		otherMorphs[paramFunc(m)] = true
	}
	for k := range curMorphs {
		if _, exists := otherMorphs[k]; exists {
			TP += 1
		} else {
			FP += 1
		}
	}
	for k := range otherMorphs {
		if _, exists := curMorphs[k]; !exists {
			TN += 1
		}
	}
	// log.Println("Results", TP, TN, FP, FN)
	return
}

type Spellouts []Spellout

func (s Spellouts) Len() int {
	return len(s)
}

func (s Spellouts) Less(i, j int) bool {
	return s[i].AsString() < s[j].AsString()
}

func (s Spellouts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Spellout) String() string {
	posStrings := make([]string, len(s))
	for i, morph := range s {
		posStrings[i] = morph.CPOS
	}
	return strings.Join(posStrings, ":")
}

func (s Spellout) AsString() string {
	strs := make([]string, len(s))
	for i, morph := range s {
		strs[i] = morph.String()
	}
	return strings.Join(strs, ";")
}

func (s Spellout) Equal(other Spellout) bool {
	if len(s) != len(other) {
		return false
	}
	for i, val := range other {
		if !s[i].Equal(val) {
			return false
		}
	}
	return true
}

func (s Spellouts) Find(other Spellout) (int, bool) {
	for i, cur := range s {
		if cur.Equal(other) {
			return i, true
		}
	}
	return 0, false
}

type Mapping struct {
	Token    Token
	Spellout Spellout
}

func (m *Mapping) Equal(other *Mapping) bool {
	return m.Token == other.Token && m.Spellout.Equal(other.Spellout)
}

func (m *Mapping) String() string {
	if len(m.Spellout) > 0 {
		return fmt.Sprintf("%v|%v-%v-%v", m.Token, m.Spellout[len(m.Spellout)-1].Form, m.Spellout[len(m.Spellout)-1].POS, m.Spellout[len(m.Spellout)-1].FeatureStr)
	} else {
		return string(m.Token)
	}
}

type Mappings []*Mapping

var _ alg.Index = make(Mappings, 1)

func (ms Mappings) Equal(otherEq util.Equaler) bool {
	other, ok := otherEq.(Mappings)
	if !ok {
		return false
	}
	if len(ms) != len(other) {
		return false
	}
	for i, m := range ms {
		if !m.Equal(other[i]) {
			return false
		}
	}
	return true
}

func (ms Mappings) Index(i int) (int, bool) {
	if i >= len(ms) {
		return 0, false
	}
	return i, true
}

type Path int

type Lattice struct {
	Token           Token
	Morphemes       Morphemes
	Spellouts       Spellouts
	Next            map[int][]int
	BottomId, TopId int
}

func (l *Lattice) Add(morphs BasicMorphemes, start, end, numToken int) {
	nextNode := start
	for i, morph := range morphs {
		newMorph := morph.EMorpheme()
		newMorph.TokenID = numToken
		newMorph.BasicDirectedEdge[0] = len(l.Morphemes)
		// log.Println("\t\t\t\tSetting first node", nextNode)
		newMorph.BasicDirectedEdge[1] = nextNode
		if i < len(morphs)-1 {
			// log.Println("\t\t\t\tSearch for outgoing node")
			exists := true
			for exists {
				// log.Println("\t\t\t\t\tFound outgoing node", nextNode)
				// log.Println("\t\t\t\t\tIn", l.Next)
				nextNode++
				_, exists = l.Next[nextNode]
			}
			// log.Println("\t\t\t\tSetting outgoing node", nextNode)
			newMorph.BasicDirectedEdge[2] = nextNode
		} else {
			// log.Println("\t\t\t\tSetting last node", end)
			newMorph.BasicDirectedEdge[2] = end
		}
		// log.Println("\t\t\tadding morph", i, morph, "at nodes", newMorph.From(), newMorph.To())
		if _, exists := l.Next[newMorph.From()]; exists {
			// log.Println("\t\t\tappending morph ID", newMorph.ID(), "to", l.Next[newMorph.From()])
			l.Next[newMorph.From()] = append(l.Next[newMorph.From()], newMorph.ID())
		} else {
			l.Next[newMorph.From()] = []int{newMorph.ID()}
			// log.Println("\t\t\tcreating new morph next list for", newMorph.ID(), "at", newMorph.From(), ":", l.Next[newMorph.From()])
		}
		l.Morphemes = append(l.Morphemes, newMorph)
	}
}

func (l *Lattice) BumpTop(from, to int, upTo int) {
	for _, morph := range l.Morphemes {
		if morph.ID() < upTo && morph.To() == from {
			// log.Println("\t\t\t\t\tBumping top for", morph)
			morph.BasicDirectedEdge[2] = to
		}
	}
	l.TopId = to
}
func (l *Lattice) AddAnalysis(prefix BasicMorphemes, hosts []BasicMorphemes, numToken int) {
	// log.Println("\t\t\tStarting with top", prevTop)
	startNode := l.BottomId
	oldestId := len(l.Morphemes)
	if prefix != nil {
		maxSameMorphNode := l.BottomId
		lastMatchingMorph := -1
		for i, m := range prefix {
			// log.Println("\t\t\t\tSearching for morpheme", m)
			edges, _ := l.Next[maxSameMorphNode]
			for _, edgeId := range edges {
				edge := l.Morphemes[edgeId]
				if edge.Morpheme.Equal(m) {
					maxSameMorphNode = edge.To()
					lastMatchingMorph = i
					break
				}
			}
			if lastMatchingMorph < i {
				break
			}
		}
		if lastMatchingMorph < len(prefix)-1 {
			prefixTail := prefix[lastMatchingMorph+1:]
			// log.Println("\t\t\tAdding rest of prefix:", prefixTail)
			addTopOffset := 0
			if val, exists := l.Next[maxSameMorphNode]; exists && len(val) > 0 {
				// log.Println("\t\t\tmaxSameMorphNode", maxSameMorphNode, "exists")
				addTopOffset = -1
			}
			endOfPrefix := l.TopId + len(prefixTail) + addTopOffset
			l.Add(prefixTail, maxSameMorphNode, endOfPrefix, numToken)
			// log.Println("\t\t\tBump Top to:", l.TopId+len(prefixTail))
			l.BumpTop(l.TopId, l.TopId+len(prefixTail), oldestId)
			// log.Println("\t\t\tSetting maxSameMorphNode:", endOfPrefix)
			maxSameMorphNode = endOfPrefix
		}
		startNode = maxSameMorphNode
	}
	// log.Println("\t\tadding host")
	for _, host := range hosts {
		// log.Println("\t\t\tAdding Host:", host, "at", startNode)
		newTop := len(host) + startNode
		l.Add(host, startNode, newTop, numToken)

		if newTop >= l.TopId {
			// log.Println("\t\t\tBump Top to:", newTop, "from:", l.TopId)
			l.BumpTop(l.TopId, newTop, oldestId)
		}
		oldestId = len(l.Morphemes)
	}
	// optionally regenerate spellout
}

func (l *Lattice) Optimize() {
	// removed := make(map[int]bool, len(l.Next))
	// for node, out := range l.Next {
	// 	if _, exists := removed[node]; !exists {
	// 		toRemove := make(map[int]int, len(out))
	// 		for i, outId1 := range out[:len(out)-2] {
	// 			if _, id1Exists := toRemove[outId1]; !id1Exists {
	// 				for _, outId2 := range out[i+1:] {
	// 					if _, id2Exists := toRemove[outId2]; !id2Exists {
	// 						m1, m2 := l.Morphemes[outId1], l.Morphemes[outId2]
	// 						if m1.Equal(m2) {
	// 							toRemove[outId2] = outId1
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func (l *Lattice) BridgeMissingMorphemes() {
	for _, m := range l.Morphemes {
		if _, exists := l.Next[m.To()]; !exists && m.To() < l.TopId {
			if _, nextExists := l.Next[m.To()+1]; nextExists {
				// log.Println("Bridging morpheme", m.Form, "from", m.To(), "to", m.To()+1)
				m.BasicDirectedEdge[2] += 1
			} else {
				// log.Println("Morpheme's next does not exist and cannot bridge! (", m.Form, m.From(), m.To(), ")")
			}
		}
	}
}

func (l *Lattice) UnionPath(other *Lattice) {
	// assume other is a "gold" path (only one "next" at each node)
	// add gold lattice path if it is an alternative to existing paths with the
	// same nodes
	formMorphs := make(map[string][]*EMorpheme)
	for _, predMorph := range l.Morphemes {
		// log.Println("Found morpheme", predMorph, "at", predMorph.From(), predMorph.To())
		if cur, exists := formMorphs[predMorph.Form]; exists {
			formMorphs[predMorph.Form] = append(cur, predMorph)
		} else {
			formMorphs[predMorph.Form] = []*EMorpheme{predMorph}
		}
	}
	var found, missingMorpheme bool

	for _, goldMorph := range other.Morphemes {
		// log.Println("At morph", goldMorph)
		if curMorphs, exists := formMorphs[goldMorph.Form]; exists {
			for _, curMorph := range curMorphs {
				// log.Println("\tComparing to morph", curMorph)
				if curMorph.Equal(goldMorph) {
					found = true
				}
			}
		} else {
			// log.Println("Warning: gold morph form", goldMorph.Form, "is not in pred lattice!")
			missingMorpheme = true
			continue
		}
		if !found {
			// log.Println("Getting example for", goldMorph.Form)
			exampleMorphs, _ := formMorphs[goldMorph.Form]
			// log.Println("Examples", exampleMorphs)
			exampleFromTos := make(map[string][]int)
			for _, example := range exampleMorphs {
				exFromToStr := fmt.Sprintf("%v-%v", example.From(), example.To())
				if _, exists := exampleFromTos[exFromToStr]; !exists {
					// log.Println("Found new pair of from to for example", exFromToStr)
					exampleFromTos[exFromToStr] = []int{example.From(), example.To()}
				}
			}
			for _, exampleFromTo := range exampleFromTos {
				// log.Println("Found example pair at", exampleFromTo[0], exampleFromTo[1])
				// log.Println("Adding missing morpheme (form with same POS/properties did not exist)", goldMorph.Form, goldMorph.POS, goldMorph.CPOS, goldMorph.FeatureStr)
				l.InfuseMorph(goldMorph, exampleFromTo[0], exampleFromTo[1], true)
			}
		}
		found = false
	}
	// a gold form was not found in the pred lattice, try to find a trivial attachment
	// fast forward equivalent epilogue morphemes, try to equate missing morpheme to
	// the fusion of the rest of the pred lattice and attach
	if missingMorpheme {
		// the gold lattice should have only one spellout
		var (
			prevPredNodeId    int = -1
			currentPredNodeId int
			currentNode       *EMorpheme
			currentNodes      []int
			nextExists        bool
		)
		currentPredNodeId = l.Bottom()
	GoldLoop:
		for i, goldMorph := range other.Spellouts[0] {
			currentNodes, nextExists = l.Next[currentPredNodeId]
			if !nextExists {
				panic("Lost in pred lattice")
			}
			for _, currentNodeId := range currentNodes {
				if currentNode = l.Morphemes[currentNodeId]; currentNode.Equal(goldMorph) {
					// gold morpheme was found, move on to the next gold morpheme
					prevPredNodeId = currentPredNodeId
					currentPredNodeId = currentNode.BasicDirectedEdge[2]
					// log.Println("Found morpheme at", goldMorph.Form)
					continue GoldLoop
				}
			}
			// log.Println("Failed to find morpheme at", goldMorph.Form)
			// if the previous inner loop did not "continue" the goldloop
			// we found the location of the missing gold morpheme
			// we try to match with the fused morpheme from this point on
			for _, fusedCandidate := range l.AllFusedFrom(currentPredNodeId) {
				if fusedCandidate == goldMorph.Form {
					// if successful, we set the start node to the current node and the end node to the
					// top of the lattice
					// log.Println("Adding missing morpheme (form did not exist)", goldMorph.Form, goldMorph.POS, goldMorph.CPOS, goldMorph.FeatureStr, ";", currentPredNodeId)
					l.InfuseMorph(goldMorph, currentPredNodeId, l.Top(), true)
					break GoldLoop
				}
			}
			log.Println("Failed to find at current morpheme, trying previous", goldMorph.Form)
			// failed to fuse from current node, try to backtrack
			// maybe previous node will succeed
			if prevPredNodeId > -1 {
				for _, fusedCandidate := range l.AllFusedFrom(prevPredNodeId) {
					if fusedCandidate == other.Spellouts[0][i-1].Form {
						// log.Println("Adding missing morpheme (form did not exist); at", goldMorph.Form, goldMorph.POS, goldMorph.CPOS, goldMorph.FeatureStr, ";", currentPredNodeId)
						l.InfuseMorph(goldMorph, currentPredNodeId, l.Top(), true)
					}
				}
			}
		}
		if len(formMorphs) == 1 {
			for formMorph, edges := range formMorphs {
				// will happen exactly once
				// test to see if pred is a concatenation of gold morphs
				// if yes, and a morpheme is missing, assume morphological
				// analysis failure and complete gold path is missing
				goldMorphs := make([]string, len(other.Morphemes))
				for i, goldMorph := range other.Morphemes {
					goldMorphs[i] = goldMorph.Form
				}
				goldConcat := strings.Join(goldMorphs, "")
				if formMorph == goldConcat {
					log.Println("Found morphological analysis failure, adding gold path", strings.Join(goldMorphs, "-"))
					nextFromLatNode := edges[0].From()
					nextToLatNode := edges[0].To() + 1
					for i, goldMorph := range other.Morphemes {
						// log.Println("Need to fuse", goldMorph, "maybe at", nextFromLatNode, nextToLatNode)
						if i == len(other.Morphemes)-1 {
							nextToLatNode = edges[0].To()
							// log.Println("Updating", nextToLatNode)
						}
						l.InfuseMorph(goldMorph, nextFromLatNode, nextToLatNode, false)
						nextFromLatNode = nextToLatNode
						nextToLatNode++
					}
					l.Spellouts = nil
					l.GenSpellouts()
				}
			}
		}
	}
}

func (l *Lattice) InfuseMorph(morph *EMorpheme, from, to int, genSpellout bool) {
	// log.Println("Infusing", morph, "at", from, to)
	newMorph := morph.Copy()
	newMorph.Morpheme.BasicDirectedEdge[1] = from
	newMorph.Morpheme.BasicDirectedEdge[2] = to
	id := len(l.Morphemes)
	newMorph.Morpheme.BasicDirectedEdge[0] = id
	l.Morphemes = append(l.Morphemes, newMorph)
	mList, _ := l.Next[newMorph.From()]
	l.Next[newMorph.From()] = append(mList, id)

	if genSpellout {
		l.Spellouts = nil
		l.GenSpellouts()
	}
}

func NewRootLattice() Lattice {
	morphs := make(Morphemes, 1)
	morphs[0] = NewRootMorpheme()
	lat := &Lattice{
		ROOT_TOKEN,
		morphs,
		nil,
		make(map[int][]int),
		0,
		0,
	}
	return *lat
}

type LatticeSentence []Lattice

var _ Sentence = LatticeSentence{}

func (ls LatticeSentence) Tokens() []string {
	res := make([]string, len(ls))
	for i, val := range ls {
		res[i] = string(val.Token)
	}
	return res
}

func (ls LatticeSentence) Equal(otherEq util.Equaler) bool {
	otherSent := otherEq.(Sentence)
	if len(otherSent.Tokens()) != len(ls) {
		return false
	}
	otherToks := otherSent.Tokens()
	curToks := ls.Tokens()
	return reflect.DeepEqual(curToks, otherToks)
}

func (l *Lattice) GetDirectedEdge(i int) graph.DirectedEdge {
	return graph.DirectedEdge(l.Morphemes[i])
}

func (l *Lattice) GetEdge(i int) graph.Edge {
	return graph.Edge(l.Morphemes[i])
}

func (l *Lattice) GetEdges() []int {
	res := make([]int, len(l.Morphemes))
	for i, _ := range l.Morphemes {
		res[i] = i
	}
	return res
}

func (l *Lattice) GetVertices() []int {
	vSet := make(map[int]bool)
	for _, edge := range l.Morphemes {
		vSet[edge.From()] = true
		vSet[edge.To()] = true
	}
	res := make([]int, 0, len(vSet))
	for k, _ := range vSet {
		res = append(res, k)
	}
	return res
}

func (l *Lattice) GetVertex(i int) graph.Vertex {
	return graph.BasicVertex(i)
}

func (l *Lattice) NumberOfEdges() int {
	return len(l.Morphemes)
}

func (l *Lattice) NumberOfVertices() int {
	return l.Top() - l.Bottom()
}

var _ graph.BoundedLattice = &Lattice{}
var _ graph.DirectedGraph = &Lattice{}

// untested..
func (l *Lattice) Inf(i, j int) int {
	iReachable := make(map[int]int)
	for path := range graph.YieldAllPaths(graph.DirectedGraph(l), l.Bottom(), i) {
		for i, el := range path {
			dist := len(path) - i - 1
			iReachable[el.ID()] = dist
		}
	}
	var bestVal, bestDist int = -1, -1
	for path := range graph.YieldAllPaths(graph.DirectedGraph(l), l.Bottom(), j) {
		for i, _ := range path {
			el := path[len(path)-i-1]
			dist, exists := iReachable[el.ID()]
			if exists {
				if bestDist == -1 || bestDist > dist {
					bestVal = el.ID()
					bestDist = dist
					break
				}
			}
		}
	}
	return bestVal
}

// untested..
func (l *Lattice) Sup(i, j int) int {
	iReachable := make(map[int]int)
	for path := range graph.YieldAllPaths(graph.DirectedGraph(l), i, l.Top()) {
		for dist, el := range path {
			iReachable[el.ID()] = dist
		}
	}
	var bestVal, bestDist int = -1, -1
	for path := range graph.YieldAllPaths(graph.DirectedGraph(l), j, l.Top()) {
		for _, el := range path {
			dist, exists := iReachable[el.ID()]
			if exists {
				if bestDist == -1 || bestDist > dist {
					bestVal = el.ID()
					bestDist = dist
					break
				}
			}
		}
	}
	return bestVal
}

func (l *Lattice) Top() int {
	return l.TopId
}

func (l *Lattice) Bottom() int {
	return l.BottomId
}

func (l *Lattice) MaxPathLen() int {
	if len(l.Morphemes) == 0 {
		return 0
	}
	return l.Top() - l.Bottom()
}

func (l *Lattice) SortMorphemes() {
	sort.Sort(l.Morphemes)
}

func (l *Lattice) SortNexts() {
	for _, next := range l.Next {
		sort.Ints(next)
	}
}

func (l *Lattice) GenToken() {
	if l.Spellouts == nil {
		panic("Can't generate token without a spellout")
	}
	if len(l.Spellouts) == 0 {
		l.Token = Token("")
		return
	}
	spellout := l.Spellouts[0]
	strs := make([]string, len(spellout))
	for i, morph := range spellout {
		strs[i] = morph.Form
	}
	l.Token = Token(strings.Join(strs, ""))
}

func (l *Lattice) GenSpellouts() {
	if l.Spellouts != nil {
		return
	}
	if len(l.Morphemes) == 0 {
		l.Spellouts = make(Spellouts, 0)
		return
	}
	var (
		pathId   int
		from, to int = l.Bottom(), l.Top()
	)
	l.Spellouts = make(Spellouts, 0, to-from)
	for path := range graph.YieldAllPaths(graph.DirectedGraph(l), from, to) {
		spellout := make(Spellout, len(path))
		for i, el := range path {
			spellout[i] = el.(*EMorpheme)
		}
		l.Spellouts = append(l.Spellouts, spellout)

		pathId++
	}
	sort.Sort(l.Spellouts)
}

func (l *Lattice) YieldPaths() chan Path {
	l.GenSpellouts()
	pathChan := make(chan Path)
	go func() {
		for i, _ := range l.Spellouts {
			pathChan <- Path(i)
		}
		close(pathChan)
	}()
	return pathChan
}

func (l *Lattice) Path(i int) Spellout {
	return l.Spellouts[i]
}

func (l *Lattice) AllFusedFrom(from int) []string {
	var (
		curMorphemeIds []int
		curMorpheme    *EMorpheme
		curNode        int = from
		exists         bool
	)
	forms := make([]string, 0, 3)
	allNext := []string{""}
	// log.Println("Fusing from", from, "in lattice", l.Token)
	for curNode < l.Top() {
		curMorphemeIds, exists = l.Next[curNode]
		if !exists {
			panic(fmt.Sprintf("Lattice does not have outgoing node Id %v", curNode))
		}
		if len(curMorphemeIds) > 1 {
			// log.Println("\tFound multiple outgoing")
			for _, curMorphemeId := range curMorphemeIds {
				curMorpheme = l.Morphemes[curMorphemeId]
				// log.Println("\tRecursing at", curMorpheme)
				for _, curStr := range l.AllFusedFrom(curMorpheme.To()) {
					allNext = append(allNext, curMorpheme.Form+curStr)
				}
			}
			break
		} else {
			curMorpheme = l.Morphemes[curMorphemeIds[0]]
			// log.Println("\tAt morpheme", curMorpheme)
			curNode = curMorpheme.To()
			if curMorpheme.From() > l.Bottom() && curMorpheme.Form == "H" {
				// log.Println("\tSkipping fused H")
				// fuse any H encountered mid lattice
				continue
			}
			forms = append(forms, curMorpheme.Form)
		}
	}
	if len(allNext) > 0 {
		results := make([]string, len(allNext))
		for i, val := range allNext {
			results[i] = strings.Join(forms, "") + val
		}
		// log.Println("Returning", results)
		return results
	} else {
		// log.Println("Returning", forms)
		return []string{strings.Join(forms, "")}
	}
}

type MorphDependencyGraph interface {
	LabeledDependencyGraph
	GetMappings() Mappings
	GetMorpheme(int) *EMorpheme
}
