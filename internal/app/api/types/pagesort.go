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

func (ps *PagingNSorting) Init(page, size, ageRange, genderStr string) error {

	gender, err := genderInit(genderStr)
	if err != nil {
		return err
	}
	ps.Filter.Gender = gender

	time, err := convertAgeRangeToDate(ageRange)
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
func convertAgeRangeToDate(ageRange string) (*AgeRange, error) {
	// ageRange = "" range 0-100
	if ageRange == "" {
		return &AgeRange{
			Gte: time.Now().AddDate(-100, 0, 0),
			Lt:  time.Now(),
		}, nil
	}

	split := strings.Split(ageRange, ",")
	// age=24 range 24
	if len(split) == 1 {

		gteInt, err := strconv.ParseInt(split[0], 10, 64)

		if err != nil {
			return nil, err
		}
		return &AgeRange{
			Gte: time.Now().AddDate(-int(gteInt), 0, 0),
			Lt:  time.Now().AddDate(-int(gteInt+1), 0, 0),
		}, nil
	}
	// return err
	if len(split) != 2 {
		return nil, errors.Errorf("Can't convert ageRange to arr", ageRange)
	}

	gteInt, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return nil, err
	}

	ltInt, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return nil, err
	}
	//age=24,24 range 24
	if gteInt == ltInt {
		return &AgeRange{
			Gte: time.Now().AddDate(-int(gteInt), 0, 0),
			Lt:  time.Now().AddDate(-int(gteInt+1), 0, 0),
		}, nil
	}
	//age=24,25 range 24->25
	return &AgeRange{
		Gte: time.Now().AddDate(-int(gteInt), 0, 0),
		Lt:  time.Now().AddDate(-int(ltInt), 0, 0),
	}, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
