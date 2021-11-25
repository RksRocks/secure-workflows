package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const HardenRunnerActionNameWithBranch = "step-security/harden-runner@main"
const HardenRunnerActionName = "step-security/harden-runner"

func SecureWorkflow(inputYaml string, svc dynamodbiface.DynamoDBAPI) (*FixWorkflowPermsReponse, error) {
	fixResponse, err := AddJobLevelPermissions(inputYaml, svc)
	if err != nil {
		return nil, err
	} else {

		if len(fixResponse.MissingActions) > 0 {
			StoreMissingActions(fixResponse.MissingActions, svc)
		}

		fixResponse.FinalOutput, _ = AddAction(fixResponse.FinalOutput, HardenRunnerActionNameWithBranch)

		fixResponse.FinalOutput, _ = PinActions(fixResponse.FinalOutput)

		if !fixResponse.HasErrors {
			fixResponse.FinalOutput, _ = AddWorkflowLevelPermissions(fixResponse.FinalOutput)
		}

		return fixResponse, nil
	}
}