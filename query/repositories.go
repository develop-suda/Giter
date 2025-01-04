// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    welcome, err := UnmarshalWelcome(bytes)
//    bytes, err = welcome.Marshal()

package query

import (
	"encoding/json"
	"time"
)

type Welcome []WelcomeElement

func UnmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Welcome) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type WelcomeElement struct {
	ID                       int64         `json:"id"`
	NodeID                   string        `json:"node_id"`
	Name                     string        `json:"name"`
	FullName                 string        `json:"full_name"`
	Private                  bool          `json:"private"`
	Owner                    Owner         `json:"owner"`
	HTMLURL                  string        `json:"html_url"`
	Description              interface{}   `json:"description"`
	Fork                     bool          `json:"fork"`
	URL                      string        `json:"url"`
	ForksURL                 string        `json:"forks_url"`
	KeysURL                  string        `json:"keys_url"`
	CollaboratorsURL         string        `json:"collaborators_url"`
	TeamsURL                 string        `json:"teams_url"`
	HooksURL                 string        `json:"hooks_url"`
	IssueEventsURL           string        `json:"issue_events_url"`
	EventsURL                string        `json:"events_url"`
	AssigneesURL             string        `json:"assignees_url"`
	BranchesURL              string        `json:"branches_url"`
	TagsURL                  string        `json:"tags_url"`
	BlobsURL                 string        `json:"blobs_url"`
	GitTagsURL               string        `json:"git_tags_url"`
	GitRefsURL               string        `json:"git_refs_url"`
	TreesURL                 string        `json:"trees_url"`
	StatusesURL              string        `json:"statuses_url"`
	LanguagesURL             string        `json:"languages_url"`
	StargazersURL            string        `json:"stargazers_url"`
	ContributorsURL          string        `json:"contributors_url"`
	SubscribersURL           string        `json:"subscribers_url"`
	SubscriptionURL          string        `json:"subscription_url"`
	CommitsURL               string        `json:"commits_url"`
	GitCommitsURL            string        `json:"git_commits_url"`
	CommentsURL              string        `json:"comments_url"`
	IssueCommentURL          string        `json:"issue_comment_url"`
	ContentsURL              string        `json:"contents_url"`
	CompareURL               string        `json:"compare_url"`
	MergesURL                string        `json:"merges_url"`
	ArchiveURL               string        `json:"archive_url"`
	DownloadsURL             string        `json:"downloads_url"`
	IssuesURL                string        `json:"issues_url"`
	PullsURL                 string        `json:"pulls_url"`
	MilestonesURL            string        `json:"milestones_url"`
	NotificationsURL         string        `json:"notifications_url"`
	LabelsURL                string        `json:"labels_url"`
	ReleasesURL              string        `json:"releases_url"`
	DeploymentsURL           string        `json:"deployments_url"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
	PushedAt                 time.Time     `json:"pushed_at"`
	GitURL                   string        `json:"git_url"`
	SSHURL                   string        `json:"ssh_url"`
	CloneURL                 string        `json:"clone_url"`
	SvnURL                   string        `json:"svn_url"`
	Homepage                 interface{}   `json:"homepage"`
	Size                     int64         `json:"size"`
	StargazersCount          int64         `json:"stargazers_count"`
	WatchersCount            int64         `json:"watchers_count"`
	Language                 *string       `json:"language"`
	HasIssues                bool          `json:"has_issues"`
	HasProjects              bool          `json:"has_projects"`
	HasDownloads             bool          `json:"has_downloads"`
	HasWiki                  bool          `json:"has_wiki"`
	HasPages                 bool          `json:"has_pages"`
	HasDiscussions           bool          `json:"has_discussions"`
	ForksCount               int64         `json:"forks_count"`
	MirrorURL                interface{}   `json:"mirror_url"`
	Archived                 bool          `json:"archived"`
	Disabled                 bool          `json:"disabled"`
	OpenIssuesCount          int64         `json:"open_issues_count"`
	License                  interface{}   `json:"license"`
	AllowForking             bool          `json:"allow_forking"`
	IsTemplate               bool          `json:"is_template"`
	WebCommitSignoffRequired bool          `json:"web_commit_signoff_required"`
	Topics                   []interface{} `json:"topics"`
	Visibility               Visibility    `json:"visibility"`
	Forks                    int64         `json:"forks"`
	OpenIssues               int64         `json:"open_issues"`
	Watchers                 int64         `json:"watchers"`
	DefaultBranch            DefaultBranch `json:"default_branch"`
}

type ROwner struct {
	Login             Login        `json:"login"`
	ID                int64        `json:"id"`
	NodeID            NodeID       `json:"node_id"`
	AvatarURL         string       `json:"avatar_url"`
	GravatarID        string       `json:"gravatar_id"`
	URL               string       `json:"url"`
	HTMLURL           string       `json:"html_url"`
	FollowersURL      string       `json:"followers_url"`
	FollowingURL      FollowingURL `json:"following_url"`
	GistsURL          GistsURL     `json:"gists_url"`
	StarredURL        StarredURL   `json:"starred_url"`
	SubscriptionsURL  string       `json:"subscriptions_url"`
	OrganizationsURL  string       `json:"organizations_url"`
	ReposURL          string       `json:"repos_url"`
	EventsURL         EventsURL    `json:"events_url"`
	ReceivedEventsURL string       `json:"received_events_url"`
	Type              Type         `json:"type"`
	UserViewType      Visibility   `json:"user_view_type"`
	SiteAdmin         bool         `json:"site_admin"`
}

type DefaultBranch string

const (
	Main DefaultBranch = "main"
)

type EventsURL string

const (
	HTTPSAPIGithubCOMUsersDevelopSudaEventsPrivacy EventsURL = "https://api.github.com/users/develop-suda/events{/privacy}"
)

type FollowingURL string

const (
	HTTPSAPIGithubCOMUsersDevelopSudaFollowingOtherUser FollowingURL = "https://api.github.com/users/develop-suda/following{/other_user}"
)

type GistsURL string

const (
	HTTPSAPIGithubCOMUsersDevelopSudaGistsGistID GistsURL = "https://api.github.com/users/develop-suda/gists{/gist_id}"
)

type Login string

const (
	DevelopSuda Login = "develop-suda"
)

type NodeID string

const (
	MDQ6VXNlcjcxMzcwNzA5 NodeID = "MDQ6VXNlcjcxMzcwNzA5"
)

type StarredURL string

const (
	HTTPSAPIGithubCOMUsersDevelopSudaStarredOwnerRepo StarredURL = "https://api.github.com/users/develop-suda/starred{/owner}{/repo}"
)

type Type string

const (
	RUser Type = "User"
)

type Visibility string

const (
	Public Visibility = "public"
)
