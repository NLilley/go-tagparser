# Stat - Tag Parser Stat Command

This small utility will attempt to parse and provide summary statistics for a tag document. If the tag document is incorrectly formed, the error will be logged to standard out.

By default the document will be read from `stdin`. If you provide the `-i="PATH"` argument, it will read the document from the provided file instead.

## Examples
### Standard In
```
echo "<happy><people /></happy>" | go run stat.go
```

### File
```
go run stat.go -i="./test.html"
```

### Potential Output
```
2024/03/18 19:45:14 Input:
<html lang="en">
    <head>
        <meta charset="utf-8" />
        <meta name="robots" content="noindex,nofollow" />
        <title>The Bestest Site Ever!</title>
    </head>
    <body>
        <h1>Graphic Design is my Passion!</h1>
        <a href="https://www.youtube.com/watch?v=p3G5IXn0K7A">🐹</a>
        <span style="color:red">Culture</span> and <span style="color: purple">Beauty</span>
    </body>
</html>

2024/03/18 19:45:14
Tag Document Statistics:
------------------------
Total Tags: 10
Total Text Contents: 6
Total Attributes: 7

Tag Histogram:
        meta    2
        title   1
        body    1
        h1      1
        a       1
        span    2
        html    1
        head    1

Attribute Histogram:
        lang    1
        charset 1
        name    1
        content 1
        href    1
        style   2
```
