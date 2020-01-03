package generate

import (
	"fmt"
	"io/ioutil"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
)

func getSchema(filename string) (*ast.Schema, error) {
	// TODO: IRL we have to get the schema from GraphQL (maybe we can generate
	// that once we can bootstrap) where it comes as JSON, not SDL, so we have
	// to convert (or add gqlparser support to convert)
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unreadable schema file %v: %v", filename, err)
	}

	schema, graphqlError := gqlparser.LoadSchema(
		&ast.Source{Name: filename, Input: string(text)})
	if graphqlError != nil {
		return nil, fmt.Errorf("invalid schema file %v: %v",
			filename, graphqlError)
	}

	return schema, nil
}

func getAndValidateQueries(filename string, schema *ast.Schema) (*ast.QueryDocument, error) {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unreadable query-spec file %v: %v", filename, err)
	}

	// The following is more or less gqlparser.LoadQuery, but we can provide a
	// name so we might as well (and we break out the two errors).
	document, graphqlError := parser.ParseQuery(
		&ast.Source{Name: filename, Input: string(text)})
	if graphqlError != nil { // ParseQuery returns type *graphql.Error, yuck
		return nil, fmt.Errorf("invalid query-spec file %v: %v", filename, graphqlError)
	}

	graphqlErrors := validator.Validate(schema, document)
	if graphqlErrors != nil {
		return nil, fmt.Errorf("query-spec does not match schema: %v", graphqlErrors)
	}

	return document, nil
}
