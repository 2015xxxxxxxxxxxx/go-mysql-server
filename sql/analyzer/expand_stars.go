package analyzer

import (
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

// expandStars replaces star expressions into lists of concrete column expressions
func expandStars(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope) (sql.Node, error) {
	span, _ := ctx.Span("expand_stars")
	defer span.Finish()

	tableAliases, err := getTableAliases(n, scope)
	if err != nil {
		return nil, err
	}

	return plan.TransformUp(n, func(n sql.Node) (sql.Node, error) {
		if n.Resolved() {
			return n, nil
		}

		switch n := n.(type) {
		case *plan.Project:
			if !n.Child.Resolved() {
				return n, nil
			}

			expressions, err := expandStarsForExpressions(a, n.Projections, n.Child.Schema(), tableAliases)
			if err != nil {
				return nil, err
			}

			return plan.NewProject(expressions, n.Child), nil
		case *plan.GroupBy:
			if !n.Child.Resolved() {
				return n, nil
			}

			aggregate, err := expandStarsForExpressions(a, n.SelectedExprs, n.Child.Schema(), tableAliases)
			if err != nil {
				return nil, err
			}

			return plan.NewGroupBy(aggregate, n.GroupByExprs, n.Child), nil
		default:
			return n, nil
		}
	})
}

func expandStarsForExpressions(a *Analyzer, exprs []sql.Expression, schema sql.Schema, tableAliases TableAliases) ([]sql.Expression, error) {
	var expressions []sql.Expression
	for _, e := range exprs {
		if star, ok := e.(*expression.Star); ok {
			var exprs []sql.Expression
			for i, col := range schema {
				lowerSource := strings.ToLower(col.Source)
				lowerTable := strings.ToLower(star.Table)
				if star.Table == "" || lowerTable == lowerSource {
					exprs = append(exprs, expression.NewGetFieldWithTable(
						i, col.Type, col.Source, col.Name, col.Nullable,
					))
				}
			}

			if len(exprs) == 0 && star.Table != "" {
				return nil, sql.ErrTableNotFound.New(star.Table)
			}

			expressions = append(expressions, exprs...)
		} else {
			expressions = append(expressions, e)
		}
	}

	a.Log("resolved * to expressions %s", expressions)
	return expressions, nil
}
