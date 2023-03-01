package domain

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	MaxSearchLimit     = 100
	DefaultSearchLimit = 20
)

// Filter is a struct for filtering
type Filter struct {
	// Page is a number of page
	// in:query
	Page int `json:"page"`
	// Limit is a count of values
	// in:query
	List int `json:"list"`
	// SortBy
	// in:query
	SortBy []string `json:"sortBy"` // firstname asc, lastname dsc
}

// Pagination is a struct for pagination
type Pagination struct {
	// Order can be asc or dsc. asc by default
	Order string `json:"order"`
	// Page is a number of page
	Page int `json:"page"`
	// Limit is a count of values
	List int `json:"list"`
	// TotalItems is a number of items
	Total *int64 `json:"total"`
}

func (f *Filter) OrderString() string {
	return strings.Join(f.SortBy, ",")
}

func (f *Filter) Validate() {
	f.SortBy = f.validatedSortBy()
	if f.Page < 0 {
		f.Page = 0
	}
	if f.List < 1 || f.List > MaxSearchLimit {
		f.List = DefaultSearchLimit
	}
}

func (f *Filter) validatedSortBy() []string {
	validatedSortByArray := make([]string, 0, len(f.SortBy))
	for _, sortBy := range f.SortBy {
		tokens := strings.Split(sortBy, " ")
		if len(tokens) != 1 && len(tokens) != 2 {
			continue
		}
		validatedSortBy := tokens[0]
		// update json model fields to database fields
		switch strings.ToLower(validatedSortBy) {
		case "createdat":
			validatedSortBy = "created_at"
		case "updatedat":
			validatedSortBy = "updated_at"
		case "paymentstatus":
			validatedSortBy = "payment_status"
		case "fullname":
			validatedSortBy = "full_name"
		}
		if len(tokens) == 2 {
			switch strings.ToLower(tokens[1]) {
			case "asc":
				validatedSortBy += " asc"
			case "desc":
				validatedSortBy += " desc"
			default:
				continue
			}
		}
		validatedSortByArray = append(validatedSortByArray, validatedSortBy)
	}
	return validatedSortByArray
}

func GetFilterFromQuery(r *http.Request) (*Filter, error) {
	var (
		page int
		list int
		err  error
	)
	params := r.URL.Query()
	if len(params["page"]) != 0 {
		page, err = strconv.Atoi(params["page"][0])
		if err != nil {
			return nil, errors.New("cannot parse page query param")
		}
	}
	if len(params["list"]) != 0 {
		list, err = strconv.Atoi(params["list"][0])
		if err != nil {
			return nil, errors.New("cannot parse list query param")
		}
	}

	filter := &Filter{
		Page:   page,
		List:   list,
		SortBy: params["sortBy"],
	}
	filter.Validate()

	return filter, nil
}
