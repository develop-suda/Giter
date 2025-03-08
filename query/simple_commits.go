package query

import (
	"sort"
	"time"
)

type SimpleCommits struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Owner  Owner    `json:"owner"`
	Branch []Branch `json:"branch"`
}

type Branch struct {
	Name  string       `json:"name"`
	Nodes []CommitNode `json:"nodes"`
}

func (r *SimpleCommits) UpdateCommittedDatesToJST() *SimpleCommits {
	for i, branch := range r.Branch {
		for j, commitNode := range branch.Nodes {
			r.Branch[i].Nodes[j].CommittedDate = commitNode.CommittedDateJST()
		}
	}
	return r
}

type Commits struct {
	Name          string    `json:"name"`
	Branch        string    `json:"branch"`
	Message       string    `json:"message"`
	URL           string    `json:"url"`
	CommittedDate time.Time `json:"committedDate"`
}

// TODO:同一コミットを削除する
// main,masterブランチを正とする
// oidで検索をかけ一致したら削除する
func (c *SimpleCommits) RemoveDuplicateCommits() {
	primaryBranches := map[string]bool{"main": true, "master": true}
	commitMap := make(map[string]bool)

	for _, branch := range c.Branch {
		if primaryBranches[branch.Name] {
			for _, commitNode := range branch.Nodes {
				commitMap[commitNode.Oid] = true
			}
		}
	}

	for i, branch := range c.Branch {
		if !primaryBranches[branch.Name] {
			var uniqueNodes []CommitNode
			for _, commitNode := range branch.Nodes {
				// HACK:こーゆー書き方もある
				if !commitMap[commitNode.Oid] {
					uniqueNodes = append(uniqueNodes, commitNode)
					commitMap[commitNode.Oid] = true
				}
			}
			c.Branch[i].Nodes = uniqueNodes
		}
	}
}

func ToCommits(commits *[]SimpleCommits) *[]Commits {
	var result []Commits
	for _, simpleCommit := range *commits {
		for _, branch := range simpleCommit.Branch {
			for _, commitNode := range branch.Nodes {
				result = append(result, Commits{
					Name:          simpleCommit.Name,
					Branch:        branch.Name,
					Message:       commitNode.Message,
					URL:           commitNode.URL,
					CommittedDate: commitNode.CommittedDate,
				})
			}
		}
	}
	SortCommits(result)
	return &result
}

func SortCommits(commits []Commits) *[]Commits {
	sort.SliceStable(commits, func(i, j int) bool { return (commits)[i].CommittedDate.After((commits)[j].CommittedDate) })
	return &commits
}
