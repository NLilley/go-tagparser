```
  ______          _______            ______                            
 / _____)        (_______)          (_____ \                           
| /  ___  ___ ___ _       ____  ____ _____) )___  ____ ___  ____  ____ 
| | (___)/ _ (___) |     / _  |/ _  |  ____/ _  |/ ___)___)/ _  )/ ___)
| \____/| |_| |  | |____( ( | ( ( | | |   ( ( | | |  |___ ( (/ /| |    
 \_____/ \___/    \______)_||_|\_|| |_|    \_||_|_|  (___/ \____)_|    
                              (_____|                                  

                                                            Go-TagParer
```

A small parser capable of reading html/xml style tag-based documents.

## Project Overview

The main parser is the `Parse` function located in the `parser` package/parser.go file. To use it, simply call the method with a runeified document input.

The parser output is tree of tags representing your document. Each tag has one or more children representing it's child tags. Where the child is raw text, a pseudo tag will be created for it with name `<text>` and attribute `text = content`

Alternatively, you can use the [/cmd/tagStat](/cmd/tagStat/README.md) command to provide a concise summary of the documents contents. To install, run `go install ./cmd/tagStat`. See the README for more detailed usage instructions.

## Parser Specification

The parser expects that:
- There is a single root tag
- Escaped characters should be left in-tact (i.e. &lt; won't be transformed to "<")

Note:
- Leading and trailing space characters will be stripped before processing
- Self closing tags are permitted
- Nameless tags are permitted (but must have no attributes). i.e. `</>` and `<>MyContent</>`
- Attributes can use either single and double quotes
- Embedded text content will have a Tag.name of `<text>`.
- - i.e. For `<p>Content</p>`, Content will be wrapped into a tag with `Tag.name = "<text>"` and attribute text == `"Content"`
 - Empty tag attributes (valueless attributes) are not supported (i.e. `<checkbox checked/>`)
- Unicode is mostly support in attribute names, values and the like
- Whitespace is stripped from either side of raw text content


## Run Tests:
```
cd ./pkg
go test
```

## Run Benchmarks:
```
cd ./pkg
go test -bench . -benchmem -count 10 > 10_runs_bench.txt
```