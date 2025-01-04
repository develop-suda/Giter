package query

type GitHubQuery struct {
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
	Message       string `json:"message"`
	URL           string `json:"url"`
	CommittedDate string `json:"committedDate"`
	Oid           string `json:"oid"`
}
