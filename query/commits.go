package query

import "time"

type CommitsQuery struct {
	User User `json:"user" graphql:"user(login: $USER_NAME)"`
}

type User struct {
	Repository Repository `json:"repository" graphql:"repository(name: $REPOSITORY_NAME)"`
}

type Repository struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Owner Owner  `json:"owner"`
	Refs  Refs   `json:"refs" graphql:"refs(refPrefix: \"refs/heads/\" first:10)"`
}

type Owner struct {
	Login string `json:"login"`
}

type Refs struct {
	Nodes []RefNode `json:"nodes"`
}

type RefNode struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}

type Target struct {
	Commit Commit `json:"commit" graphql:"... on Commit"`
}

type Commit struct {
	History History `json:"history" graphql:"history(first:100)"`
}

type History struct {
	Nodes []CommitNode `json:"nodes"`
}

type CommitNode struct {
	Message       string    `json:"message"`
	URL           string    `json:"url"`
	CommittedDate time.Time `json:"committedDate"`
	Oid           string    `json:"oid"`
}

func (c *CommitNode) CommittedDateJST() time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return c.CommittedDate.In(jst)
}

func (q *CommitsQuery) ToSimpleCommits() *SimpleCommits {
	repo := q.User.Repository
	sRepo := &SimpleCommits{
		Name:  repo.Name,
		URL:   repo.URL,
		Owner: repo.Owner,
	}

	for _, refNode := range repo.Refs.Nodes {
		branch := Branch{Name: refNode.Name}
		for _, commitNode := range refNode.Target.Commit.History.Nodes {
			branch.Nodes = append(branch.Nodes, commitNode)
		}
		sRepo.Branch = append(sRepo.Branch, branch)
	}

	return sRepo
}
