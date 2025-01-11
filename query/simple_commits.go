package query

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
