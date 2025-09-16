# PHP SPX Graph Visualizer

SPX Graph Visualizer is a tool for analyzing and visualizing performance profiles created with [SPX PHP Profiler](https://github.com/NoiseByNorthwest/php-spx).

## Description
I really don't like viewing the profile as a Flamegraph. I think it's difficult.

## Features
- **Graph visualization** — interactive call graphs with function names
- **Zoom** — support zoom on the graph
- **Compression txt support** — works with .txt and .txt.gz files
- **Report** — view in browser or save to HTML report

## Usage

### Start web server
```bash
./spx-graph --file profile.txt.gz # http://localhost:8080

./spx-graph --file profile.txt.gz --port 9090 # http://localhost:9090
```

### Save HTML report
```bash
./spx-graph --file profile.txt.gz -o result.html
```


## Input data format

Works with SPX profiles containing sections:
```
[events]
0 1 0 0
1 1 33780 272
1 0 39230 536
...

[functions]
/var/www/html/index.php
/var/www/html/vendor/autoload.php
MyClass::myMethod
...
```

## License
MIT License

## Contributing

Pull requests and issues are welcome!

## Links
- [PHP SPX](https://github.com/NoiseByNorthwest/php-spx) — PHP profiler
