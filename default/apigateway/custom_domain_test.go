package apigateway

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestCustomDomainCdk(t *testing.T) {

	t.Run("createCustomDomain/ok", func(t *testing.T) {
		stack := createTestStack(t)
		hostedZone := createTestHostedZone(t, stack)
		api := createTestAPIGateway(t, stack, nil)
		props := CustomDomainProps{
			Api:           api,
			HostedZone:    hostedZone,
			SubdomainName: jsii.String("test"),
		}

		NewCustomDomain(stack, jsii.String("testCustomDomain"), &props)

		template := assertions.Template_FromStack(stack, nil)
		template.HasResourceProperties(jsii.String("AWS::ApiGateway::DomainName"), map[string]any{
			"DomainName": "test.testing.verde.systems",
		})

		template.HasResource(jsii.String("AWS::ApiGateway::BasePathMapping"), map[string]any{
			"Properties": map[string]any{},
		})

		template.HasResourceProperties(jsii.String("AWS::CertificateManager::Certificate"), map[string]any{
			"DomainName":   "test.testing.verde.systems",
			"KeyAlgorithm": "RSA_2048",
		})

		// in Route53 records "Name" should be domain name with "." at the end
		// more information https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53-recordset.html#cfn-route53-recordset-name
		template.HasResourceProperties(jsii.String("AWS::Route53::RecordSet"), map[string]any{
			"Name": "test.testing.verde.systems.",
			"Type": "A",
		})

		template.HasResourceProperties(jsii.String("AWS::Route53::RecordSet"), map[string]any{
			"Name": "test.testing.verde.systems.",
			"Type": "CAA",
		})
	})

	t.Run("createCustomDomain/validation_panic", func(t *testing.T) {
		stack := createTestStack(t)

		defer func() {
			if r := recover(); r == nil {
				t.Error("panic expected")
			}
		}()

		createCustomDomain(stack, jsii.String("testCustomDomain"), nil)

	})

	t.Run("NewCustomDomain", func(t *testing.T) {
		stack := createTestStack(t)
		hostedZone := createTestHostedZone(t, stack)
		api := createTestAPIGateway(t, stack, nil)
		props := CustomDomainProps{
			Api:           api,
			HostedZone:    hostedZone,
			SubdomainName: jsii.String("test"),
		}

		NewCustomDomain(stack, jsii.String("testCustomDomain"), &props)

		template := assertions.Template_FromStack(stack, nil)

		template.HasResourceProperties(jsii.String("AWS::ApiGateway::DomainName"), map[string]any{
			"Tags": []map[string]any{
				{"Key": "test:operations:construct:NewCustomDomain", "Value": "0.0.1"},
			},
		})

		template.HasResourceProperties(jsii.String("AWS::CertificateManager::Certificate"), map[string]any{
			"Tags": []map[string]any{
				{"Key": "Name", "Value": "testCustomDomain"},
				{"Key": "test:operations:construct:NewCustomDomain", "Value": "0.0.1"},
			},
		})
	})

}
