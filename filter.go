package main

import (
	"encoding/json"
	"net/http"
)

type Filter struct {
	ProjectName         string `json:"name"`
	FilteredProjectName string `json:"filtered"`
}

type FilterService interface {
}

func NewFilterService(url string) FilterService {
	return &filterService{
		URL:     url,
		TimeLog: []TimeLog{},
		Filters: []Filter{},
	}
}

type filterService struct {
	URL     string
	TimeLog []TimeLog
	Filters []Filter
}

func (s *filterService) GetFilter() (*Filter, error) {
	body, err := http.Get(s.URL)

	if err != nil {
		return nil, err
	}

	defer body.Body.Close()

	var f Filter
	err = json.NewDecoder(body.Body).Decode(&f)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (s *filterService) Filter() []TimeLog {

	tls := s.TimeLog
	for _, fl := range s.Filters {
		for _, tl := range tls {
			if tl.ProjectName == fl.ProjectName {
				tl.ProjectName = fl.FilteredProjectName
			}
		}
	}
	return tls
}
