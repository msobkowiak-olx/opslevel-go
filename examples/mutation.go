package opslevel_example

import (
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"
)

func init() {
	var mutation struct {
		Payload struct {
			Aliases []string
			OwnerId string
			Errors  []opslevel.OpsLevelErrors
		} `graphql:"aliasCreate(input: $input)"`
	}
	variables := opslevel.PayloadVariables{
		"input": opslevel.AliasCreateInput{
			Alias:   "MyNewAlias",
			OwnerId: "XXXXXXXXXXX",
		},
	}

	client := opslevel.NewClient("xxx")
	if err := client.Mutate(&mutation, variables); err != nil {
		panic(err)
	}
	for _, alias := range mutation.Payload.Aliases {
		fmt.Println(alias)
	}
}
