package lambda

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoDefaultFunctionProps struct {
	awscdklambdagoalpha.GoFunctionProps
}

const defaultLambdaTimeoutSeconds float64 = 29
const defaultLambdaMemorySize float64 = 128

func NewGoFunction(scope constructs.Construct, id *string, props *GoDefaultFunctionProps) awscdklambdagoalpha.GoFunction {
	p := awscdklambdagoalpha.GoFunctionProps{}
	if props != nil {
		p = props.GoFunctionProps
	}

	if err := props.Validate(); err != nil {
		panic(err)
	}

	return createLambdaGoFunction(scope, id, &p)
}

func NewGoRestApiFunction(scope constructs.Construct, id *string, props *GoDefaultFunctionProps) awscdklambdagoalpha.GoFunction {
	p := awscdklambdagoalpha.GoFunctionProps{}
	if props != nil {
		p = props.GoFunctionProps
	}

	if err := props.Validate(); err != nil {
		panic(err)
	}

	if p.Timeout == nil || *p.Timeout.ToSeconds(nil) > defaultLambdaTimeoutSeconds {
		p.Timeout = awscdk.Duration_Seconds(jsii.Number(defaultLambdaTimeoutSeconds))
	}

	return createLambdaGoFunction(scope, id, &p)
}

func (props *GoDefaultFunctionProps) Validate() error {
	if props.GoFunctionProps.Entry == nil {
		return fmt.Errorf("Lambda entry point required")
	}

	return nil
}

func createLambdaGoFunction(scope constructs.Construct, id *string, props *awscdklambdagoalpha.GoFunctionProps) awscdklambdagoalpha.GoFunction {
	p := awscdklambdagoalpha.GoFunctionProps{}
	if props != nil {
		p = *props
	}

	if p.Runtime == nil {
		p.Runtime = awslambda.Runtime_PROVIDED_AL2()
	}

	if p.Architecture == nil {
		p.Architecture = awslambda.Architecture_ARM_64()
	}

	if p.Timeout == nil {
		p.Timeout = awscdk.Duration_Seconds(jsii.Number(defaultLambdaTimeoutSeconds))
	}

	if p.MemorySize == nil {
		p.MemorySize = jsii.Number(defaultLambdaMemorySize)
	}

	if p.Bundling == nil {
		ldflags := fmt.Sprintf(`-ldflags "%s"`, "-s -w")
		p.Bundling = &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: &[]*string{jsii.String(ldflags)},
		}
	}

	return awscdklambdagoalpha.NewGoFunction(scope, id, &p)
}
