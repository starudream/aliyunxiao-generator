package api

import (
	"strings"
)

type WorkItem struct {
	Identifier      string `json:"identifier"`
	Subject         string `json:"subject"`
	SpaceIdentifier string `json:"spaceIdentifier"`
	SpaceName       string `json:"spaceName"`

	WorkitemTypeName string `json:"workitemTypeName"`

	WorkitemType struct {
		DisplayName string `json:"displayName"`
	} `json:"workitemType"`
}

func ListWorkItem(day string) ([]*WorkItem, error) {
	r := R().SetQueryParam("day", day)
	return Exec[[]*WorkItem](r, "GET", "/workitem/workitem/workTable/workitemTime/list")
}

type WorkItemTime struct {
	WorkitemIdentifier string `json:"workitemIdentifier"`
	RecordedHours      int    `json:"recordedHours"`
}

func ListWorkItemTime(day string, identifiers []string) (map[string]int, error) {
	if len(identifiers) == 0 {
		return nil, nil
	}
	r := R().SetQueryParam("day", day).SetQueryParam("workitemIdentifiers", strings.Join(identifiers, ","))
	resp, err := Exec[[]*WorkItemTime](r, "GET", "/workitem/workitem/time/stats/user/getDayWorkitemTimeByWorkitems")
	if err != nil {
		return nil, err
	}
	data := map[string]int{}
	for _, v := range resp {
		data[v.WorkitemIdentifier] = v.RecordedHours
	}
	return data, nil
}
