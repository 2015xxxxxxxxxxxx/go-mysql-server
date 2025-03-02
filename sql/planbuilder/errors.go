// Copyright 2023 Dolthub, Inc.
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

package planbuilder

import (
	"regexp"

	"gopkg.in/src-d/go-errors.v1"
)

var (
	errInvalidDescribeFormat = errors.NewKind("invalid format %q for DESCRIBE, supported formats: %s")

	errInvalidSortOrder = errors.NewKind("invalid sort order: %s")

	ErrPrimaryKeyOnNullField = errors.NewKind("All parts of PRIMARY KEY must be NOT NULL")

	// TODO: We parse table options in Vitess, but we return them to GMS as a single string, so GMS has to reparse
	//       table options using these regexes. It would be cleaner for Vitess to parse them into structures so that
	//       GMS could just pull out the data it needs, without having to regex it out from a single string. Because
	//       of how the original quotes are lost and single quotes are always used, we also currently lose some info
	//       about the original statement – notably, we can't parse comments that contain single quotes because Vitess
	//       strips the original double quotes and sends the comment string to GMS wrapped in single quotes.
	tableCharsetOptionRegex   = regexp.MustCompile(`(?i)(DEFAULT)?\s+CHARACTER\s+SET((\s*=?\s*)|\s+)([A-Za-z0-9_]+)`)
	tableCollationOptionRegex = regexp.MustCompile(`(?i)(DEFAULT)?\s+COLLATE((\s*=?\s*)|\s+)([A-Za-z0-9_]+)`)
	tableCommentOptionRegex   = regexp.MustCompile(`(?i)\s+COMMENT((\s*=?\s*)|\s+)('([^']+)')`)

	// ErrUnionSchemasDifferentLength is returned when the two sides of a
	// UNION do not have the same number of columns in their schemas.
	ErrUnionSchemasDifferentLength = errors.NewKind(
		"cannot union two queries whose schemas are different lengths; left has %d column(s) right has %d column(s).",
	)

	ErrQualifiedOrderBy = errors.NewKind("Table '%s' from one of the SELECTs cannot be used in global ORDER clause")

	ErrOrderByBinding = errors.NewKind("bindings in sort clauses not supported yet")

	ErrFailedToParseStats = errors.NewKind("failed to parse data: %s\n%s")
)
