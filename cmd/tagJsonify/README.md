# TagJsonify - Dump TAG documents to JSON

This small utility will attempt to parse a tag document and dump it to JSON. All output will be written to Standard Out - To capture it, pipe it to a file.

If you provide the `-i="PATH"` argument, the program will read the document from the provided file. If you provide the `-stdin` argument, the input will be read from standard in. Otherwise help information will be provided.

## Examples
### Standard In
```
echo "<happy><people /></happy>" | go run tagJsonify.go -stdin
```

### File
```
go run tagJsonify.go -i="./test.html"
```

### Potential Output
```
{
    "_name": "happy",
    "_children": [
        {
            "_name": "people"
        }
    ]
}
```