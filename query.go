package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Feature struct {
	Geometry struct {
		Coordinates [2]float64 `json:"coordinates"`
	} `json:"geometry"`

	Properties struct {
		EventType             *string    `json:"eventtype"`
		PublicID              *string    `json:"publicid"`
		ModificationTime      *time.Time `json:"modificationtime"`
		OriginTime            *time.Time `json:"origintime"`
		OriginError           *float64   `json:"originerror"`
		EarthModel            *string    `json:"earthmodel"`
		EvaluationMethod      *string    `json:"evaluationmethod"`
		EvaluationStatus      *string    `json:"evaluationstatus"`
		EvaluationMode        *string    `json:"evaluationmode"`
		Latitude              *float64   `json:"latitude"`
		Longitude             *float64   `json:"longitude"`
		Depth                 *float64   `json:"depth"`
		DepthType             *string    `json:"depthtype"`
		UsedPhaseCount        *int32     `json:"usedphasecount"`
		UsedStationCount      *int32     `json:"usedstationcount"`
		AzimuthalGap          *float64   `json:"azimuthalgap"`
		MinimumDistance       *float64   `json:"minimumdistance"`
		Magnitude             *float64   `json:"magnitude"`
		MagnitudeType         *string    `json:"magnitudetype"`
		MagnitudeStationCount *int32     `json:"magnitudestationcount"`
		MagnitudeUncertainty  *float64   `json:"magnitudeuncertainty"`
	} `json:"properties"`
}

type Search struct {
	Features []Feature `json:"features"`
}

type Query struct {
	service string
	limit   int
	filters []string
}

func NewQuery(service string, limit int) *Query {
	return &Query{
		service: service,
		limit:   limit,
	}
}

func (q *Query) Values() *url.Values {

	v := url.Values{}
	v.Set("service", "WFS")
	v.Set("version", "1.0.0")
	v.Set("request", "GetFeature")
	v.Set("typeName", "geonet:quake_search_v1")
	if q.limit > 0 {
		v.Set("maxFeatures", strconv.Itoa(q.limit))
	}
	v.Set("outputFormat", "json")

	return &v
}

func (q *Query) AddFilter(filter string) {
	q.filters = append(q.filters, filter)
}

func (q Query) URL() *url.URL {

	v := q.Values()

	var cql_filter string
	if len(q.filters) > 0 {
		cql_filter = "&cql_filter=" + strings.Join(q.filters, "+and+")
	}

	u := url.URL{
		Scheme:   "http",
		Host:     q.service,
		Path:     "geonet/ows",
		RawQuery: v.Encode() + cql_filter,
	}

	return &u
}

func (q Query) Search() (*Search, error) {

	res, err := http.Get(q.URL().String())
	if err != nil {
		return nil, err
	}

	var body []byte
	if res.Body != nil {
		defer res.Body.Close()
		if body, err = ioutil.ReadAll(res.Body); err != nil {
			return nil, err
		}
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error [%d]: %s", res.StatusCode, string(body))
	}

	search := Search{}
	err = json.Unmarshal(body, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}
