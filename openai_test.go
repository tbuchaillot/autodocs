package main

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) CreateChatCompletion(
	ctx context.Context,
	request openai.ChatCompletionRequest,
) (response openai.ChatCompletionResponse, err error) {
	args := m.Called(ctx, request)
	return args.Get(0).(openai.ChatCompletionResponse), args.Error(1)
}

func TestCreateComment(t *testing.T) {
	mockClient := new(MockClient)
	completion := &completion{
		client: mockClient,
	}

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

	mockClient.On("CreateChatCompletion", mock.Anything, mock.Anything).Return(
		openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: "Test comment",
					},
				},
			},
		}, nil)

	comment := completion.CreateComment(fset, funcDecl)
	assert.Equal(t, "Test comment", comment)
}
