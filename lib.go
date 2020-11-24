package elastic

import (
	"context"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type CreateIndexBodyFunc func() (string, error)

func CreateIndexIfNotExists(ctx context.Context, client *Client, index string, bodyFn CreateIndexBodyFunc) (err error) {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return
	}
	if !exists {
		var body string
		body, err = bodyFn()
		if err != nil {
			return
		}
		err = createIndex(ctx, client, index, body)
	}
	return
}

func DeleteIndexIfExists(ctx context.Context, client *Client, index string) (err error) {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return
	}
	if !exists {
		return
	}
	_, err = client.DeleteIndex(index).Do(ctx)
	return
}

type SourceMapToFeature func(id string, sourceMap map[string]interface{}) (*geojson.Feature, error)

type EnvelopeQuery struct {
	name     string
	bound    orb.Bound
	relation string
}

func (q *EnvelopeQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})

	sq := make(map[string]interface{})
	source["geo_shape"] = sq

	params := make(map[string]interface{})
	sq[q.name] = params

	// using the Elasticsearch’s envelope GeoJSON extension
	// coordinates order: [UpperLeft, LowerRight]
	// https://www.elastic.co/guide/en/elasticsearch/reference/6.2/geo-shape.html#_envelope
	shape := make(map[string]interface{})
	shape["type"] = "envelope"
	shape["coordinates"] = []orb.Point{q.bound.LeftTop(), q.bound.RightBottom()}

	params["shape"] = shape
	params["relation"] = q.relation

	// source = map[string]interface{}{
	// 	"geo_shape": map[string]interface{}{
	// 		q.name: map[string]interface{}{
	// 			"shape": map[string]interface{}{
	// 				"type":        "envelope",
	// 				"coordinates": []orb.Point{q.bound.Min, q.bound.Max},
	// 			},
	// 			"relation": q.relation,
	// 		},
	// 	},
	// }

	return source, nil
}

func GeoIntersectionQuery(name string, bound orb.Bound) Query {
	return &EnvelopeQuery{
		name:     name,
		bound:    bound,
		relation: "intersects",
	}
}
