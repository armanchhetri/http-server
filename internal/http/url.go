package http

import "strings"

type URL struct {
	Raw string
	// Host        string
	Path        string
	QueryString string
}

type Query map[string]string

func (url *URL) Query() Query {
	queryMap := make(map[string]string)
	queries := strings.Split(url.QueryString, "&")
	for _, q := range queries {
		keyVals := strings.Split(q, "=")
		queryMap[keyVals[0]] = keyVals[1]
	}
	return Query(queryMap)
}
