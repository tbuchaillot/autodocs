package main

import (
	"bytes"
	"context"
	"go/ast"
	"go/printer"
	"go/token"
	"html/template"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

type Completion interface {
	CreateComment(fset *token.FileSet, funcDecl *ast.FuncDecl) string
}

type completion struct {
	client Client
}

type Client interface {
	CreateChatCompletion(
		ctx context.Context,
		request openai.ChatCompletionRequest,
	) (response openai.ChatCompletionResponse, err error)
}

func NewCompletion(key string) Completion {
	return &completion{
		client: openai.NewClient(key),
	}
}

func (c *completion) CreateComment(fset *token.FileSet, funcDecl *ast.FuncDecl) string {
	buf := new(bytes.Buffer)

	printer.Fprint(buf, fset, funcDecl)

	// Get the function declaration string from the buffer
	fnString := buf.String()

	// Reset the buffer for the next iteration
	buf.Reset()

	data := struct{ Fn string }{
		Fn: fnString,
	}

	tmpl, err := template.New("test").Parse(fnTmpl)
	if err != nil {
		log.Fatal(err)
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		log.Fatal(err)
	}

	return c.generateComment(tpl.String())
}

func (c *completion) generateComment(input string) string {
	msgs := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	}}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    msgs,
			Temperature: 1.0,
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

var fnTmpl = `
You are a Golang DX expert specialized in documentation and good practices. You are asked to write a comment for a function, only the comment, NOT the function itself.
The composition of the comment MUST have the following structure:
// {function detailed description}
// {function usage Example}


The function is:
{{.Fn}}
`
