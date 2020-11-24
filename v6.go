// +build es6

package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	lib "github.com/olivere/elastic"
	"github.com/paulmach/orb/geojson"
)

type Error = lib.Error

type Client = lib.Client

type Query = lib.Query

type BoolQuery = lib.BoolQuery

type TermQuery = lib.TermQuery

type TermsQuery = lib.TermsQuery

type RangeQuery = lib.RangeQuery

type Script = lib.Script

type ScriptQuery = lib.ScriptQuery

type BulkIndexRequest = lib.BulkIndexRequest

type BulkUpdateRequest = lib.BulkUpdateRequest

type BulkDeleteRequest = lib.BulkDeleteRequest

type BulkResponse = lib.BulkResponse

type BulkResponseItem = lib.BulkResponseItem

type SearchResult = lib.SearchResult

func NewMatchAllQuery() *lib.MatchAllQuery {
	return lib.NewMatchAllQuery()
}

func NewBoolQuery() *lib.BoolQuery {
	return lib.NewBoolQuery()
}

func NewTermQuery(name string, value interface{}) *lib.TermQuery {
	return lib.NewTermQuery(name, value)
}

func NewTermsQuery(name string, values ...interface{}) *lib.TermsQuery {
	return lib.NewTermsQuery(name, values)
}

func NewRangeQuery(name string) *lib.RangeQuery {
	return lib.NewRangeQuery(name)
}

func NewScript(script string) *lib.Script {
	return lib.NewScript(script)
}

func NewScriptQuery(script *lib.Script) *lib.ScriptQuery {
	return lib.NewScriptQuery(script)
}

func NewBulkIndexRequest() *lib.BulkIndexRequest {
	return lib.NewBulkIndexRequest()
}

func NewBulkDeleteRequest() *lib.BulkDeleteRequest {
	return lib.NewBulkDeleteRequest()
}

func NewSearchSource() *lib.SearchSource {
	return lib.NewSearchSource()
}

func NewClient(httpClient *http.Client, urls ...string) (*Client, error) {
	retrier := lib.NewBackoffRetrier(lib.NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond))
	return lib.NewClient(
		lib.SetURL(urls...),
		lib.SetHttpClient(httpClient),
		lib.SetHealthcheck(false),
		lib.SetSniff(false),
		lib.SetRetrier(retrier),
		//lib.SetHealthcheckTimeout(3*time.Second),
		//lib.SetErrorLog(stdlog.New(os.Stderr, "", 0)),
	)
}

func createIndex(ctx context.Context, client *Client, index string, body string) (err error) {
	// NOTE: includeTypeName default to true for 6.x and false for 7.x and removed for 8.x
	_, err = client.CreateIndex(index).IncludeTypeName(false).BodyString(body).Do(ctx)
	return
}

func PutMapping(ctx context.Context, client *Client, index string, mapping map[string]interface{}) (err error) {
	resp, err := client.PutMapping().Index(index).IncludeTypeName(false).BodyJson(mapping).Do(ctx)
	if err != nil {
		return
	}
	if !resp.Acknowledged {
		err = errors.New("elastic: server acknowledged false")
		return
	}
	return
}

func GetMapping(ctx context.Context, client *Client, index string) (mapping map[string]interface{}, err error) {
	mapping, err = client.GetMapping().Index(index).IncludeTypeName(false).Do(ctx)
	return
}

func SearchHitToFeature(hit *lib.SearchHit, sourceMapToFeature SourceMapToFeature) (*geojson.Feature, error) {
	var sourceMap map[string]interface{}
	err := json.Unmarshal(*hit.Source, &sourceMap)
	if err != nil {
		return nil, err
	}
	return sourceMapToFeature(hit.Id, sourceMap)
}

func UnmarshalGetResult(result *lib.GetResult, v interface{}) (err error) {
	err = json.Unmarshal(*result.Source, v)
	return
}

func UnmarshalSearchHit(hit *lib.SearchHit, v interface{}) (err error) {
	err = json.Unmarshal(*hit.Source, v)
	return
}
