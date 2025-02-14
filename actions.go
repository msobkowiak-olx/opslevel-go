package opslevel

import (
	"fmt"
)

type CustomActionsId struct {
	Aliases []string `graphql:"aliases"`
	Id      ID       `graphql:"id"`
}

type CustomActionsExternalAction struct {
	CustomActionsId

	Description    string `graphql:"description"`
	LiquidTemplate string `graphql:"liquidTemplate"`
	Name           string `graphql:"name"`

	CustomActionsWebhookAction `graphql:"... on CustomActionsWebhookAction"`
}

type CustomActionsWebhookAction struct {
	Headers    JSON                        `graphql:"headers" scalar:"true"`
	HTTPMethod CustomActionsHttpMethodEnum `graphql:"httpMethod"`
	WebhookURL string                      `graphql:"webhookUrl"`
}

type CustomActionsTriggerDefinition struct {
	Action                 CustomActionsId                                 `graphql:"action"`
	Aliases                []string                                        `graphql:"aliases"`
	Description            string                                          `graphql:"description"`
	Filter                 FilterId                                        `graphql:"filter"`
	Id                     ID                                              `graphql:"id"`
	ManualInputsDefinition string                                          `graphql:"manualInputsDefinition"`
	Name                   string                                          `graphql:"name"`
	Owner                  TeamId                                          `graphql:"owner"`
	Published              bool                                            `graphql:"published"`
	Timestamps             Timestamps                                      `graphql:"timestamps"`
	AccessControl          CustomActionsTriggerDefinitionAccessControlEnum `graphql:"accessControl"`
	ResponseTemplate       string                                          `graphql:"responseTemplate"`
	EntityType             CustomActionsEntityTypeEnum                     `graphql:"entityType"`
}

func (c *CustomActionsTriggerDefinition) ExtendedTeamAccess(client *Client, variables *PayloadVariables) (*TeamConnection, error) {
	var q struct {
		Account struct {
			CustomActionsTriggerDefinition struct {
				ExtendedTeamAccess TeamConnection `graphql:"extendedTeamAccess(after: $after, first: $first)"`
			} `graphql:"customActionsTriggerDefinition(input: $input)"`
		}
	}
	if c.Id == "" {
		return nil, fmt.Errorf("Unable to get teams with ExtendedTeamAccess, invalid CustomActionsTriggerDefinition id: '%s'", c.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["input"] = *NewIdentifier(string(c.Id))

	if err := client.Query(&q, *variables, WithName("ExtendedTeamAccessList")); err != nil {
		return nil, err
	}

	for q.Account.CustomActionsTriggerDefinition.ExtendedTeamAccess.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.CustomActionsTriggerDefinition.ExtendedTeamAccess.PageInfo.End
		resp, err := c.ExtendedTeamAccess(client, variables)
		if err != nil {
			return nil, err
		}
		q.Account.CustomActionsTriggerDefinition.ExtendedTeamAccess.Nodes = append(q.Account.CustomActionsTriggerDefinition.ExtendedTeamAccess.Nodes, resp.Nodes...)
		q.Account.CustomActionsTriggerDefinition.ExtendedTeamAccess.PageInfo = resp.PageInfo
		q.Account.CustomActionsTriggerDefinition.ExtendedTeamAccess.TotalCount += resp.TotalCount
	}
	return &q.Account.CustomActionsTriggerDefinition.ExtendedTeamAccess, nil
}

type CustomActionsExternalActionsConnection struct {
	Nodes      []CustomActionsExternalAction
	PageInfo   PageInfo
	TotalCount int
}

type CustomActionsTriggerDefinitionsConnection struct {
	Nodes      []CustomActionsTriggerDefinition
	PageInfo   PageInfo
	TotalCount int
}

type CustomActionsWebhookActionCreateInput struct {
	Name           string                      `json:"name"`
	Description    *string                     `json:"description,omitempty"`
	LiquidTemplate string                      `json:"liquidTemplate"`
	WebhookURL     string                      `json:"webhookUrl"`
	HTTPMethod     CustomActionsHttpMethodEnum `json:"httpMethod"`
	Headers        JSON                        `json:"headers"`
}

type CustomActionsWebhookActionUpdateInput struct {
	Id             ID                          `json:"id"`
	Name           *string                     `json:"name,omitempty"`
	Description    *string                     `json:"description,omitempty"`
	LiquidTemplate *string                     `json:"liquidTemplate,omitempty"`
	WebhookURL     *string                     `json:"webhookUrl,omitempty"`
	HTTPMethod     CustomActionsHttpMethodEnum `json:"httpMethod,omitempty"`
	Headers        *JSON                       `json:"headers,omitempty"`
}

type CustomActionsTriggerDefinitionCreateInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Owner       ID      `json:"ownerId"`
	// In the API actionID is `ID!` but that's because of the CustomActionsWebhookActionCreateInput
	// But we are not implementing that because it is used for the UI, so we need to enforce an actionId is given
	Action ID  `json:"actionId"`
	Filter *ID `json:"filterId,omitempty"`
	// This is being explicitly left out to reduce the complexity of the implementation
	// action *CustomActionsWebhookActionCreateInput
	ManualInputsDefinition string                                          `json:"manualInputsDefinition"`
	Published              *bool                                           `json:"published,omitempty"`
	AccessControl          CustomActionsTriggerDefinitionAccessControlEnum `json:"accessControl"`
	ResponseTemplate       string                                          `json:"responseTemplate"`
	EntityType             CustomActionsEntityTypeEnum                     `json:"entityType"`
	ExtendedTeamAccess     *[]IdentifierInput                              `json:"extendedTeamAccess,omitempty"`
}

type CustomActionsTriggerDefinitionUpdateInput struct {
	Id          ID      `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Owner       *ID     `json:"ownerId,omitempty"`
	Action      *ID     `json:"actionId,omitempty"`
	Filter      *ID     `json:"filterId,omitempty"`
	// This is being explicitly left out to reduce the complexity of the implementation
	// action *CustomActionsWebhookActionCreateInput
	ManualInputsDefinition *string                                         `json:"manualInputsDefinition,omitempty"`
	Published              *bool                                           `json:"published,omitempty"`
	AccessControl          CustomActionsTriggerDefinitionAccessControlEnum `json:"accessControl,omitempty"`
	ResponseTemplate       *string                                         `json:"responseTemplate,omitempty"`
	EntityType             CustomActionsEntityTypeEnum                     `json:"entityType,omitempty"`
	ExtendedTeamAccess     *[]IdentifierInput                              `json:"extendedTeamAccess,omitempty"`
}

func (client *Client) CreateWebhookAction(input CustomActionsWebhookActionCreateInput) (*CustomActionsExternalAction, error) {
	var m struct {
		Payload struct {
			WebhookAction CustomActionsExternalAction
			Errors        []OpsLevelErrors
		} `graphql:"customActionsWebhookActionCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("WebhookActionCreate"))
	return &m.Payload.WebhookAction, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) GetCustomAction(input IdentifierInput) (*CustomActionsExternalAction, error) {
	var q struct {
		Account struct {
			Action CustomActionsExternalAction `graphql:"customActionsExternalAction(input: $input)"`
		}
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Query(&q, v, WithName("ExternalActionGet"))
	if q.Account.Action.Id == "" {
		err = fmt.Errorf("CustomActionsExternalAction with ID '%s' or Alias '%s' not found", input.Id, input.Alias)
	}
	return &q.Account.Action, HandleErrors(err, nil)
}

func (client *Client) ListCustomActions(variables *PayloadVariables) (CustomActionsExternalActionsConnection, error) {
	var q struct {
		Account struct {
			Actions CustomActionsExternalActionsConnection `graphql:"customActionsExternalActions(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	if err := client.Query(&q, *variables, WithName("ExternalActionList")); err != nil {
		return CustomActionsExternalActionsConnection{}, err
	}
	for q.Account.Actions.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Actions.PageInfo.End
		resp, err := client.ListCustomActions(variables)
		if err != nil {
			return CustomActionsExternalActionsConnection{}, err
		}
		q.Account.Actions.Nodes = append(q.Account.Actions.Nodes, resp.Nodes...)
		q.Account.Actions.PageInfo = resp.PageInfo
	}
	return q.Account.Actions, nil
}

func (client *Client) UpdateWebhookAction(input CustomActionsWebhookActionUpdateInput) (*CustomActionsExternalAction, error) {
	var m struct {
		Payload struct {
			WebhookAction CustomActionsExternalAction
			Errors        []OpsLevelErrors
		} `graphql:"customActionsWebhookActionUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("WebhookActionUpdate"))
	return &m.Payload.WebhookAction, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) DeleteWebhookAction(input IdentifierInput) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"customActionsWebhookActionDelete(resource: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("WebhookActionDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

func (client *Client) CreateTriggerDefinition(input CustomActionsTriggerDefinitionCreateInput) (*CustomActionsTriggerDefinition, error) {
	var m struct {
		Payload struct {
			TriggerDefinition CustomActionsTriggerDefinition
			Errors            []OpsLevelErrors
		} `graphql:"customActionsTriggerDefinitionCreate(input: $input)"`
	}
	if input.AccessControl == "" {
		input.AccessControl = CustomActionsTriggerDefinitionAccessControlEnumEveryone
	}
	if input.EntityType == "" {
		input.EntityType = CustomActionsEntityTypeEnumService
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("TriggerDefinitionCreate"))
	return &m.Payload.TriggerDefinition, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) GetTriggerDefinition(input IdentifierInput) (*CustomActionsTriggerDefinition, error) {
	var q struct {
		Account struct {
			Definition CustomActionsTriggerDefinition `graphql:"customActionsTriggerDefinition(input: $input)"`
		}
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Query(&q, v, WithName("TriggerDefinitionGet"))
	if q.Account.Definition.Id == "" {
		err = fmt.Errorf("CustomActionsTriggerDefinition with ID '%s' or Alias '%s' not found", input.Id, input.Alias)
	}
	return &q.Account.Definition, HandleErrors(err, nil)
}

func (client *Client) ListTriggerDefinitions(variables *PayloadVariables) (CustomActionsTriggerDefinitionsConnection, error) {
	var q struct {
		Account struct {
			Definitions CustomActionsTriggerDefinitionsConnection `graphql:"customActionsTriggerDefinitions(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	if err := client.Query(&q, *variables, WithName("TriggerDefinitionList")); err != nil {
		return CustomActionsTriggerDefinitionsConnection{}, err
	}
	for q.Account.Definitions.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Definitions.PageInfo.End
		resp, err := client.ListTriggerDefinitions(variables)
		if err != nil {
			return CustomActionsTriggerDefinitionsConnection{}, err
		}
		q.Account.Definitions.Nodes = append(q.Account.Definitions.Nodes, resp.Nodes...)
		q.Account.Definitions.PageInfo = resp.PageInfo
		q.Account.Definitions.TotalCount += resp.TotalCount
	}
	return q.Account.Definitions, nil
}

func (client *Client) UpdateTriggerDefinition(input CustomActionsTriggerDefinitionUpdateInput) (*CustomActionsTriggerDefinition, error) {
	var m struct {
		Payload struct {
			TriggerDefinition CustomActionsTriggerDefinition
			Errors            []OpsLevelErrors
		} `graphql:"customActionsTriggerDefinitionUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("TriggerDefinitionUpdate"))
	return &m.Payload.TriggerDefinition, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) DeleteTriggerDefinition(input IdentifierInput) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"customActionsTriggerDefinitionDelete(resource: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("TriggerDefinitionDelete"))
	return HandleErrors(err, m.Payload.Errors)
}
