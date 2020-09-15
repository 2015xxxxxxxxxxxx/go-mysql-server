package plan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liquidata-inc/go-mysql-server/sql"
	"github.com/liquidata-inc/go-mysql-server/sql/expression"
)

func TestSet(t *testing.T) {
	require := require.New(t)

	ctx := sql.NewContext(context.Background(), sql.WithSession(sql.NewBaseSession()))

	s := NewSet(
		[]sql.Expression {
			expression.NewSetField(expression.NewUnresolvedColumn("foo"), expression.NewLiteral("bar", sql.LongText)),
			expression.NewSetField(expression.NewUnresolvedColumn("@@baz"), expression.NewLiteral(int64(1), sql.Int64)),
		}...,
	)

	_, err := s.RowIter(ctx, nil)
	require.NoError(err)

	typ, v := ctx.Get("foo")
	require.Equal(sql.LongText, typ)
	require.Equal("bar", v)

	typ, v = ctx.Get("baz")
	require.Equal(sql.Int64, typ)
	require.Equal(int64(1), v)
}

func TestSetDesfault(t *testing.T) {
	require := require.New(t)

	ctx := sql.NewContext(context.Background(), sql.WithSession(sql.NewBaseSession()))

	s := NewSet(
		[]sql.Expression{
			expression.NewSetField(expression.NewUnresolvedColumn("auto_increment_increment"), expression.NewLiteral(int64(123), sql.Int64)),
			expression.NewSetField(expression.NewUnresolvedColumn("@@sql_select_limit"), expression.NewLiteral(int64(1), sql.Int64)),
		}...,
	)

	_, err := s.RowIter(ctx, nil)
	require.NoError(err)

	typ, v := ctx.Get("auto_increment_increment")
	require.Equal(sql.Int64, typ)
	require.Equal(int64(123), v)

	typ, v = ctx.Get("sql_select_limit")
	require.Equal(sql.Int64, typ)
	require.Equal(int64(1), v)

	s = NewSet(
		expression.NewSetField(expression.NewUnresolvedColumn("auto_increment_increment"), expression.NewDefaultColumn("")),
		expression.NewSetField(expression.NewUnresolvedColumn("@@sql_select_limit"), expression.NewDefaultColumn("")),
	)

	_, err = s.RowIter(ctx, nil)
	require.NoError(err)

	defaults := sql.DefaultSessionConfig()

	typ, v = ctx.Get("auto_increment_increment")
	require.Equal(defaults["auto_increment_increment"].Typ, typ)
	require.Equal(defaults["auto_increment_increment"].Value, v)

	typ, v = ctx.Get("sql_select_limit")
	require.Equal(defaults["sql_select_limit"].Typ, typ)
	require.Equal(defaults["sql_select_limit"].Value, v)

}
