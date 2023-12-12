package requests

import (
	"boilerplate/internal/domain"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func DecodePaginationQuery(r *http.Request) (domain.Pagination, error) {
	pageStr := r.URL.Query().Get("page")
	countStr := r.URL.Query().Get("count")

	p := domain.Pagination{
		Page:         1,
		CountPerPage: 20,
	}

	if pageStr != "" {
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			log.Print(err)
			return domain.Pagination{}, fmt.Errorf("problems in parsing 'page' query parameter")
		}

		p.Page = page
	}

	if countStr != "" {
		count, err := strconv.ParseUint(countStr, 10, 64)
		if err != nil {
			log.Print(err)
			return domain.Pagination{}, fmt.Errorf("problems in parsing 'count' query parameter")
		}

		p.CountPerPage = count
	}

	return p, nil
}
