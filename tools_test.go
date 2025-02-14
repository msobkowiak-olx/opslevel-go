package opslevel_test

import (
	"testing"

	ol "github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2023"
)

func TestCreateTool(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation ToolCreate($input:ToolCreateInput!){toolCreate(input: $input){tool{category,categoryAlias,displayName,environment,id,url,service{id,aliases}},errors{message,path}}}"`,
		`{ "input": { "category": "other", "displayName": "example", "serviceId": "{{ template "id1_string" }}", "url": "https://example.com" }}`,
		`{"data": { "toolCreate": { "tool": {{ template "tool_1" }}, "errors": [] }}}`,
	)
	client := BestTestClient(t, "toolCreate", testRequest)
	// Act
	result, err := client.CreateTool(ol.ToolCreateInput{
		Category:    ol.ToolCategoryOther,
		DisplayName: "example",
		ServiceId:   id1,
		Url:         "https://example.com",
	})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, id1, result.Service.Id)
	autopilot.Equals(t, ol.ToolCategoryOther, result.Category)
	autopilot.Equals(t, "Example", result.DisplayName)
	autopilot.Equals(t, "https://example.com", result.Url)
}

func TestUpdateTool(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation ToolUpdate($input:ToolUpdateInput!){toolUpdate(input: $input){tool{category,categoryAlias,displayName,environment,id,url,service{id,aliases}},errors{message,path}}}"`,
		`{ "input": { {{ template "id1" }}, "category": "deployment" }}`,
		`{"data": { "toolUpdate": { "tool": {{ template "tool_1_update" }}, "errors": [] }}}`,
	)
	client := BestTestClient(t, "toolUpdate", testRequest)
	// Act
	result, err := client.UpdateTool(ol.ToolUpdateInput{
		Id:       id1,
		Category: ol.ToolCategoryDeployment,
	})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, ol.ToolCategoryDeployment, result.Category)
	autopilot.Equals(t, "prod", result.Environment)
}

func TestDeleteTool(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation ToolDelete($input:ToolDeleteInput!){toolDelete(input: $input){errors{message,path}}}"`,
		`{ "input": { {{ template "id1" }} } }`,
		`{"data": { "toolDelete": { "errors": [] }}}`,
	)
	client := BestTestClient(t, "toolDelete", testRequest)
	// Act
	err := client.DeleteTool(id1)
	// Assert
	autopilot.Ok(t, err)
}
