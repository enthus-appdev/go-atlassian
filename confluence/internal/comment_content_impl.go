package internal

import (
	"context"
	"fmt"
	model "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
	"github.com/ctreminiom/go-atlassian/v2/service"
	"github.com/ctreminiom/go-atlassian/v2/service/confluence"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// NewCommentService creates a new instance of CommentService.
// It takes a service.Connector as input and returns a pointer to CommentService.
func NewCommentService(client service.Connector) *CommentService {
	return &CommentService{
		internalClient: &internalCommentImpl{c: client},
	}
}

// CommentService provides methods to interact with comment operations in Confluence.
type CommentService struct {
	// internalClient is the connector interface for comment operations.
	internalClient confluence.CommentConnector
}

// Gets returns the comments on a piece of content.
//
// GET /wiki/rest/api/content/{id}/child/comment
//
// https://docs.go-atlassian.io/confluence-cloud/content/comments#get-content-comments
func (c *CommentService) Gets(ctx context.Context, contentID string, expand, location []string, startAt, maxResults int) (*model.ContentPageScheme, *model.ResponseScheme, error) {
	return c.internalClient.Gets(ctx, contentID, expand, location, startAt, maxResults)
}

type internalCommentImpl struct {
	c service.Connector
}

func (i *internalCommentImpl) Gets(ctx context.Context, contentID string, expand, location []string, startAt, maxResults int) (*model.ContentPageScheme, *model.ResponseScheme, error) {

	// Start tracing span
	tracer := otel.Tracer("github.com/ctreminiom/go-atlassian/v2/confluence")
	ctx, span := tracer.Start(ctx, "confluence.content.comment.gets")
	defer span.End()

	if contentID == "" {
		span.SetStatus(codes.Error, model.ErrNoContentID.Error())
		span.RecordError(model.ErrNoContentID)
		return nil, nil, model.ErrNoContentID
	}

	query := url.Values{}
	query.Add("start", strconv.Itoa(startAt))
	query.Add("limit", strconv.Itoa(maxResults))

	if len(expand) != 0 {
		query.Add("expand", strings.Join(expand, ","))
	}

	if len(location) != 0 {
		query.Add("location", strings.Join(location, ","))
	}

	endpoint := fmt.Sprintf("wiki/rest/api/content/%v/child/comment?%v", contentID, query.Encode())

	// Set span attributes for the HTTP request
	span.SetAttributes(
		attribute.String("http.method", http.MethodGet),
		attribute.String("http.url", endpoint),
		attribute.String("component", "go-atlassian"),
		attribute.String("module", "confluence"),
		attribute.String("operation", "content.comment.gets"),
		attribute.String("content.id", contentID),
		attribute.Int("pagination.start", startAt),
		attribute.Int("pagination.limit", maxResults),
	)

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, "", nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, nil, err
	}

	page := new(model.ContentPageScheme)
	response, err := i.c.Call(request, page)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, response, err
	}

	// Set response attributes
	if response != nil {
		span.SetAttributes(attribute.Int("http.status_code", response.Code))
		if response.Code >= 400 {
			span.SetStatus(codes.Error, http.StatusText(response.Code))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}

	return page, response, nil
}
