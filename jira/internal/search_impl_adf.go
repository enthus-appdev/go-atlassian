package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	model "github.com/enthus-appdev/go-atlassian/pkg/infra/models"
	"github.com/enthus-appdev/go-atlassian/service"
	"github.com/enthus-appdev/go-atlassian/service/jira"
)

// SearchADFService provides methods to manage advanced document format (ADF) searches in Jira Service Management.
type SearchADFService struct {
	// internalClient is the connector interface for ADF search operations.
	internalClient jira.SearchADFConnector
}

// Checks checks whether one or more issues would be returned by one or more JQL queries.
//
// POST /rest/api/{2-3}/jql/match
//
// https://docs.go-atlassian.io/jira-software-cloud/issues/search#check-issues-against-jql
func (s *SearchADFService) Checks(ctx context.Context, payload *model.IssueSearchCheckPayloadScheme) (*model.IssueMatchesPageScheme, *model.ResponseScheme, error) {
	return s.internalClient.Checks(ctx, payload)
}

// Get search issues using JQL query under the HTTP Method GET
//
// GET /rest/api/3/search
//
// https://docs.go-atlassian.io/jira-software-cloud/issues/search#search-for-issues-using-jql-get
func (s *SearchADFService) Get(ctx context.Context, jql string, fields, expands []string, startAt, maxResults int, validate string) (*model.IssueSearchScheme, *model.ResponseScheme, error) {
	return s.internalClient.Get(ctx, jql, fields, expands, startAt, maxResults, validate)
}

// Post search issues using JQL query under the HTTP Method POST
//
// POST /rest/api/3/search
//
// https://docs.go-atlassian.io/jira-software-cloud/issues/search#search-for-issues-using-jql-get
func (s *SearchADFService) Post(ctx context.Context, jql string, fields, expands []string, startAt, maxResults int, validate string) (*model.IssueSearchScheme, *model.ResponseScheme, error) {
	return s.internalClient.Post(ctx, jql, fields, expands, startAt, maxResults, validate)
}

type internalSearchADFImpl struct {
	c       service.Connector
	version string
}

func (i *internalSearchADFImpl) Checks(ctx context.Context, payload *model.IssueSearchCheckPayloadScheme) (*model.IssueMatchesPageScheme, *model.ResponseScheme, error) {

	endpoint := fmt.Sprintf("rest/api/%v/jql/match", i.version)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, "", payload)
	if err != nil {
		return nil, nil, err
	}

	issues := new(model.IssueMatchesPageScheme)
	response, err := i.c.Call(request, issues)
	if err != nil {
		return nil, response, err
	}

	return issues, response, nil
}

func (i *internalSearchADFImpl) Get(ctx context.Context, jql string, nextPageToken *string, maxResults int, fields, expands, properties []string, fieldsByKey, failFast bool, reconcileIssues []int) (*model.IssueSearchScheme, *model.ResponseScheme, error) {

	if jql == "" {
		return nil, nil, model.ErrNoJQL
	}

	params := url.Values{}
	params.Add("jql", jql)
	if nextPageToken != nil {
		params.Add("nextPageToken", *nextPageToken)
	}
	params.Add("maxResults", strconv.Itoa(maxResults))
	params.Add("fieldsByKey", strconv.FormatBool(fieldsByKey))
	params.Add("failFast", strconv.FormatBool(failFast))

	if len(expands) != 0 {
		params.Add("expand", strings.Join(expands, ","))
	}

	if len(fields) != 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	if len(properties) != 0 {
		params.Add("properties", strings.Join(properties, ","))
	}

	if len(reconcileIssues) != 0 {
		params.Add("reconcileIssues", strings.Join(strings.Fields(fmt.Sprint(reconcileIssues)), ","))
	}

	endpoint := fmt.Sprintf("rest/api/%v/search/jql?%v", i.version, params.Encode())

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, "", nil)
	if err != nil {
		return nil, nil, err
	}

	issues := new(model.IssueSearchScheme)
	response, err := i.c.Call(request, issues)
	if err != nil {
		return nil, response, err
	}

	return issues, response, nil
}

func (i *internalSearchADFImpl) Post(ctx context.Context, jql string, nextPageToken *string, maxResults int, fields, expands, properties []string, fieldsByKey, failFast bool, reconcileIssues []int) (*model.IssueSearchScheme, *model.ResponseScheme, error) {

	payload := struct {
		Jql             string   `json:"jql,omitempty"`
		NextPageToken   *string  `json:"nextPageToken,omitempty"`
		MaxResults      int      `json:"maxResults,omitempty"`
		Fields          []string `json:"fields,omitempty"`
		Expand          []string `json:"expand,omitempty"`
		Properties      []string `json:"properties,omitempty"`
		FieldsByKey     bool     `json:"fieldsByKey,omitempty"`
		FailFast        bool     `json:"failFast,omitempty"`
		ReconcileIssues []int    `json:"reconcileIssues,omitempty"`
	}{
		Jql:             jql,
		NextPageToken:   nextPageToken,
		MaxResults:      maxResults,
		Fields:          fields,
		Expand:          expands,
		Properties:      properties,
		FieldsByKey:     fieldsByKey,
		FailFast:        failFast,
		ReconcileIssues: reconcileIssues,
	}

	endpoint := fmt.Sprintf("rest/api/%v/search/jql", i.version)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, "", payload)
	if err != nil {
		return nil, nil, err
	}

	issues := new(model.IssueSearchScheme)
	response, err := i.c.Call(request, issues)
	if err != nil {
		return nil, response, err
	}

	return issues, response, nil
}
