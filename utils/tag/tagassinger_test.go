package tag

import (
	"fmt"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/jsii-runtime-go"
)

func TestTagAssigner(t *testing.T) {

	t.Run("AddDevOpsStackNameTag", func(t *testing.T) {
		stack := createTestStack(t)
		createTestAPIGateway(t, stack, nil)

		AddDevOpsStackNameTag(stack, &TagProps{TagValue: jsii.String("value")})

		want := []map[string]string{
			{
				"Key":   fmt.Sprintf("%s:%s:%s", orgTabLabel, string(devOpsUseCaseTagLabel), "stack-name"),
				"Value": "value",
			},
		}

		template := assertions.Template_FromStack(stack, nil)

		template.HasResourceProperties(jsii.String("AWS::ApiGateway::RestApi"), map[string]any{
			"Tags": want,
		})
	})
	t.Run("AddDevOpsConstructTag", func(t *testing.T) {
		stack := createTestStack(t)
		createTestAPIGateway(t, stack, nil)

		AddDevOpsConstructTag(stack, &TagProps{TagValue: jsii.String("0.0.1"), TagName: jsii.String("TestConstruct")})

		want := []map[string]string{
			{
				"Key":   fmt.Sprintf("%s:%s:%s", orgTabLabel, string(devOpsUseCaseConstructTagLabel), "TestConstruct"),
				"Value": "0.0.1",
			},
		}

		template := assertions.Template_FromStack(stack, nil)

		template.HasResourceProperties(jsii.String("AWS::ApiGateway::RestApi"), map[string]any{
			"Tags": want,
		})
	})
}

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
