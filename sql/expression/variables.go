// Copyright 2020-2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

// SystemVar is an expression that returns the value of a system variable. It's also used as the expression on the left
// hand side of a SET statement for a system variable.
type SystemVar struct {
	Name           string
	Collation      sql.CollationID
	Scope          sql.SystemVariableScope
	SpecifiedScope string
}

var _ sql.Expression = (*SystemVar)(nil)
var _ sql.CollationCoercible = (*SystemVar)(nil)

// NewSystemVar creates a new SystemVar expression for the system variable named |name| with the specified |scope|.
// The |specifiedScope| parameter indicates the exact scope that was specified in the original reference to this
// system variable, and is used to ensure we output a column name in a result set that exactly matches how the
// system variable was originally referenced. If the |specifiedScope| parameter is empty, then the scope was not
// originally specified and any scope has been inferred.
func NewSystemVar(name string, scope sql.SystemVariableScope, specifiedScope string) *SystemVar {
	return &SystemVar{name, sql.CollationID(0), scope, specifiedScope}
}

// Children implements the sql.Expression interface.
func (v *SystemVar) Children() []sql.Expression { return nil }

// Eval implements the sql.Expression interface.
func (v *SystemVar) Eval(ctx *sql.Context, _ sql.Row) (interface{}, error) {
	switch v.Scope {
	case sql.SystemVariableScope_Session:
		// "character_set_database" and "collation_database" are special system variables, in that they're set whenever
		// the current database is changed. Rather than attempting to synchronize the session variables of all
		// outstanding contexts whenever a database's collation is updated, we just pull the values from the database
		// directly. MySQL also plans to make these system variables immutable (from the user's perspective). This isn't
		// exactly the same as MySQL's behavior, but this is the intent of their behavior, which is also way easier to
		// implement.
		switch strings.ToLower(v.Name) {
		case "character_set_database":
			return v.Collation.CharacterSet().String(), nil
		case "collation_database":
			return v.Collation.String(), nil
		default:
			val, err := ctx.GetSessionVariable(ctx, v.Name)
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	case sql.SystemVariableScope_Global:
		_, val, ok := sql.SystemVariables.GetGlobal(v.Name)
		if !ok {
			return nil, sql.ErrUnknownSystemVariable.New(v.Name)
		}
		return val, nil
	default: // should never happen
		return nil, fmt.Errorf("unknown scope `%v` on system variable `%s`", v.Scope, v.Name)
	}
}

// Type implements the sql.Expression interface.
func (v *SystemVar) Type() sql.Type {
	if sysVar, _, ok := sql.SystemVariables.GetGlobal(v.Name); ok {
		return sysVar.Type
	}
	return types.Null
}

// CollationCoercibility implements the interface sql.CollationCoercible.
func (v *SystemVar) CollationCoercibility(ctx *sql.Context) (collation sql.CollationID, coercibility byte) {
	typ := v.Type()
	if types.IsText(typ) {
		collation, _ = typ.CollationCoercibility(ctx)
		return collation, 3
	}
	return typ.CollationCoercibility(ctx)
}

// IsNullable implements the sql.Expression interface.
func (v *SystemVar) IsNullable() bool { return false }

// Resolved implements the sql.Expression interface.
func (v *SystemVar) Resolved() bool { return true }

// String implements the sql.Expression interface.
func (v *SystemVar) String() string {
	// If the scope wasn't explicitly provided, then don't include it in the string representation
	if v.SpecifiedScope == "" {
		return fmt.Sprintf("@@%s", v.Name)
	} else {
		return fmt.Sprintf("@@%s.%s", v.SpecifiedScope, v.Name)
	}
}

// WithChildren implements the Expression interface.
func (v *SystemVar) WithChildren(children ...sql.Expression) (sql.Expression, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(v, len(children), 0)
	}
	return v, nil
}

// UserVar is an expression that returns the value of a user variable. It's also used as the expression on the left hand
// side of a SET statement for a user var.
type UserVar struct {
	Name     string
	exprType sql.Type
}

var _ sql.Expression = (*UserVar)(nil)
var _ sql.CollationCoercible = (*UserVar)(nil)

// NewUserVar creates a UserVar with a name, but no type information, for use as the left-hand value
// in a SetField assignment Expression. This method should not be used when the user variable is
// being used as a value, since the correct type information will not be available.
func NewUserVar(name string) *UserVar {
	return &UserVar{Name: name, exprType: types.Null}
}

// NewUserVarWithType creates a UserVar with its type resolved, so that it can be used as a value
// in other expressions.
func NewUserVarWithType(name string, t sql.Type) *UserVar {
	return &UserVar{Name: name, exprType: t}
}

// Children implements the sql.Expression interface.
func (v *UserVar) Children() []sql.Expression { return nil }

// Eval implements the sql.Expression interface.
func (v *UserVar) Eval(ctx *sql.Context, _ sql.Row) (interface{}, error) {
	_, val, err := ctx.GetUserVariable(ctx, v.Name)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// Type implements the sql.Expression interface.
func (v *UserVar) Type() sql.Type {
	return v.exprType
}

// CollationCoercibility implements the interface sql.CollationCoercible.
func (v *UserVar) CollationCoercibility(ctx *sql.Context) (collation sql.CollationID, coercibility byte) {
	collation, _ = v.exprType.CollationCoercibility(ctx)
	return collation, 2
}

// IsNullable implements the sql.Expression interface.
func (v *UserVar) IsNullable() bool { return true }

// Resolved implements the sql.Expression interface.
func (v *UserVar) Resolved() bool { return true }

// String implements the sql.Expression interface.
func (v *UserVar) String() string { return "@" + v.Name }

// WithChildren implements the Expression interface.
func (v *UserVar) WithChildren(children ...sql.Expression) (sql.Expression, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(v, len(children), 0)
	}
	return v, nil
}
