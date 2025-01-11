package query

type SCommits struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Owner  Owner    `json:"owner"`
	Branch []Branch `json:"branch"`
}

type Branch struct {
	Name  string       `json:"name"`
	Nodes []CommitNode `json:"nodes"`
}

func (r *SCommits) UpdateCommittedDatesToJST() *SCommits {
	for i, branch := range r.Branch {
		for j, commitNode := range branch.Nodes {
			r.Branch[i].Nodes[j].CommittedDate = commitNode.CommittedDateJST()
		}
	}
	return r
}
