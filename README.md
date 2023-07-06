# AutoDocs

AutoDocs is a command-line tool that uses OpenAI to automatically generate documentation for Go functions that don't have comments. It can process a single Go file or all Go files in a directory.

## Requirements

AutoDocs uses the OpenAI API to generate documentation. To use AutoDocs, you need to have an OpenAI API key and set it as an environment variable:

```bash
export OPENAI_API_KEY=your-api-key
```

Replace your-api-key with your actual OpenAI API key.

## Installation

You can install AutoDocs directly using the go install command:
```bash
go install github.com/tbuchaillot/autodocs@latest
```

This will download the source code, compile it, and move the autodocs executable to $GOPATH/bin.

Make sure `$GOPATH/bin` is in your PATH to run the autodocs command from anywhere.

## Usage

You can run AutoDocs on a single Go file with the `-f` flag:

```bash
autodocs -f myfile.go
```

This will add documentation to myfile.go functions.

You can also run AutoDocs on all Go files in a directory with the `-d` flag:

```bash
autodocs -d mydir
```

This will add documentation to all .go files functions in mydir.


## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

