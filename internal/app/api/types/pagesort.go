// Pagination and Sorting
package types

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type PagingNSorting struct {
	Size   int    `json:"size" default:"100"`
	Page   int    `json:"page" default:"1"`
	Filter Filter `json:"filter"`
}
type Filter struct {
	AgeRange AgeRange `json:"age"`
	Gender   []string `json:"gender" default:"" bson:"gender,omitempty"`
}
type AgeRange struct {
	Gte time.Time `json:"gte"`
	Lt  time.Time `json:"lt"`
}
type Pagination struct {
	TotalItems      int    `json:"totalItems"`
	TotalPages      int    `json:"totalPages"`
	CurrentPage     int    `json:"currentPage"`
	MaxItemsPerPage int    `json:"maxItemsPerPage"`
	Filter          Filter `json:"filter"`
}

func (ps *PagingNSorting) Init(page, size, minAge, maxAge, genderStr string) error {

	gender, err := genderInit(genderStr)
	if err != nil {
		return err
	}
	ps.Filter.Gender = gender

	time, err := convertAgeRangeToDate(minAge, maxAge)
	if err != nil {
		return err
	}
	ps.Filter.AgeRange = *time

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

	ps.Page = int(pageInt)
	ps.Size = int(sizeInt)

	return nil
}

func genderInit(gender string) ([]string, error) {
	genderArray := []string{"Male", "Female", "Other"}

	if gender != "" {
		split := strings.Split(gender, ",")
		if len(split) != 1 {
			for _, v := range split {
				if !stringInSlice(v, genderArray) {
					return nil, errors.Errorf("gender not in arr {Male, Female, Other}", gender)
				}
			}
			return split, nil
		}

		return []string{gender}, nil
	}

	return genderArray, nil
}

// convertAgeRangeToDate, ex: 18 year olds (now 2021) -> year birdday 2003 ->
// range birdday 2003 2002
func convertAgeRangeToDate(minAgeStr, maxAgeStr string) (*AgeRange, error) {
	// range 0-100
	min, err := convertAgeStr(minAgeStr, 0)
	if err != nil {
		return nil, err
	}

	max, err := convertAgeStr(maxAgeStr, 100)
	if err != nil {
		return nil, err
	}

	return &AgeRange{
		Gte: time.Now().AddDate(-int(max+1), 0, 0),
		Lt:  time.Now().AddDate(-int(min), 0, 0),
	}, nil
}

func convertAgeStr(ageStr string, defaultValue int) (int, error) {
	if ageStr == "" {
		return defaultValue, nil
	} else {
		value, err := strconv.ParseInt(ageStr, 10, 64)
		if err != nil {
			return -1, err
		}
		return int(value), nil
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
