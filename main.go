package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocomment",
	Short: "gocomment adds TODO comments to Go functions without comments",
	Long: `A Fast and Flexible Commenter built with
                love by in go.`,
	Run: func(cmd *cobra.Command, args []string) {
		key := os.Getenv("OPENAI_API_KEY")
		if key == "" {
			log.Fatal("Please set OPENAI_API_KEY environment variable")
			return
		}
		completion := NewCompletion(key)

		// Get flags
		file, _ := cmd.Flags().GetString("file")
		dir, _ := cmd.Flags().GetString("dir")

		if file != "" && strings.HasSuffix(file, ".go") {
			addCommentsIfMissing(file, completion)
		} else if dir != "" {
			err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
					addCommentsIfMissing(path, completion)
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Please provide a file or directory")
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringP("file", "f", "", "Go file to process")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "Directory to process")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func addCommentsIfMissing(filename string, completion Completion) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	comments := []*ast.CommentGroup{}
	ast.Inspect(node, func(n ast.Node) bool {
		// collect comments
		c, ok := n.(*ast.CommentGroup)
		if ok {
			comments = append(comments, c)
		}
		// handle function declarations without documentation
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			if fn.Doc.Text() == "" {
				comment := &ast.Comment{
					Text:  completion.CreateComment(fset, fn),
					Slash: fn.Pos() - 1,
				}
				// create CommentGroup and set it to the function's documentation comment
				cg := &ast.CommentGroup{
					List: []*ast.Comment{comment},
				}
				fn.Doc = cg
				log.Printf("adding comment to %s func %s", filename, fn.Name.Name)
			}
		}
		return true
	})
	// set ast's comments to the collected comments
	node.Comments = comments

	// write new ast to file
	err = WriteFile(filename, fset, node)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func WriteFile(filename string, fset *token.FileSet, content *ast.File) error {
	// write new ast to file
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := printer.Fprint(f, fset, content); err != nil {
		return err
	}
	return nil
}
