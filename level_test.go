package opslevel_test

import (
	"testing"

	ol "github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2023"
)

func TestCreateRubricLevels(t *testing.T) {
	testRequest := NewTestRequest(
		`"mutation LevelCreate($input:LevelCreateInput!){levelCreate(input: $input){level{alias,description,id,index,name},errors{message,path}}}"`,
		`{"input": { "name": "Kyle", "description": "Created By Kyle", "index": 4 }}`,
		`{"data": { "levelCreate": { "level": { "alias": "kyle", "description": "Created By Kyle", "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNDgw", "index": 4, "name": "Kyle" }, "errors": [] }}}`,
	)

	client := BestTestClient(t, "rubric/level_create", testRequest)
	// Act
	result, _ := client.CreateLevel(ol.LevelCreateInput{
		Name:        "Kyle",
		Description: "Created By Kyle",
		Index:       ol.NewInt(4),
	})
	// Assert
	autopilot.Equals(t, "kyle", result.Alias)
	autopilot.Equals(t, 4, result.Index)
}

func TestGetRubricLevel(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"query LevelGet($id:ID!){account{level(id: $id){alias,description,id,index,name}}}"`,
		`{"id": "Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tsaXN0LzYyMg"}`,
		`{"data": {
        "account": {
          "level": {
            "alias": "bronze",
            "description": "Services in this level satisfy critical checks. This is the minimum standard to ship to production.",
            "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvMzE3",
            "index": 1,
            "name": "Bronze"
          }}}}`,
	)

	client := BestTestClient(t, "rubric/level_get", testRequest)
	// Act
	result, err := client.GetLevel("Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tsaXN0LzYyMg")
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, "Bronze", result.Name)
}

func TestGetMissingRubricLevel(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"query LevelGet($id:ID!){account{level(id: $id){alias,description,id,index,name}}}"`,
		`{"id": "Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tsaXN0LzYyMg"}`,
		`{"data": { "account": { "level": null }}}`,
	)

	client := BestTestClient(t, "rubric/level_get_missing", testRequest)
	// Act
	_, err := client.GetLevel("Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tsaXN0LzYyMg")
	// Assert
	autopilot.Assert(t, err != nil, "This test should throw an error.")
}

func TestListRubricLevels(t *testing.T) {
	// Arrange
	client := ATestClient(t, "rubric/level/list")
	// Act
	result, _ := client.ListLevels()
	// Assert
	autopilot.Equals(t, 4, len(result))
	autopilot.Equals(t, "Bronze", result[1].Name)
}

func TestUpdateRubricLevel(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation LevelUpdate($input:LevelUpdateInput!){levelUpdate(input: $input){level{alias,description,id,index,name},errors{message,path}}}"`,
		`{"input": { {{ template "id1" }}, "name": "{{ template "name1" }}", "description": "{{ template "description" }}" }}`,
		`{"data": { "levelUpdate": { "level": {{ template "level_1" }}, "errors": [] }}}`,
	)

	client := BestTestClient(t, "rubric/level_update", testRequest)
	// Act
	result, _ := client.UpdateLevel(ol.LevelUpdateInput{
		Id:          id1,
		Name:        "Example",
		Description: ol.NewString("An example description"),
	})
	// Assert
	autopilot.Equals(t, "example", result.Alias)
	autopilot.Equals(t, "Example", result.Name)
	autopilot.Equals(t, "An example description", result.Description)
}

func TestUpdateRubricLevelNoName(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation LevelUpdate($input:LevelUpdateInput!){levelUpdate(input: $input){level{alias,description,id,index,name},errors{message,path}}}"`,
		`{"input": { {{ template "id1" }}, "description": "{{ template "description" }}" } }`,
		`{"data": { "levelUpdate": { "level": {{ template "level_1" }}, "errors": [] }}}`,
	)

	client := BestTestClient(t, "rubric/level_update_noname", testRequest)
	// Act
	result, _ := client.UpdateLevel(ol.LevelUpdateInput{
		Id:          id1,
		Description: ol.NewString("An example description"),
	})
	// Assert
	autopilot.Equals(t, "example", result.Alias)
	autopilot.Equals(t, "Example", result.Name)
	autopilot.Equals(t, "An example description", result.Description)
}

func TestUpdateRubricLevelEmptyDescription(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation LevelUpdate($input:LevelUpdateInput!){levelUpdate(input: $input){level{alias,description,id,index,name},errors{message,path}}}"`,
		`{"input": { {{ template "id1" }}, "name": "{{ template "name1" }}", "description": "" }}`,
		`{"data": { "levelUpdate": { "level": {{ template "level_1" }}, "errors": [] }}}`,
	)

	client := BestTestClient(t, "rubric/level_update_emptydescription", testRequest)
	// Act
	result, _ := client.UpdateLevel(ol.LevelUpdateInput{
		Id:          id1,
		Name:        "Example",
		Description: ol.NewString(""),
	})
	// Assert
	autopilot.Equals(t, "example", result.Alias)
	autopilot.Equals(t, "Example", result.Name)
	autopilot.Equals(t, "An example description", result.Description)
}

func TestUpdateRubricLevelNoDescription(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation LevelUpdate($input:LevelUpdateInput!){levelUpdate(input: $input){level{alias,description,id,index,name},errors{message,path}}}"`,
		`{"input": { {{ template "id1" }}, "name": "{{ template "name1" }}" }}`,
		`{"data": { "levelUpdate": { "level": {{ template "level_1" }}, "errors": [] }}}`,
	)

	client := BestTestClient(t, "rubric/level_update_nodescription", testRequest)
	// Act
	result, _ := client.UpdateLevel(ol.LevelUpdateInput{
		Id:   id1,
		Name: "Example",
	})
	// Assert
	autopilot.Equals(t, "example", result.Alias)
	autopilot.Equals(t, "Example", result.Name)
	autopilot.Equals(t, "An example description", result.Description)
}

func TestDeleteRubricLevels(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation LevelDelete($input:LevelDeleteInput!){levelDelete(input: $input){deletedLevelId,errors{message,path}}}"`,
		`{"input": { "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNDgw" }}`,
		`{"data": { "levelDelete": { "deletedLevelId": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNDgw", "errors": [] }}}`,
	)

	client := BestTestClient(t, "rubric/level_delete", testRequest)
	// Act
	err := client.DeleteLevel("Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNDgw")
	// Assert
	autopilot.Equals(t, nil, err)
}
