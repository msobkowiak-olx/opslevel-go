package opslevel

type HasRecentDeployCheckFragment struct {
	Days int `graphql:"days"`
}

type CheckHasRecentDeployCreateInput struct {
	CheckCreateInput

	Days int `json:"days"`
}

type CheckHasRecentDeployUpdateInput struct {
	CheckUpdateInput

	Days *int `json:"days,omitempty"`
}

func (client *Client) CreateCheckHasRecentDeploy(input CheckHasRecentDeployCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkHasRecentDeployCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckHasRecentDeployCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckHasRecentDeploy(input CheckHasRecentDeployUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkHasRecentDeployUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckHasRecentDeployUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
