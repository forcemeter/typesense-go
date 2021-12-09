package typesense

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/typesense/mocks"
)

func newMultiSearchParams() *api.MultiSearchParams {
	return &api.MultiSearchParams{
		Q:              pointer.String("text"),
		QueryBy:        pointer.String("company_name"),
		MaxHits:        pointer.Interface("all"),
		Prefix:         pointer.String("true"),
		FilterBy:       pointer.String("num_employees:=100"),
		SortBy:         pointer.String("num_employees:desc"),
		FacetBy:        pointer.String("year_started"),
		MaxFacetValues: pointer.Int(10),
		FacetQuery:     pointer.String("facetQuery"),
		NumTypos:       pointer.Int(2),
		Page:           pointer.Int(1),
		PerPage:        pointer.Int(10),
		GroupBy:        pointer.String("country"),
		GroupLimit:     pointer.Int(3),
		IncludeFields:  pointer.String("company_name"),
	}
}

func newMultiSearchBodyParams() api.MultiSearchSearchesParameter {
	return api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: "companies",
				MultiSearchParameters: api.MultiSearchParameters{
					Q:       pointer.String("text"),
					QueryBy: pointer.String("company_name"),
				},
			},
			{
				Collection: "companies",
				MultiSearchParameters: api.MultiSearchParameters{
					Q:       pointer.String("text"),
					QueryBy: pointer.String("company_name"),
				},
			},
		},
	}
}

func newMultiSearchResult() *api.MultiSearchResult {
	return &api.MultiSearchResult{
		Results: []api.SearchResult{
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]int{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]int{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
		},
	}
}

func TestMultiSearchResultDeserialization(t *testing.T) {
	inputJSON := `{
			"results": [
				{
					"facet_counts": [],
					"found": 1,
					"search_time_ms": 1,
					"hits": [
						{
						"highlights": [
							{
								"field": "company_name",
								"snippet": "<mark>Stark</mark> Industries",
								"matched_tokens": ["Stark"]
							}
						],
						"document": {
								"id": "124",
								"company_name": "Stark Industries",
								"num_employees": 5215,
								"country": "USA"
							}
						}
					]
				},
				{
					"facet_counts": [],
					"found": 1,
					"search_time_ms": 1,
					"hits": [
						{
						"highlights": [
							{
								"field": "company_name",
								"snippet": "<mark>Stark</mark> Industries",
								"matched_tokens": ["Stark"]
							}
						],
						"document": {
								"id": "124",
								"company_name": "Stark Industries",
								"num_employees": 5215,
								"country": "USA"
							}
						}
					]
				}
		]
	}`
	expected := &api.MultiSearchResult{
		Results: []api.SearchResult{
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]int{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]int{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
		},
	}
	result := &api.MultiSearchResult{}
	err := json.Unmarshal([]byte(inputJSON), result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMultiSearch(t *testing.T) {
	expectedParams := newMultiSearchParams()
	expectedResult := newMultiSearchResult()
	expectedBody := newMultiSearchBodyParams()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := newMultiSearchResult()

	mockAPIClient.EXPECT().
		MultiSearchWithResponse(gomock.Not(gomock.Nil()), expectedParams, api.MultiSearchJSONRequestBody(expectedBody)).Return(&api.MultiSearchResponse{
		JSON200: mockedResult,
	}, nil).Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newMultiSearchParams()
	body := newMultiSearchBodyParams()
	result, err := client.MultiSearch.Perform(params, body)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestMultiSearchOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	expectedParams := newMultiSearchParams()
	expectedBody := newMultiSearchBodyParams()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		MultiSearchWithResponse(gomock.Not(gomock.Nil()), expectedParams, api.MultiSearchJSONRequestBody(expectedBody)).
		Return(&api.MultiSearchResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newMultiSearchParams()
	_, err := client.MultiSearch.Perform(params, newMultiSearchBodyParams())
	assert.NotNil(t, err)
}

func TestMultiSearchOnApiClientError(t *testing.T) {
	expectedParams := newMultiSearchParams()
	expectedBody := newMultiSearchBodyParams()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		MultiSearchWithResponse(gomock.Not(gomock.Nil()), expectedParams, api.MultiSearchJSONRequestBody(expectedBody)).
		Return(nil, errors.New("failed request")).Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newMultiSearchParams()
	_, err := client.MultiSearch.Perform(params, newMultiSearchBodyParams())
	assert.NotNil(t, err)
}

