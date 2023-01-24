// Copyright 2022 Dolthub, Inc.
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

package plan

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/binlogreplication"
)

// BinlogReplicaControllerCommand represents a SQL statement that requires a BinlogReplicaController
// (e.g. Start Replica, Show Replica Status).
type BinlogReplicaControllerCommand interface {
	sql.Node

	// WithBinlogReplicaController returns a new instance of this BinlogReplicaController, with the binlog replica
	// controller configured.
	WithBinlogReplicaController(controller binlogreplication.BinlogReplicaController) sql.Node
}

// ChangeReplicationSource is the plan node for the "CHANGE REPLICATION SOURCE TO" statement.
// https://dev.mysql.com/doc/refman/8.0/en/change-replication-source-to.html
type ChangeReplicationSource struct {
	replicaController binlogreplication.BinlogReplicaController
	Options           []binlogreplication.ReplicationOption
}

var _ sql.Node = (*ChangeReplicationSource)(nil)
var _ BinlogReplicaControllerCommand = (*ChangeReplicationSource)(nil)

func NewChangeReplicationSource(options []binlogreplication.ReplicationOption) *ChangeReplicationSource {
	return &ChangeReplicationSource{
		Options: options,
	}
}

// WithBinlogReplicaController implements the BinlogReplicaControllerCommand interface.
func (c *ChangeReplicationSource) WithBinlogReplicaController(controller binlogreplication.BinlogReplicaController) sql.Node {
	nc := *c
	nc.replicaController = controller
	return &nc
}

func (c *ChangeReplicationSource) Resolved() bool {
	return true
}

func (c *ChangeReplicationSource) String() string {
	sb := strings.Builder{}
	sb.WriteString("CHANGE REPLICATION SOURCE TO ")
	for i, option := range c.Options {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s = %s", option.Name, option.Value))
	}
	return sb.String()
}

func (c *ChangeReplicationSource) Schema() sql.Schema {
	return nil
}

func (c *ChangeReplicationSource) Children() []sql.Node {
	return nil
}

func (c *ChangeReplicationSource) RowIter(ctx *sql.Context, _ sql.Row) (sql.RowIter, error) {
	if c.replicaController == nil {
		return nil, fmt.Errorf("no replication controller available")
	}

	err := c.replicaController.SetReplicationSourceOptions(ctx, c.Options)
	return sql.RowsToRowIter(), err
}

func (c *ChangeReplicationSource) WithChildren(children ...sql.Node) (sql.Node, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(c, len(children), 0)
	}

	newNode := *c
	return &newNode, nil
}

func (c *ChangeReplicationSource) CheckPrivileges(_ *sql.Context, _ sql.PrivilegedOperationChecker) bool {
	// TODO: implement privilege checks
	return true
}

// ChangeReplicationFilter is a plan node for the "CHANGE REPLICATION FILTER" statement.
// https://dev.mysql.com/doc/refman/8.0/en/change-replication-filter.html
type ChangeReplicationFilter struct {
	replicaController binlogreplication.BinlogReplicaController
	Options           []binlogreplication.ReplicationOption
}

var _ sql.Node = (*ChangeReplicationFilter)(nil)
var _ BinlogReplicaControllerCommand = (*ChangeReplicationFilter)(nil)

func NewChangeReplicationFilter(options []binlogreplication.ReplicationOption) *ChangeReplicationFilter {
	return &ChangeReplicationFilter{
		Options: options,
	}
}

// WithBinlogReplicaController implements the BinlogReplicaControllerCommand interface.
func (c *ChangeReplicationFilter) WithBinlogReplicaController(controller binlogreplication.BinlogReplicaController) sql.Node {
	nc := *c
	nc.replicaController = controller
	return &nc
}

func (c *ChangeReplicationFilter) Resolved() bool {
	return true
}

func (c *ChangeReplicationFilter) String() string {
	sb := strings.Builder{}
	sb.WriteString("CHANGE REPLICATION FILTER ")
	for i, option := range c.Options {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(option.Name)
		sb.WriteString(" = ")
		// TODO: Fix this to use better typing
		sb.WriteString(fmt.Sprintf("%s", option.Value))
	}
	return sb.String()
}

func (c *ChangeReplicationFilter) Schema() sql.Schema {
	return nil
}

func (c *ChangeReplicationFilter) Children() []sql.Node {
	return nil
}

func (c *ChangeReplicationFilter) RowIter(ctx *sql.Context, row sql.Row) (sql.RowIter, error) {
	if c.replicaController == nil {
		return nil, fmt.Errorf("no replication controller available")
	}

	err := c.replicaController.SetReplicationFilterOptions(ctx, c.Options)
	return sql.RowsToRowIter(), err
}

func (c *ChangeReplicationFilter) WithChildren(children ...sql.Node) (sql.Node, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(c, len(children), 0)
	}

	newNode := *c
	return &newNode, nil
}

func (c *ChangeReplicationFilter) CheckPrivileges(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	// TODO: implement privilege checks
	return true
}

// StartReplica is a plan node for the "START REPLICA" statement.
// https://dev.mysql.com/doc/refman/8.0/en/start-replica.html
type StartReplica struct {
	replicaController binlogreplication.BinlogReplicaController
}

var _ sql.Node = (*StartReplica)(nil)
var _ BinlogReplicaControllerCommand = (*StartReplica)(nil)

func NewStartReplica() *StartReplica {
	return &StartReplica{}
}

// WithBinlogReplicaController implements the BinlogReplicaControllerCommand interface.
func (s *StartReplica) WithBinlogReplicaController(controller binlogreplication.BinlogReplicaController) sql.Node {
	nc := *s
	nc.replicaController = controller
	return &nc
}

func (s *StartReplica) Resolved() bool {
	return true
}

func (s *StartReplica) String() string {
	return "START REPLICA"
}

func (s *StartReplica) Schema() sql.Schema {
	return nil
}

func (s *StartReplica) Children() []sql.Node {
	return nil
}

func (s *StartReplica) RowIter(ctx *sql.Context, _ sql.Row) (sql.RowIter, error) {
	if s.replicaController == nil {
		return nil, fmt.Errorf("no replication controller available")
	}

	err := s.replicaController.StartReplica(ctx)
	return sql.RowsToRowIter(), err
}

func (s *StartReplica) WithChildren(children ...sql.Node) (sql.Node, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(s, len(children), 0)
	}

	newNode := *s
	return &newNode, nil
}

func (s *StartReplica) CheckPrivileges(_ *sql.Context, _ sql.PrivilegedOperationChecker) bool {
	// TODO: implement privilege checks
	return true
}

// StopReplica is the plan node for the "STOP REPLICA" statement.
// https://dev.mysql.com/doc/refman/8.0/en/stop-replica.html
type StopReplica struct {
	replicaController binlogreplication.BinlogReplicaController
}

var _ sql.Node = (*StopReplica)(nil)
var _ BinlogReplicaControllerCommand = (*StopReplica)(nil)

func NewStopReplica() *StopReplica {
	return &StopReplica{}
}

// WithBinlogReplicaController implements the BinlogReplicaControllerCommand interface.
func (s *StopReplica) WithBinlogReplicaController(controller binlogreplication.BinlogReplicaController) sql.Node {
	nc := *s
	nc.replicaController = controller
	return &nc
}

func (s *StopReplica) Resolved() bool {
	return true
}

func (s *StopReplica) String() string {
	return "STOP REPLICA"
}

func (s *StopReplica) Schema() sql.Schema {
	return nil
}

func (s *StopReplica) Children() []sql.Node {
	return nil
}

func (s *StopReplica) RowIter(ctx *sql.Context, _ sql.Row) (sql.RowIter, error) {
	if s.replicaController == nil {
		return nil, fmt.Errorf("no replication controller available")
	}

	err := s.replicaController.StopReplica(ctx)
	return sql.RowsToRowIter(), err
}

func (s *StopReplica) WithChildren(children ...sql.Node) (sql.Node, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(s, len(children), 0)
	}

	newNode := *s
	return &newNode, nil
}

func (s *StopReplica) CheckPrivileges(_ *sql.Context, _ sql.PrivilegedOperationChecker) bool {
	// TODO: implement privilege checks
	return true
}

// ResetReplica is a plan node for the "RESET REPLICA" statement.
// https://dev.mysql.com/doc/refman/8.0/en/reset-replica.html
type ResetReplica struct {
	replicaController binlogreplication.BinlogReplicaController
	All               bool
}

var _ sql.Node = (*ResetReplica)(nil)
var _ BinlogReplicaControllerCommand = (*ResetReplica)(nil)

func NewResetReplica(all bool) *ResetReplica {
	return &ResetReplica{
		All: all,
	}
}

// WithBinlogReplicaController implements the BinlogReplicaControllerCommand interface.
func (r *ResetReplica) WithBinlogReplicaController(controller binlogreplication.BinlogReplicaController) sql.Node {
	nc := *r
	nc.replicaController = controller
	return &nc
}

func (r *ResetReplica) Resolved() bool {
	return true
}

func (r *ResetReplica) String() string {
	sb := strings.Builder{}
	sb.WriteString("RESET REPLICA")
	if r.All {
		sb.WriteString(" ALL")
	}
	return sb.String()
}

func (r *ResetReplica) Schema() sql.Schema {
	return nil
}

func (r *ResetReplica) Children() []sql.Node {
	return nil
}

func (r *ResetReplica) RowIter(ctx *sql.Context, row sql.Row) (sql.RowIter, error) {
	if r.replicaController == nil {
		return nil, fmt.Errorf("no replication controller available")
	}

	err := r.replicaController.ResetReplica(ctx, r.All)
	return sql.RowsToRowIter(), err
}

func (r *ResetReplica) WithChildren(children ...sql.Node) (sql.Node, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(r, len(children), 0)
	}

	newNode := *r
	return &newNode, nil
}

func (r *ResetReplica) CheckPrivileges(ctx *sql.Context, opChecker sql.PrivilegedOperationChecker) bool {
	// TODO: implement privilege checks
	return true
}
