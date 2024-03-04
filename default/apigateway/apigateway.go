package apigateway

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PublicLambdaRestApiProps struct {
	awsapigateway.LambdaRestApiProps
	HostedZone    awsroute53.IHostedZone
	SubdomainName *string
}

type PrivateLambdaRestApiProps struct {
	awsapigateway.LambdaRestApiProps
	VpcEndpoints []awsec2.IVpcEndpoint
}

// Create Public API Gateway
func NewPublicLambdaRestApi(scope constructs.Construct, id *string, props *PublicLambdaRestApiProps) awsapigateway.LambdaRestApi {
	p := awsapigateway.LambdaRestApiProps{}
	if props != nil {
		p = props.LambdaRestApiProps
	}

	// ??: Should we check EndpointTypes?
	if p.EndpointConfiguration == nil && p.EndpointTypes != nil {
		p.EndpointConfiguration = &awsapigateway.EndpointConfiguration{
			Types: p.EndpointTypes,
		}
	}

	api := createLambdaApiGateway(scope, id, &p)
	if props.HostedZone != nil && props.SubdomainName != nil {
		createCustomDomain(scope, id, &CustomDomainProps{
			Api:           api,
			HostedZone:    props.HostedZone,
			SubdomainName: props.SubdomainName,
		})
	}

	return api
}

func NewPrivateAPIGateway(scope constructs.Construct, id *string, props *PrivateLambdaRestApiProps) awsapigateway.LambdaRestApi {
	p := awsapigateway.LambdaRestApiProps{}
	if props != nil {
		p = props.LambdaRestApiProps
	}

	if p.EndpointConfiguration == nil {
		p.EndpointConfiguration = &awsapigateway.EndpointConfiguration{
			Types: &[]awsapigateway.EndpointType{
				awsapigateway.EndpointType_PRIVATE,
			},
			VpcEndpoints: &props.VpcEndpoints,
		}
	}

	return createLambdaApiGateway(scope, id, &p)
}

func (props *PrivateLambdaRestApiProps) Validate() {
	t := *props.EndpointTypes
	if t == nil || len(t) != 1 || t[0] != awsapigateway.EndpointType_PRIVATE {
		panic("Endpoint type should be only PRIVATE")
	}

	if props.EndpointConfiguration == nil || props.EndpointConfiguration.VpcEndpoints == nil || len(props.VpcEndpoints) == 0 {
		panic("VPC endpoints required for private API Gateway!")
	}
}

func createLambdaApiGateway(scope constructs.Construct, id *string, props *awsapigateway.LambdaRestApiProps) awsapigateway.LambdaRestApi {
	if props.CloudWatchRole == nil {
		props.CloudWatchRole = jsii.Bool(false)
	}

	if props.DeployOptions == nil {
		props.DeployOptions = &awsapigateway.StageOptions{
			LoggingLevel:   awsapigateway.MethodLoggingLevel_INFO,
			TracingEnabled: jsii.Bool(true),
		}
	}

	if props.Proxy == nil {
		props.Proxy = jsii.Bool(false)
	}

	if props.Proxy == nil {
		props.Proxy = jsii.Bool(false)
	}

	if props.EndpointExportName == nil {
		props.EndpointExportName = id
	}

	return awsapigateway.NewLambdaRestApi(scope, id, props)
}
