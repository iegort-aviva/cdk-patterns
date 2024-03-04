package vpc

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func LookUpDefault(scope constructs.Construct) awsec2.IVpc {
	// ADD ENV logic here
	return awsec2.Vpc_FromLookup(scope, jsii.String("defaultVPC"), &awsec2.VpcLookupOptions{})
}
