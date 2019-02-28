package webapi

import (
	"bytes"
	"fmt"
	"github.com/gonuts/commander"
	"io"
	"log"
	"strings"
	"sync"
	"yap/app"
	"yap/nlp/format/lattice"
	"yap/nlp/format/raw"
	"yap/nlp/parser/ma"
	"yap/nlp/parser/xliter8"
	nlp "yap/nlp/types"
	"yap/util"
)

var (
	maLock   sync.Mutex
	maHebrew xliter8.Interface
	maData   *ma.BGULex
)

func HebrewMorphAnalyazerInitialize(cmd *commander.Command, args []string) {
	prefixLocation, found := util.LocateFile(app.HebMaPrefixFile, app.HEB_MA_DEFAULT_DATA_DIRS)
	if !found {
		panic(fmt.Sprintf("Lexicon prefix file not found: %v", app.HebMaPrefixFile))
	}
	lexiconLocation, found := util.LocateFile(app.HebMaLexiconFile, app.HEB_MA_DEFAULT_DATA_DIRS)
	if !found {
		panic(fmt.Sprintf("Lexicon file not found: %v", app.HebMaLexiconFile))
	}
	app.HebMaPrefixFile = prefixLocation
	app.HebMaLexiconFile = lexiconLocation
	app.HebMAConfigOut()
	maData = new(ma.BGULex)
	maData.MAType = "spmrl"
	log.Println("Reading Morphological Analyzer BGU Prefixes")
	maData.LoadPrefixes(app.HebMaPrefixFile)
	log.Println("Reading Morphological Analyzer BGU Lexicon")
	maData.LoadLex(app.HebMaLexiconFile, app.HebMaNnpnofeats)
	log.Println()
	maData.AlwaysNNP = app.HebMaAlwaysnnp
	maData.LogOOV = app.HebMaShowoov

}

func HebrewMorphAnalyzeRawSentences(input string) string {
	maLock.Lock()
	var (
		reader io.Reader
		sents  []nlp.BasicSentence
		err    error
	)
	reader = strings.NewReader(input)
	sents, err = raw.Read(reader, 0)
	if err != nil {
		panic(fmt.Sprintf("Failed reading raw input - %v", err))
	}
	log.Println("Running Hebrew Morphological Analysis")
	log.Println("input:\n", input)
	stats := new(ma.AnalyzeStats)
	stats.Init()
	maData.Stats = stats
	//prefix := log.Prefix()
	lattices := make([]nlp.LatticeSentence, len(sents))
	oovInd := make([]interface{}, len(sents))
	for i, sent := range sents {
		//log.SetPrefix(fmt.Sprintf("%v graph# %v ", prefix, i))
		lattices[i], oovInd[i] = maData.Analyze(sent.Tokens())
	}
	log.Println()
	output := lattice.Sentence2LatticeCorpus(lattices, maHebrew)
	buf := new(bytes.Buffer)
	err = lattice.Write(buf, output)
	maLock.Unlock()
	return buf.String()
}

func HebrewMorphAnalyzeBasicSentences(sents []nlp.BasicSentence) string {
	maLock.Lock()

	stats := new(ma.AnalyzeStats)
	stats.Init()
	maData.Stats = stats

	lattices := make([]nlp.LatticeSentence, len(sents))
	//oovInd := make([]interface{}, len(sents))
	for i, sent := range sents {
		lattices[i], _ = maData.Analyze(sent.Tokens())
	}

	output := lattice.Sentence2LatticeCorpus(lattices, maHebrew)
	buf := new(bytes.Buffer)
	lattice.Write(buf, output)

	maLock.Unlock()
	return buf.String()
}
