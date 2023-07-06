package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCompletion struct {
	mock.Mock
}

func (m *MockCompletion) CreateComment(fset *token.FileSet, funcDecl *ast.FuncDecl) string {
	args := m.Called(fset, funcDecl)
	return args.String(0)
}

func TestAddCommentsIfMissing(t *testing.T) {
	mockCompletion := new(MockCompletion)
	fset := token.NewFileSet()
	src := `
		package main
		func HelloWorld() {
			println("Hello, World!")
		}
	`
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	var funcDecl *ast.FuncDecl
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			funcDecl = fd
			break
		}
	}

	if funcDecl == nil {
		t.Fatal("No function declaration found")
	}

	mockCompletion.On("CreateComment", mock.Anything, mock.Anything).Return("Test comment")

	tmpfile, err := ioutil.TempFile("", "example.*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some Go code to the temporary file
	if _, err := tmpfile.Write([]byte(src)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	err = addCommentsIfMissing(tmpfile.Name(), mockCompletion)
	assert.NoError(t, err)
}

func TestWriteFile(t *testing.T) {
	fset := token.NewFileSet()
	src := `
		package main
		func HelloWorld() {
			println("Hello, World!")
		}
	`
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "example.*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	err = WriteFile(tmpfile.Name(), fset, file)
	assert.NoError(t, err)
}
