package apigateway

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/jsii-runtime-go"
)

func createTestStack(t *testing.T) awscdk.Stack {
	t.Helper()

	app := awscdk.NewApp(nil)
	stack := awscdk.NewStack(app, jsii.String("testStack"), &awscdk.StackProps{
		Env: &awscdk.Environment{
			Account: jsii.String("012345678901"),
			Region:  jsii.String("eu-west-1"),
		},
	})

	return stack
}

func createTestHostedZone(t *testing.T, stack awscdk.Stack) awsroute53.HostedZone {
	t.Helper()

	return awsroute53.NewHostedZone(stack, jsii.String("testHostedZone"), &awsroute53.HostedZoneProps{
		ZoneName: jsii.String("testing.verde.systems"),
	})
}

func createTestAPIGateway(t *testing.T, stack awscdk.Stack, props *awsapigateway.RestApiProps) awsapigateway.RestApi {
	t.Helper()

	apiName := "testApi"
	if props == nil {
		props = &awsapigateway.RestApiProps{
			RestApiName: jsii.String(apiName),
		}
	}

	api := awsapigateway.NewRestApi(stack, jsii.String(apiName), props)
	api.Root().AddMethod(jsii.String("GET"), nil, nil)

	return api
}
