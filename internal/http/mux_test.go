package http

import (
	"testing"
)

func TestInsert(t *testing.T) {
	type args struct {
		node *TrieNode
		path []string
	}
	tests := []struct {
		name string
		args args
		want *TrieNode
	}{
		{
			name: "builds a trie and returns leaf node",
			args: args{
				&TrieNode{},
				[]string{"path", "a", "b", "<param>"},
			},
			want: &TrieNode{lastNode: true, pathParam: "param"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Insert(tt.args.node, tt.args.path); got.lastNode == false || (got.pathParam != tt.want.pathParam) {
				t.Errorf("Insert() = %v, want pathParam %v and lastNode %v", got, tt.want.pathParam, tt.want.lastNode)
			}
		})
	}
}
