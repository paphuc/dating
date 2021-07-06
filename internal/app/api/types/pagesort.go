// Pagination and Sorting
package types

import "strconv"

type PagingNSorting struct {
	Size int `json:"size" default:"100"`
	Page int `json:"page" default:"1"`
}

func (ps *PagingNSorting) Init(page, size string) error {

	if size == "" && page == "" {
		ps.Size = 100
		ps.Page = 1
		return nil
	}

	sizeInt, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return err
	}

	if page == "" {
		page = "1" // default page = 1
	}

	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return err
	}

	ps.Page = int(sizeInt)
	ps.Size = int(pageInt)

	return nil
}
