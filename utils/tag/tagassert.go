package tag

import (
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type UseCaseTagLabel string

const companyName = "test"
const (
	devOpsUseCaseTagLabel UseCaseTagLabel = "operations"
	// ??: Should it be different tags for `default` and `patterns`?
	devOpsUseCaseConstructTagLabel UseCaseTagLabel = "operations:construct"
)

type TagProps struct {
	awscdk.TagProps
	TagName  *string
	TagValue *string
}

func (props *TagProps) Validate() error {
	return nil
}

func AddDevOpsStackName(scope constructs.Construct, props *TagProps) {
	p := awscdk.TagProps{}
	if props != nil {
		p = props.TagProps
	}

	name := createTagName(devOpsUseCaseTagLabel, "stack-name")
	awscdk.Tags_Of(scope).Add(jsii.String(name), props.TagValue, &p)
}

func AddDevOpsConstructName(scope constructs.Construct, props *TagProps) {
	p := awscdk.TagProps{}
	if props != nil {
		p = props.TagProps
	}

	name := createTagName(devOpsUseCaseConstructTagLabel, *props.TagName)
	awscdk.Tags_Of(scope).Add(jsii.String(name), props.TagValue, &p)
}

func createTagName(ucTagLevel UseCaseTagLabel, name string) string {
	els := []string{companyName, string(ucTagLevel), name}
	return strings.Join(els, ":")
}
