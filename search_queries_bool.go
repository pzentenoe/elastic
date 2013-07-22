// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A bool query matches documents matching boolean
// combinations of other queries.
// For more details, see:
// http://www.elasticsearch.org/guide/reference/query-dsl/bool-query.html
type BoolQuery struct {
	Query
	mustClauses        []Query
	shouldClauses      []Query
	mustNotClauses     []Query
	boost              *float32
	disableCoord       *bool
	minimumShouldMatch string
}

// Creates a new bool query.
func NewBoolQuery() BoolQuery {
	q := BoolQuery{
		mustClauses:    make([]Query, 0),
		shouldClauses:  make([]Query, 0),
		mustNotClauses: make([]Query, 0),
	}
	return q
}

func (q BoolQuery) Must(query Query) BoolQuery {
	q.mustClauses = append(q.mustClauses, query)
	return q
}

func (q BoolQuery) MustNot(query Query) BoolQuery {
	q.mustNotClauses = append(q.mustNotClauses, query)
	return q
}

func (q BoolQuery) Should(query Query) BoolQuery {
	q.shouldClauses = append(q.shouldClauses, query)
	return q
}

func (q BoolQuery) Boost(boost float32) BoolQuery {
	q.boost = &boost
	return q
}

func (q BoolQuery) DisableCoord(disableCoord bool) BoolQuery {
	q.disableCoord = &disableCoord
	return q
}

func (q BoolQuery) MinimumShouldMatch(minimumShouldMatch string) BoolQuery {
	q.minimumShouldMatch = minimumShouldMatch
	return q
}

// Creates the query source for the bool query.
func (q BoolQuery) Source() interface{} {
	// {
	//	"bool" : {
	//		"must" : {
	//			"term" : { "user" : "kimchy" }
	//		},
	//		"must_not" : {
	//			"range" : {
	//				"age" : { "from" : 10, "to" : 20 }
	//			}
	//		},
	//		"should" : [
	//			{
	//				"term" : { "tag" : "wow" }
	//			},
	//			{
	//				"term" : { "tag" : "elasticsearch" }
	//			}
	//		],
	//		"minimum_number_should_match" : 1,
	//		"boost" : 1.0
	//	}
	// }

	query := make(map[string]interface{})

	boolClause := make(map[string]interface{})
	query["bool"] = boolClause

	// must
	if len(q.mustClauses) == 1 {
		boolClause["must"] = q.mustClauses[0].Source()
	} else if len(q.mustClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range q.mustClauses {
			clauses = append(clauses, subQuery.Source())
		}
		boolClause["must"] = clauses
	}

	// must_not
	if len(q.mustNotClauses) == 1 {
		boolClause["must_not"] = q.mustNotClauses[0].Source()
	} else if len(q.mustNotClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range q.mustNotClauses {
			clauses = append(clauses, subQuery.Source())
		}
		boolClause["must_not"] = clauses
	}

	// should
	if len(q.shouldClauses) == 1 {
		boolClause["should"] = q.shouldClauses[0].Source()
	} else if len(q.shouldClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range q.shouldClauses {
			clauses = append(clauses, subQuery.Source())
		}
		boolClause["should"] = clauses
	}

	return query
}