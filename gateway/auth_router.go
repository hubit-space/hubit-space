package pagination

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationParams struct {
	Collection          *mongo.Collection
	GlobalSearchFields  []string
	GlobalFilterAliases map[string]string
	OutputAliases       map[string]string
}

type PaginateRequest struct {
	Search      string
	Page        int
	PerPage     int
	Order       string
	Direction   string
	Filter      string
	Between     string
	FilterNotIn string
}

type PaginationResult struct {
	Status       int         `json:"status"`
	CurrentPage  int         `json:"current_page"`
	Data         interface{} `json:"data"`
	From         int         `json:"from"`
	To           int         `json:"to"`
	PerPage      int         `json:"per_page"`
	Total        int64       `json:"total"`
	TotalPages   int         `json:"total_pages"`
	PreviousPage int         `json:"prev,omitempty"`
	NextPage     int         `json:"next,omitempty"`
}

func Paginate(ctx *gin.Context, params PaginationParams) (*PaginationResult, error) {
	pageParam, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		pageParam = 1
	}
	perPageParam, err := strconv.Atoi(ctx.Query("per_page"))
	if err != nil {
		perPageParam = 10
	}

	searchParam := ctx.Query("search")
	orderParam := ctx.Query("order")
	directionParam := ctx.Query("direction")
	filterParam := ctx.Query("filter")
	betweenParam := ctx.Query("between")
	filterNotInParam := ctx.Query("filterNotIn")

	req := PaginateRequest{
		Page:        pageParam,
		PerPage:     perPageParam,
		Order:       orderParam,
		Direction:   directionParam,
		Filter:      filterParam,
		Between:     betweenParam,
		Search:      searchParam,
		FilterNotIn: filterNotInParam,
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	filter := bson.M{}
	if req.Search != "" {
		searchParts := []bson.M{}
		for _, field := range params.GlobalSearchFields {
			searchParts = append(searchParts, bson.M{
				field: bson.M{"$regex": req.Search, "$options": "i"},
			})
		}
		filter["$or"] = searchParts
	}

	if req.Filter != "" {
		for k, v := range parseKeyValuePairs(req.Filter, params.GlobalFilterAliases) {

			if intVal, err := strconv.Atoi(v); err == nil {
				filter[k] = intVal
				continue
			}

			lowerVal := strings.ToLower(v)
			if lowerVal == "true" || lowerVal == "false" {
				filter[k] = lowerVal == "true"
				continue
			}

			filter[k] = bson.M{"$regex": v, "$options": "i"}
		}
	}

	if req.Between != "" {
		for k, v := range parseKeyValuePairs(req.Between, params.GlobalFilterAliases) {
			rangeParts := strings.Split(v, ":")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid between format")
			}
			filter[k] = bson.M{"$gte": rangeParts[0], "$lte": rangeParts[1]}
		}
	}

	if req.FilterNotIn != "" {
		for k, v := range parseKeyValuePairs(req.FilterNotIn, params.GlobalFilterAliases) {
			lowerVal := strings.ToLower(v)
			switch lowerVal {
			case "true", "false":
				boolVal := lowerVal == "true"
				filter[k] = bson.M{"$ne": boolVal}
			default:
				filter[k] = bson.M{"$ne": v}
			}
		}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((req.Page - 1) * req.PerPage))
	findOptions.SetLimit(int64(req.PerPage))

	orderField := req.Order
	if orderField == "" {
		orderField = "_id"
	}

	direction := 1
	if strings.ToLower(req.Direction) == "desc" {
		direction = -1
	}

	findOptions.SetSort(bson.D{{Key: orderField, Value: direction}})

	log.Println("Filter:", filter)
	log.Println("FindOptions:", findOptions)

	cursor, err := params.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var data []bson.M
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	// Apply output aliases
	for _, item := range data {
		for from, to := range params.OutputAliases {
			if val, ok := item[from]; ok {
				item[to] = val
				delete(item, from)
			}
		}
	}

	total, err := params.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PerPage)))
	prev := req.Page - 1
	if prev < 1 {
		prev = 1
	}
	next := req.Page + 1
	if next > totalPages {
		next = totalPages
	}

	return &PaginationResult{
		Status:       200,
		CurrentPage:  req.Page,
		Data:         data,
		From:         (req.Page-1)*req.PerPage + 1,
		To:           (req.Page-1)*req.PerPage + len(data),
		PerPage:      req.PerPage,
		Total:        total,
		TotalPages:   totalPages,
		PreviousPage: prev,
		NextPage:     next,
	}, nil
}

func parseKeyValuePairs(input string, aliasMap map[string]string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(input, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key, value := strings.ToLower(kv[0]), kv[1]
		if alias, exists := aliasMap[key]; exists {
			result[alias] = value
		}
	}
	return result
}
