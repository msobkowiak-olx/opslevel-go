package opslevel_test

import (
	"testing"

	ol "github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2022"
)

func TestCreateWebhookAction(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation WebhookActionCreate($input:CustomActionsWebhookActionCreateInput!){customActionsWebhookActionCreate(input: $input){webhookAction{{ template "custom_actions_request" }},errors{message,path}}}",
		"variables":{"input":{"headers":"{\"Content-Type\":\"application/json\"}","httpMethod":"POST","liquidTemplate":"{\"token\": \"XXX\", \"ref\":\"main\", \"action\": \"rollback\"}","name":"Deploy Rollback","webhookUrl":"https://gitlab.com/api/v4/projects/1/trigger/pipeline"}}
	}`
	response := `{"data": {"customActionsWebhookActionCreate": {
      "webhookAction": {{ template "custom_action1" }},
      "errors": []
  }}}`

	//fmt.Print(Templated(request))
	//fmt.Print(Templated(response))
	//panic(1)

	client := ABetterTestClient(t, "custom_actions/create_action", request, response)

	// Act
	action, err := client.CreateWebhookAction(ol.CustomActionsWebhookActionCreateInput{
		Name:           "Deploy Rollback",
		LiquidTemplate: "{\"token\": \"XXX\", \"ref\":\"main\", \"action\": \"rollback\"}",
		Headers: ol.JSON{
			"Content-Type": "application/json",
		},
		HTTPMethod: ol.CustomActionsHttpMethodEnumPost,
		WebhookURL: "https://gitlab.com/api/v4/projects/1/trigger/pipeline",
	})

	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Deploy Rollback", action.Name)
}

func TestListCustomActions(t *testing.T) {
	//Arrange
	requests := []TestRequest{
		{`{"query": "query ExternalActionList($after:String!$first:Int!){account{customActionsExternalActions(after: $after, first: $first){nodes{aliases,id,description,liquidTemplate,name,... on CustomActionsWebhookAction{headers,httpMethod,webhookUrl}},pageInfo{hasNextPage,hasPreviousPage,startCursor,endCursor},totalCount}}}",
			{{ template "pagination_initial_query_variables" }}
			}`,
			`{
				"data": {
					"account": {
						"customActionsExternalActions": {
							"nodes": [
								{
									{{ template "custom_action1_response" }}
								},
								{
									{{ template "custom_action2_response" }} 
								}
							],
							{{ template "pagination_initial_pageInfo_response" }},
							"totalCount": 2
						  }}}}`},
		{`{"query": "query ExternalActionList($after:String!$first:Int!){account{customActionsExternalActions(after: $after, first: $first){nodes{aliases,id,description,liquidTemplate,name,... on CustomActionsWebhookAction{headers,httpMethod,webhookUrl}},pageInfo{hasNextPage,hasPreviousPage,startCursor,endCursor},totalCount}}}",
			{{ template "pagination_second_query_variables" }}
			}`,
			`{
				"data": {
					"account": {
						"customActionsExternalActions": {
							"nodes": [
								{
									{{ template "custom_action3_response" }}
								}
							],
							{{ template "pagination_second_pageInfo_response" }},
							"totalCount": 1
						  }}}}`},
	}
	// An easy way to see the results of templating is by uncommenting this
	//fmt.Print(Templated(request))
	//fmt.Print(Templated(response))
	//panic(1)

	client := APaginatedTestClient(t, "custom_actions/list_actions", requests...)
	// Act
	response, err := client.ListCustomActions(nil)
	result := response.Nodes
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, 3, len(result))
	autopilot.Equals(t, "Deploy Freeze", result[1].Name)
	autopilot.Equals(t, "Page On-Call", result[2].Name)
	autopilot.Equals(t, "application/json", result[0].Headers["Content-Type"])
}

func TestUpdateWebhookAction(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation WebhookActionUpdate($input:CustomActionsWebhookActionUpdateInput!){customActionsWebhookActionUpdate(input: $input){webhookAction{{ template "custom_actions_request" }},errors{message,path}}}",
		"variables":{"input":{"id": "123456789", "httpMethod":"PUT"}}
	}`
	response := `{"data": {"customActionsWebhookActionUpdate": {
     "webhookAction": {{ template "custom_action1" }},
     "errors": []
 }}}`

	client := ABetterTestClient(t, "custom_actions/update_action", request, response)

	// Act
	action, err := client.UpdateWebhookAction(ol.CustomActionsWebhookActionUpdateInput{
		Id:         "123456789",
		HTTPMethod: ol.CustomActionsHttpMethodEnumPut,
	})

	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Deploy Rollback", action.Name)
}

func TestUpdateWebhookAction2(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation WebhookActionUpdate($input:CustomActionsWebhookActionUpdateInput!){customActionsWebhookActionUpdate(input: $input){webhookAction{{ template "custom_actions_request" }},errors{message,path}}}",
		"variables":{"input":{"id": "123456789","description":"","headers":"{\"Accept\":\"application/json\"}"}}
	}`
	response := `{"data": {"customActionsWebhookActionUpdate": {
     "webhookAction": {{ template "custom_action1" }},
     "errors": []
 }}}`

	client := ABetterTestClient(t, "custom_actions/update_action2", request, response)
	headers := ol.JSON{
		"Accept": "application/json",
	}

	// Act
	action, err := client.UpdateWebhookAction(ol.CustomActionsWebhookActionUpdateInput{
		Id:          "123456789",
		Description: ol.NewString(""),
		Headers:     &headers,
	})

	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Deploy Rollback", action.Name)
}

func TestDeleteWebhookAction(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation WebhookActionDelete($input:IdentifierInput!){customActionsWebhookActionDelete(resource: $input){errors{message,path}}}",
		"variables":{"input":{"id": "123456789"}}
	}`
	response := `{"data": {"customActionsWebhookActionDelete": {
     "errors": []
 }}}`

	client := ABetterTestClient(t, "custom_actions/delete_action", request, response)

	// Act
	err := client.DeleteWebhookAction(ol.IdentifierInput{
		Id: "123456789",
	})

	// Assert
	autopilot.Ok(t, err)
}

func TestCreateTriggerDefinition(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation TriggerDefinitionCreate($input:CustomActionsTriggerDefinitionCreateInput!){customActionsTriggerDefinitionCreate(input: $input){triggerDefinition{{ template "custom_actions_trigger_request" }},errors{message,path}}}",
		"variables":{"input":{"actionId":"123456789", "description":"Disables the Deploy Freeze","entityType":"SERVICE","filterId":"987654321","manualInputsDefinition":"", "name":"Deploy Rollback","ownerId":"123456789", "accessControl": "everyone", "responseTemplate": ""}}
	}`
	response := `{"data": {"customActionsTriggerDefinitionCreate": {
     "triggerDefinition": {{ template "custom_action_trigger1" }},
     "errors": []
 }}}`

	client := ABetterTestClient(t, "custom_actions/create_trigger", request, response)
	// Act
	trigger, err := client.CreateTriggerDefinition(ol.CustomActionsTriggerDefinitionCreateInput{
		Name:        "Deploy Rollback",
		Description: ol.NewString("Disables the Deploy Freeze"),
		Action:      "123456789",
		Owner:       "123456789",
		Filter:      ol.NewID("987654321"),
	})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Release", trigger.Name)
}

func TestCreateTriggerDefinitionWithGlobalEntityType(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation TriggerDefinitionCreate($input:CustomActionsTriggerDefinitionCreateInput!){customActionsTriggerDefinitionCreate(input: $input){triggerDefinition{{ template "custom_actions_trigger_request" }},errors{message,path}}}",
		"variables":{"input":{"actionId":"123456789", "description":"Disables the Deploy Freeze","entityType":"GLOBAL","filterId":"987654321","manualInputsDefinition":"", "name":"Deploy Rollback","ownerId":"123456789", "accessControl": "everyone", "responseTemplate": ""}}
	}`
	response := `{"data": {"customActionsTriggerDefinitionCreate": {
     "triggerDefinition": {{ template "custom_action_trigger1" }},
     "errors": []
 }}}`

	client := ABetterTestClient(t, "custom_actions/create_trigger_with_global_entity", request, response)
	// Act
	trigger, err := client.CreateTriggerDefinition(ol.CustomActionsTriggerDefinitionCreateInput{
		Name:        "Deploy Rollback",
		Description: ol.NewString("Disables the Deploy Freeze"),
		Action:      "123456789",
		Owner:       "123456789",
		Filter:      ol.NewID("987654321"),
		EntityType:  "GLOBAL",
	})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Release", trigger.Name)
}

func TestGetTriggerDefinition(t *testing.T) {
	//Arrange
	request := `{"query":
		"query TriggerDefinitionGet($input:IdentifierInput!){account{customActionsTriggerDefinition(input: $input){{ template "custom_actions_trigger_request" }}}}",
		"variables":{"input":{"id":"123456789"}}
	}`
	response := `{"data": {"account": {
      "customActionsTriggerDefinition": {{ template "custom_action_trigger2" }}
  }}}`

	client := ABetterTestClient(t, "custom_actions/get_trigger", request, response)
	// Act
	trigger, err := client.GetTriggerDefinition(ol.IdentifierInput{Id: "123456789"})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Release", trigger.Name)
	autopilot.Equals(t, "Uses Ruby", trigger.Filter.Name)
	autopilot.Equals(t, "123456789", string(trigger.Owner.Id))
}

func TestListTriggerDefinitions(t *testing.T) {
	//Arrange
	requests := []TestRequest{
		{`{"query": "query TriggerDefinitionList($after:String!$first:Int!){account{customActionsTriggerDefinitions(after: $after, first: $first){nodes{action{aliases,id},aliases,description,filter{id,name},id,manualInputsDefinition,name,owner{alias,id},published,timestamps{createdAt,updatedAt},accessControl,responseTemplate,entityType},pageInfo{hasNextPage,hasPreviousPage,startCursor,endCursor},totalCount}}}",
			{{ template "pagination_initial_query_variables" }}
			}`,
			`{
				"data": {
					"account": {
						"customActionsTriggerDefinitions": {
							"nodes": [
								{
									{{ template "custom_action_trigger1_response" }}
								},
								{
									{{ template "custom_action_trigger2_response" }} 
								}
							],
							{{ template "pagination_initial_pageInfo_response" }},
							"totalCount": 2
						  }}}}`},
		{`{"query": "query TriggerDefinitionList($after:String!$first:Int!){account{customActionsTriggerDefinitions(after: $after, first: $first){nodes{action{aliases,id},aliases,description,filter{id,name},id,manualInputsDefinition,name,owner{alias,id},published,timestamps{createdAt,updatedAt},accessControl,responseTemplate,entityType},pageInfo{hasNextPage,hasPreviousPage,startCursor,endCursor},totalCount}}}",
			{{ template "pagination_second_query_variables" }}
			}`,
			`{
				"data": {
					"account": {
						"customActionsTriggerDefinitions": {
							"nodes": [
								{
									{{ template "custom_action_trigger3_response" }}
								}
							],
							{{ template "pagination_second_pageInfo_response" }},
							"totalCount": 1
						  }}}}`},
	}

	// An easy way to see the results of templating is by uncommenting this
	//fmt.Println(Templated(requests[0].Response))
	//panic(true)

	//"{account{customActionsTriggerDefinitions(after: $after, first: $first){nodes{action{aliases,id},aliases,description,filter{id,name},id,manualInputsDefinition,name,owner{alias,id},published,timestamps{createdAt,updatedAt},accessControl,responseTemplate}}}}",

	client := APaginatedTestClient(t, "custom_actions/list_triggers", requests...)
	// Act
	triggers, err := client.ListTriggerDefinitions(nil)
	result := triggers.Nodes
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, 3, triggers.TotalCount)
	autopilot.Equals(t, "Release", result[1].Name)
	autopilot.Equals(t, "Uses Ruby", result[1].Filter.Name)
	autopilot.Equals(t, "123456789", string(result[1].Owner.Id))
	autopilot.Equals(t, "Rollback", result[2].Name)
	autopilot.Equals(t, "Uses Go", result[2].Filter.Name)
	autopilot.Equals(t, "123456781", string(result[2].Owner.Id))
}

func TestUpdateTriggerDefinition(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation TriggerDefinitionUpdate($input:CustomActionsTriggerDefinitionUpdateInput!){customActionsTriggerDefinitionUpdate(input: $input){triggerDefinition{{ template "custom_actions_trigger_request" }},errors{message,path}}}",
		"variables":{"input":{"id":"123456789", "filterId":null}}
	}`
	response := `{"data": {"customActionsTriggerDefinitionUpdate": {
     "triggerDefinition": {{ template "custom_action_trigger1" }},
     "errors": []
 }}}`

	client := ABetterTestClient(t, "custom_actions/update_trigger", request, response)
	// Act
	trigger, err := client.UpdateTriggerDefinition(ol.CustomActionsTriggerDefinitionUpdateInput{
		Id:     "123456789",
		Filter: ol.NewID(),
	})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Release", trigger.Name)
}

func TestUpdateTriggerDefinition2(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation TriggerDefinitionUpdate($input:CustomActionsTriggerDefinitionUpdateInput!){customActionsTriggerDefinitionUpdate(input: $input){triggerDefinition{{ template "custom_actions_trigger_request" }},errors{message,path}}}",
		"variables":{"input":{"id":"123456789", "name":"test", "description": ""}}
	}`
	response := `{"data": {"customActionsTriggerDefinitionUpdate": {
     "triggerDefinition": {{ template "custom_action_trigger1" }},
     "errors": []
 }}}`

	client := ABetterTestClient(t, "custom_actions/update_trigger2", request, response)
	// Act
	trigger, err := client.UpdateTriggerDefinition(ol.CustomActionsTriggerDefinitionUpdateInput{
		Id:          "123456789",
		Name:        ol.NewString("test"),
		Description: ol.NewString(""),
	})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "Release", trigger.Name)
}

func TestDeleteTriggerDefinition(t *testing.T) {
	//Arrange
	request := `{"query":
		"mutation TriggerDefinitionDelete($input:IdentifierInput!){customActionsTriggerDefinitionDelete(resource: $input){errors{message,path}}}",
		"variables":{"input":{"id":"123456789"}}
	}`
	response := `{"data": {"customActionsTriggerDefinitionDelete": {
     "errors": []
 }}}`
	responseErr := `{"data": {"customActionsTriggerDefinitionDelete": {
     "errors": [{{ template "error1" }}]
 }}}`

	client := ABetterTestClient(t, "custom_actions/delete_trigger", request, response)
	clientErr := ABetterTestClient(t, "custom_actions/delete_trigger_err", request, responseErr)
	clientErr2 := ABetterTestClient(t, "custom_actions/delete_trigger_err2", request, "")
	// Act
	err := client.DeleteTriggerDefinition(ol.IdentifierInput{
		Id: "123456789",
	})
	err2 := clientErr.DeleteTriggerDefinition(ol.IdentifierInput{
		Id: "123456789",
	})
	err3 := clientErr2.DeleteTriggerDefinition(ol.IdentifierInput{
		Id: "123456789",
	})
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, `OpsLevel API Errors:
	- 'one.two.three' Example Error
`, err2.Error())
	autopilot.Assert(t, err3 != nil, "Expected error was not thrown")
}
