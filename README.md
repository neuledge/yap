
# YAP Service

An RestAPI server based on the [YAP parser](https://github.com/OnlpLab/yap)

![Joint Morph-Syntactic Processing](screenshot.png)

A live demo of parsing Hebrew texts is provided [here](http://onlp.openu.org.il/).

## Quick Start

YAP can run on Windows, Linux and MacOS.

**Windows users:** YAP doesn't handle Windows style text files that have [BOM](https://en.wikipedia.org/wiki/Byte_order_mark) marks and [CRLF](https://en.wikipedia.org/wiki/Newline) newlines.
So if you're running on Windows and YAP doesn't work make sure you don't have CRLF line endings and no BOM mark.

### Requirements

- [Go](http://www.golang.org)
- [dep](https://golang.github.io/dep/)
- [Git](https://git-scm.com/downloads)
- bzip2
- 6GB RAM

### Compilation

The following instructions are for Linux but similarly this can be done on Windows and MacOS.

- Make sure you have Go and Git installed and on the command PATH.
- Setup a Go environment:
  - Create a directory (usually per workspace/project) ``mkdir yapproj; cd yapproj``
  - Set ``$GOPATH`` environment variable to your workspace: ``export GOPATH=path/to/yapproj``
  - In the workspace directory create the src subdirectory: ``mkdir src``
  - cd into the src directory ``cd src``
- Clone the repository in the src folder of the workspace ``git clone https://github.com/neuledge/yap-service.git``
- Unzip the models and build the application:

```console
$ cd yap
$ bunzip2 data/*.bz2
$ dep ensure
$ go build .
$ ./yap
./yap - invoke yap as a standalone app or as an api server

Commands:

    api         start api server
    dep         runs dependency training/parsing
    hebma       run lexicon-based morphological analyzer on raw input
    joint       runs joint morpho-syntactic training and parsing
    ma          run data-driven morphological analyzer on raw input
    malearn     generate a data-driven morphological analysis dictionary for a set of files
    md          runs standalone morphological disambiguation training and parsing

Use "./yap help <command>" for more information about a command
```

### Running YAP as a RESTful API server

1. YAP can run as a server listening on port 8000:

    ```console
    $ ./yap api
    ```

2. You can then send HTTP GET requests with json objects in the request body and receive back a json object containing the 3 output levels:

    ```
POST /parse

{
	"sentences": [
		["כשעה", "נסיעה", "ממיאמי", ",", "נמצא", "מרכז", "הקניות", "הגדול", "הזה", ",", "שהוא", "למעשה", "מתחם", "של", "מאות", "חנויות", "ואאוטלטים", "."],
		["מתחם", "קניות"]
	]
}
    ```
    
    Response:
    ```json
[
    [
        {
            "token": 0,
            "form": "כ",
            "lemma": "כ",
            "CPOS": "PREPOSITION",
            "POS": "PREPOSITION",
            "features": {},
            "head": 6,
            "dep": "prepmod"
        },
        {
            "token": 0,
            "form": "שעה",
            "lemma": "שעה",
            "CPOS": "NN",
            "POS": "NN",
            "features": {
                "gen": "F",
                "num": "S"
            },
            "head": 0,
            "dep": "pobj"
        },
        {
            "token": 1,
            "form": "נסיעה",
            "lemma": "נסיעה",
            "CPOS": "NN",
            "POS": "NN",
            "features": {
                "gen": "F",
                "num": "S"
            },
            "head": 1,
            "dep": "nn"
        },
        {
            "token": 2,
            "form": "מ",
            "lemma": "מ",
            "CPOS": "PREPOSITION",
            "POS": "PREPOSITION",
            "features": {},
            "head": 2,
            "dep": "prepmod"
        },
        {
            "token": 2,
            "form": "מיאמי",
            "lemma": "מיאמי",
            "CPOS": "NNP",
            "POS": "NNP",
            "features": {},
            "head": 3,
            "dep": "pobj"
        },
        {
            "token": 3,
            "form": ",",
            "lemma": "",
            "CPOS": "yyCM",
            "POS": "yyCM",
            "features": {},
            "head": 6,
            "dep": "punct"
        },
        {
            "token": 4,
            "form": "נמצא",
            "lemma": "נמצא",
            "CPOS": "VB",
            "POS": "VB",
            "features": {
                "gen": "M",
                "num": "S",
                "per": "3",
                "tense": "PAST"
            },
            "head": -1,
            "dep": "ROOT"
        },
        {
            "token": 5,
            "form": "מרכז",
            "lemma": "מרכז",
            "CPOS": "NNT",
            "POS": "NNT",
            "features": {
                "gen": "M",
                "num": "S"
            },
            "head": 6,
            "dep": "subj"
        },
        {
            "token": 6,
            "form": "ה",
            "lemma": "ה",
            "CPOS": "DEF",
            "POS": "DEF",
            "features": {},
            "head": 9,
            "dep": "def"
        },
        {
            "token": 6,
            "form": "קניות",
            "lemma": "קנייה",
            "CPOS": "NN",
            "POS": "NN",
            "features": {
                "gen": "F",
                "num": "P"
            },
            "head": 7,
            "dep": "gobj"
        },
        {
            "token": 7,
            "form": "ה",
            "lemma": "ה",
            "CPOS": "DEF",
            "POS": "DEF",
            "features": {},
            "head": 11,
            "dep": "def"
        },
        {
            "token": 7,
            "form": "גדול",
            "lemma": "גדול",
            "CPOS": "JJ",
            "POS": "JJ",
            "features": {
                "gen": "M",
                "num": "S"
            },
            "head": 7,
            "dep": "amod"
        },
        {
            "token": 8,
            "form": "ה",
            "lemma": "ה",
            "CPOS": "DEF",
            "POS": "DEF",
            "features": {},
            "head": 13,
            "dep": "def"
        },
        {
            "token": 8,
            "form": "זה",
            "lemma": "זה",
            "CPOS": "PRP",
            "POS": "PRP",
            "features": {
                "gen": "M",
                "num": "S",
                "per": "3"
            },
            "head": 7,
            "dep": "amod"
        },
        {
            "token": 9,
            "form": ",",
            "lemma": "",
            "CPOS": "yyCM",
            "POS": "yyCM",
            "features": {},
            "head": 7,
            "dep": "punct"
        },
        {
            "token": 10,
            "form": "ש",
            "lemma": "ש",
            "CPOS": "REL",
            "POS": "REL",
            "features": {},
            "head": 7,
            "dep": "rcmod"
        },
        {
            "token": 10,
            "form": "הוא",
            "lemma": "הוא",
            "CPOS": "PRP",
            "POS": "PRP",
            "features": {
                "gen": "M",
                "num": "S",
                "per": "3"
            },
            "head": 18,
            "dep": "subj"
        },
        {
            "token": 11,
            "form": "למעשה",
            "lemma": "למעשה",
            "CPOS": "RB",
            "POS": "RB",
            "features": {},
            "head": 18,
            "dep": "parataxis"
        },
        {
            "token": 12,
            "form": "מתחם",
            "lemma": "מתחם",
            "CPOS": "NN",
            "POS": "NN",
            "features": {
                "gen": "M",
                "num": "S"
            },
            "head": 15,
            "dep": "relcomp"
        },
        {
            "token": 13,
            "form": "של",
            "lemma": "של",
            "CPOS": "POS",
            "POS": "POS",
            "features": {},
            "head": 18,
            "dep": "posspmod"
        },
        {
            "token": 14,
            "form": "מאות",
            "lemma": "מאה",
            "CPOS": "CDT",
            "POS": "CDT",
            "features": {
                "gen": "F",
                "num": "P"
            },
            "head": 21,
            "dep": "num"
        },
        {
            "token": 15,
            "form": "חנויות",
            "lemma": "חנות",
            "CPOS": "NN",
            "POS": "NN",
            "features": {
                "gen": "F",
                "num": "P"
            },
            "head": 22,
            "dep": "conj"
        },
        {
            "token": 16,
            "form": "ו",
            "lemma": "ו",
            "CPOS": "CONJ",
            "POS": "CONJ",
            "features": {},
            "head": 19,
            "dep": "gobj"
        },
        {
            "token": 16,
            "form": "אאוטלטים",
            "lemma": "אאוטלטים",
            "CPOS": "NN",
            "POS": "NN",
            "features": {
                "gen": "F",
                "num": "S"
            },
            "head": 22,
            "dep": "conj"
        },
        {
            "token": 17,
            "form": ".",
            "lemma": "",
            "CPOS": "yyDOT",
            "POS": "yyDOT",
            "features": {},
            "head": 6,
            "dep": "punct"
        }
    ],
    [
        {
            "token": 0,
            "form": "מתחם",
            "lemma": "מתחם",
            "CPOS": "NNT",
            "POS": "NNT",
            "features": {
                "gen": "M",
                "num": "S"
            },
            "head": -1,
            "dep": "ROOT"
        },
        {
            "token": 1,
            "form": "קניות",
            "lemma": "קנייה",
            "CPOS": "NN",
            "POS": "NN",
            "features": {
                "gen": "F",
                "num": "P"
            },
            "head": 0,
            "dep": "gobj"
        }
    ]
]
    ```

## License

This software is released under the terms of the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0).

The Apache license does not apply to the BGU Lexicon. Please contact Reut Tsarfaty regarding licensing of the lexicon.
