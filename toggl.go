package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Toggl struct {
	APIToken string
	BaseURL  string
}

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectMap map[int]string

type TimeEntry struct {
	ID          int    `json:"id"`
	ProjectId   int    `json:"project_id"`
	Start       string `json:"start"`
	Stop        string `json:"stop"`
	Description string `json:"description"`
}

func NewToggl(apiToken, baseURL string) *Toggl {
	return &Toggl{
		APIToken: apiToken,
		BaseURL:  baseURL,
	}
}

func (t *Toggl) GetTimeLogs(c echo.Context) ([]TimeLog, error) {
	ps, err := t.GetProjects(c)
	if err != nil {
		return nil, err
	}

	ts, err := t.GetTimeEntries(c)
	if err != nil {
		return nil, err
	}

	return t.Convert(c, ps, ts)
}

func (t *Toggl) Convert(c echo.Context, ps []Project, ts []TimeEntry) ([]TimeLog, error) {

	pm := make(ProjectMap)

	for _, p := range ps {
		pm[p.Id] = p.Name
	}

	var timeLogs []TimeLog

	for _, t := range ts {
		timeLogs = append(timeLogs, TimeLog{
			ID:          t.ID,
			Title:       t.Description,
			ProjectName: pm[t.ProjectId],
			Start:       toTime(t.Start),
			Stop:        toTime(t.Stop),
		})
	}

	return timeLogs, nil
}

func toTime(s string) time.Time {
	tm, _ := time.Parse(time.RFC3339, s)

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	return tm.In(jst)
}

func (t *Toggl) GetProjects(c echo.Context) ([]Project, error) {
	res, err := t.get(c, "/me/projects")
	if err != nil {
		return nil, err
	}

	var projects []Project

	err = json.Unmarshal(res, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (t *Toggl) GetTimeEntries(c echo.Context) ([]TimeEntry, error) {
	res, err := t.get(c, "/me/time_entries")
	if err != nil {
		return nil, err
	}

	var timeEntries []TimeEntry

	err = json.Unmarshal(res, &timeEntries)
	if err != nil {
		return nil, err
	}

	return timeEntries, nil
}

func (t *Toggl) get(c echo.Context, path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet,
		t.BaseURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(t.APIToken, "api_token")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
