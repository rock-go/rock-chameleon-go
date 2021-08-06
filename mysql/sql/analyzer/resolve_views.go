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

package analyzer

import (
	"fmt"

	"github.com/rock-go/rock-chameleon-go/mysql/sql"
	"github.com/rock-go/rock-chameleon-go/mysql/sql/plan"
)

func resolveViews(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope) (sql.Node, error) {
	span, _ := ctx.Span("resolve_views")
	defer span.Finish()

	return plan.TransformUp(n, func(n sql.Node) (sql.Node, error) {
		if n.Resolved() {
			return n, nil
		}

		t, ok := n.(*plan.UnresolvedTable)
		if !ok {
			return n, nil
		}

		name := t.Name()
		db := t.Database
		if db == "" {
			db = ctx.GetCurrentDatabase()
		}

		view, err := ctx.View(db, name)
		if err == nil {
			a.Log("view resolved: %q", name)

			// If this view is being asked for with an AS OF clause, then attempt to apply it to every table in the view.
			if t.AsOf != nil || t.Database != "" {
				a.Log("applying AS OF clause and database qualifier to view definition")

				// TODO: this direct editing of children is necessary because the view definition is declared as an opaque node,
				//  meaning that plan.TransformUp won't touch its children. It's only supposed to be touched by the
				//  resolve_subqueries function, which invokes the entire analyzer on the node. This is the only place we have
				//  to make this exception so far, but there may be others.
				children := view.Definition().Children()
				if len(children) == 1 {
					child, err := plan.TransformUp(children[0], func(n2 sql.Node) (sql.Node, error) {
						t2, ok := n2.(*plan.UnresolvedTable)
						if !ok {
							return n2, nil
						}

						if t.AsOf != nil {
							a.Log("applying AS OF clause to view " + t2.Name())
							if t2.AsOf != nil {
								return nil, sql.ErrIncompatibleAsOf.New(
									fmt.Sprintf("cannot combine AS OF clauses %s and %s",
										t.AsOf.String(), t2.AsOf.String()))
							}
							t2, _ = t2.WithAsOf(t.AsOf)
						}

						if t.Database != "" {
							a.Log("applying database clause to view " + t2.Name())
							if t2.Database == "" {
								t2, _ = t2.WithDatabase(db)
							}
						}

						return t2, nil
					})

					if err != nil {
						return nil, err
					}

					return view.Definition().WithChildren(child)
				}
			}

			return view.Definition(), nil
		}

		if sql.ErrNonExistingView.Is(err) {
			return n, nil
		}

		return nil, err
	})
}
