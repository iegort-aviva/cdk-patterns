package apigateway

import (
	"fmt"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"github.com/iegort-aviva/cdk-patterns/utils/tag"
)

type CustomDomainProps struct {
	Api           awsapigateway.IRestApi
	HostedZone    awsroute53.IHostedZone
	SubdomainName *string
}

// Creating custom domain for API Gateway
func NewCustomDomain(scope constructs.Construct, id *string, props *CustomDomainProps) {
	tag.AddDevOpsConstructTag(scope, &tag.TagProps{TagName: jsii.String("NewCustomDomain"), TagValue: jsii.String("0.0.1")})

	createCustomDomain(scope, id, props)
}

func (props *CustomDomainProps) Validate() error {
	errorMessages := make([]string, 0, 3)
	if props.Api == nil {
		errorMessages = append(errorMessages, "api construct required")
	}

	if props.HostedZone == nil {
		errorMessages = append(errorMessages, "route53 hosted zone required")
	}

	if props.SubdomainName == nil {
		errorMessages = append(errorMessages, "subdomain name required")
	}

	if len(errorMessages) > 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errorMessages, ", "))
}

func createCustomDomain(scope constructs.Construct, id *string, props *CustomDomainProps) map[string]constructs.Construct {
	p := CustomDomainProps{}
	if props == nil {
		props = &p
	}

	if err := props.Validate(); err != nil {
		panic(err)
	}

	customDomain := fmt.Sprintf("%s.%s", *props.SubdomainName, *props.HostedZone.ZoneName())

	// create certificate
	acmCertId := fmt.Sprintf("%s-%s", *id, "Certificate")
	acmCert := awscertificatemanager.NewCertificate(scope, jsii.String(acmCertId), &awscertificatemanager.CertificateProps{
		DomainName:      jsii.String(customDomain),
		CertificateName: id,
		KeyAlgorithm:    awscertificatemanager.KeyAlgorithm_RSA_2048(),
		Validation:      awscertificatemanager.CertificateValidation_FromDns(props.HostedZone),
	})

	// create custom domain
	customDomainName := awsapigateway.NewDomainName(scope, id, &awsapigateway.DomainNameProps{
		Certificate: acmCert,
		DomainName:  jsii.String(customDomain),
	})

	// add api mapping to the custom domain
	pathMappingId := fmt.Sprintf("%s-%s", *id, "BasePathMapping")
	awsapigateway.NewBasePathMapping(scope, jsii.String(pathMappingId), &awsapigateway.BasePathMappingProps{
		DomainName: customDomainName,
		RestApi:    props.Api,
	})

	// add Route 53 records
	aRecordId := fmt.Sprintf("%s-%s", *id, "ARecord")
	awsroute53.NewARecord(scope, jsii.String(aRecordId), &awsroute53.ARecordProps{
		Zone:       props.HostedZone,
		RecordName: props.SubdomainName,
		Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewApiGatewayDomain(customDomainName)),
	})

	caaRecordId := fmt.Sprintf("%s-%s", *id, "CaaRecord")
	awsroute53.NewCaaRecord(scope, jsii.String(caaRecordId), &awsroute53.CaaRecordProps{
		Zone:       props.HostedZone,
		RecordName: props.SubdomainName,
		Values: &[]*awsroute53.CaaRecordValue{
			{
				Flag:  jsii.Number(0),
				Tag:   awsroute53.CaaTag_ISSUE,
				Value: jsii.String("amazon.com"),
			},
			{
				Flag:  jsii.Number(0),
				Tag:   awsroute53.CaaTag_ISSUEWILD,
				Value: jsii.String("amazon.com"),
			},
		},
	})

	services := make(map[string]constructs.Construct)
	services["Certificate"] = acmCert
	services["CustomDomain"] = customDomainName

	return services
}

func addTags(scope constructs.Construct, tagsVal map[string]string) {
	tags := awscdk.Tags_Of(scope)
	for tName, tValue := range tagsVal {
		tags.Add(jsii.String(tName), jsii.String(tValue), nil)
	}
}
