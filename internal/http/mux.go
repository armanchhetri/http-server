package http

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

type TrieNode struct {
	children    []*TrieNode
	handlerFunc HandlerFunc
	word        string
	lastNode    bool
	pathParam   string
}

func Insert(node *TrieNode, path []string) *TrieNode {
	if len(path) == 0 {
		node.lastNode = true
		return node
	}
	word := path[0]
	if strings.HasPrefix(word, "<") && strings.HasSuffix(word, ">") {
		node.pathParam = string(word[1 : len(word)-1])
		return Insert(node, path[1:])
	}

	for _, child := range node.children {
		if child.word == word {
			return Insert(child, path[1:])
		}
	}
	tn := TrieNode{word: word}
	node.children = append(node.children, &tn)
	return Insert(&tn, path[1:])
}

func Search(node *TrieNode, path []string, paramAcc pathParams) *TrieNode {
	if len(path) == 0 && node.lastNode {
		return node
	} else if len(path) == 0 {
		return nil
	}
	word := path[0]
	// only add parameter if not already there
	_, ok := paramAcc[node.pathParam]
	if !ok && node.pathParam != "" {
		paramAcc[node.pathParam] = word
		return Search(node, path[1:], paramAcc)

	}
	for _, child := range node.children {
		if child.word == word {
			return Search(child, path[1:], paramAcc)
		}
	}
	return nil
}

type HandlerFunc func(ResponseWriter, *Request)

type pathParams map[string]string

type Mux struct {
	handlerTree *TrieNode
}

func (m *Mux) ServeHTTP(rw ResponseWriter, r *Request) {
	handlerFunc, param := m.getHandler(r.URL.Path)
	if handlerFunc == nil {
		rw.WriteStatus(StatusNotFound)
		rw.Write([]byte{})
		return
	}
	r.PathParam = param
	handlerFunc(rw, r)
}

func (m *Mux) getHandler(path string) (HandlerFunc, pathParams) {
	pathWords := strings.Split(path, "/")
	pathParam := make(pathParams)
	endNode := Search(m.handlerTree, pathWords, pathParam)
	if endNode == nil {
		return nil, pathParam
	}
	return endNode.handlerFunc, pathParam
}

func (m *Mux) Register(endPath string, handler HandlerFunc) {
	log.Infof("Registering route %s\n", endPath)
	urlWords := strings.Split(endPath, "/")

	endNode := Insert(m.handlerTree, urlWords)
	endNode.handlerFunc = handler
}

func NewMux() *Mux {
	handlerTree := &TrieNode{}
	return &Mux{handlerTree: handlerTree}
}
