package opslevel_test

import (
	"fmt"
	"testing"

	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2023"
)

func TestCreateSecret(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation SecretsVaultsSecretCreate($alias:String!$input:SecretInput!){secretsVaultsSecretCreate(alias: $alias, input: $input){secret{alias,id,owner{alias,id},timestamps{createdAt,updatedAt}},errors{message,path}}}"`,
		`{{ template "secret_create_vars" }}`,
		`{{ template "secret_create_response" }}`,
	)
	client := BestTestClient(t, "secrets/create", testRequest)
	fmt.Println(client)
	// Act
	secretInput := opslevel.SecretInput{
		Owner: opslevel.IdentifierInput{Id: id2},
		Value: "my-secret",
	}
	result, err := client.CreateSecret("alias1", secretInput)

	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, id2, result.Owner.Id)
}

func TestGetSecret(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"query SecretsVaultsSecret($input:IdentifierInput!){account{secretsVaultsSecret(input: $input){alias,id,owner{alias,id},timestamps{createdAt,updatedAt}}}}"`,
		`{{ template "secret_get_vars" }}`,
		`{{ template "secret_get_response" }}`,
	)
	client := BestTestClient(t, "secret/get", testRequest)
	// Act
	result, err := client.GetSecret(string(id2))
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id2, result.ID)
}

func TestListSecrets(t *testing.T) {
	// Arrange
	testRequestOne := NewTestRequest(
		`"query SecretList($after:String!$first:Int!){account{secretsVaultsSecrets(after: $after, first: $first){nodes{alias,id,owner{alias,id},timestamps{createdAt,updatedAt}},{{ template "pagination_request" }}}}}"`,
		`{{ template "pagination_initial_query_variables" }}`,
		`{{ template "secret_list_response_1" }}`,
	)
	testRequestTwo := NewTestRequest(
		`"query SecretList($after:String!$first:Int!){account{secretsVaultsSecrets(after: $after, first: $first){nodes{alias,id,owner{alias,id},timestamps{createdAt,updatedAt}},{{ template "pagination_request" }}}}}"`,
		`{{ template "pagination_second_query_variables" }}`,
		`{{ template "secret_list_response_2" }}`,
	)
	requests := []TestRequest{testRequestOne, testRequestTwo}

	client := BestTestClient(t, "secrets/list", requests...)
	// Act
	secretsVaultsSecretConnection, err := client.ListSecretsVaultsSecret(nil)
	secretNode := secretsVaultsSecretConnection.Nodes
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, 3, secretsVaultsSecretConnection.TotalCount)
	autopilot.Equals(t, "example_2", secretNode[1].Alias)
	autopilot.Equals(t, secretNode[1].Alias, secretNode[1].Owner.Alias)
}

func TestUpdateSecret(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation SecretsVaultsSecretUpdate($input:SecretInput!$secret:IdentifierInput!){secretsVaultsSecretUpdate(input: $input, secret: $secret){secret{alias,id,owner{alias,id},timestamps{createdAt,updatedAt}},errors{message,path}}}"`,
		`{{ template "secret_update_vars" }}`,
		`{{ template "secret_update_response" }}`,
	)
	client := BestTestClient(t, "secrets/update", testRequest)
	// Act
	secretInput := opslevel.SecretInput{
		Owner: opslevel.IdentifierInput{Id: id2},
		Value: "secret_value_2",
	}
	result, err := client.UpdateSecret(string(id2), secretInput)
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, id2, result.ID)
	autopilot.Equals(t, id2, result.Owner.Id)
}

func TestDeleteSecrets(t *testing.T) {
	// Arrange
	testRequest := NewTestRequest(
		`"mutation SecretsVaultsSecretDelete($input:IdentifierInput!){secretsVaultsSecretDelete(resource: $input){errors{message,path}}}"`,
		`{{ template "secret_delete_vars" }}`,
		`{{ template "secret_delete_response" }}`,
	)

	client := BestTestClient(t, "secrets/delete", testRequest)
	// Act
	err := client.DeleteSecret(string(id1))
	// Assert
	autopilot.Equals(t, nil, err)
}
