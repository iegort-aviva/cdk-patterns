[<- Home](./../../README.md)

[TBD]

1. Naming convention:

2. Tracking usage:

To facilitate construct usage tracking, the implementation includes appending a tag with the key `<company-name>:operations:construct:<tagName>Version` to denote the specific version `<version>` being used. This tagging strategy is integral for monitoring and managing the lifecycle of constructs within our infrastructure, ensuring compliance with internal standards and enabling efficient version control and upgrade paths.
[impl](../../utils/tag/)

3. Implementation:

There are three strategies for compliance implementation in AWS CDK:
Implementing proactive compliance strategies in AWS CDK involves creating mechanisms to automatically modify non-compliant resources, throw errors for non-compliance issues, and use internal, customized abstractions to enforce compliance. These strategies ensure that your cloud infrastructure not only meets current compliance standards but is also positioned to adapt to future changes in those standards. Let's explore each strategy in detail:

> // Discuss which approach when to use

### 1. Proactively Modifying Non-Compliant Resources

This strategy involves creating CDK constructs or utilizing CDK aspects that automatically adjust resource configurations to meet compliance standards during the synthesis phase, before the resources are deployed.

**Implementation Steps:**

- **Utilize CDK Aspects:** Implement aspects in AWS CDK, which are a way to visit and modify the construct tree. Use them to inspect and automatically modify non-compliant resource properties to align with compliance policies.

- **Custom Compliance Constructs:** Develop custom constructs that encapsulate compliant configurations for common resources. These constructs should include logic to automatically adjust parameters that do not meet predefined compliance standards.

**Advantages:**

- **Automated Compliance Enforcement:** Automatically adjusts resource configurations to ensure compliance, reducing manual review and intervention. This ensures that deployments are compliant by default, significantly reducing the risk of compliance violations.
- **Increased Efficiency:** By automating compliance adjustments, development teams can focus more on building features rather than worrying about compliance details. This leads to faster development cycles and quicker time-to-market for new features or services.
- **Scalability:** As your infrastructure grows, maintaining compliance across numerous resources can become challenging. This approach allows for scalable compliance enforcement, applying consistent standards across all resources automatically.

**Disadvantages:**

- **Complexity in Implementation:** Developing aspects or custom constructs that automatically adjust resource configurations can be complex and requires deep understanding of both AWS CDK and the specific compliance requirements.
- **Risk of Overriding Intentional Configurations:** There's a risk that automatic modifications could override configurations intentionally set by developers, potentially leading to unexpected behaviors or configurations that do not align with the application's requirements.
- **Maintenance Overhead:** As compliance standards evolve, these custom constructs and aspects need to be regularly updated, which can become a significant maintenance burden.

### 2. Throwing Errors for Non-Compliance Issues

This approach involves implementing checks within the CDK deployment pipeline that validate resource configurations against compliance standards. If a non-compliance issue is detected, the process throws an error, halting the deployment.

**Implementation Steps:**

- **Pre-Deployment Compliance Checks:** Integrate compliance checking tools or scripts in the CI/CD pipeline that analyze the generated CloudFormation templates or CDK constructs against compliance rules.

- **Custom Validation Functions:** Within custom constructs or the CDK app, implement validation functions that throw errors if the resource configurations do not meet compliance standards.

**Advantages:**

- **Immediate Feedback:** Provides developers with immediate feedback on compliance violations, allowing for quick corrections before deployment. This helps in maintaining a high standard of compliance and security from the early stages of development.
- **Prevention of Non-Compliant Deployments:** By halting deployments that do not meet compliance standards, this strategy ensures that only compliant resources are provisioned. This significantly reduces the risk of introducing vulnerabilities or non-compliant resources into the production environment.
- **Enhanced Governance:** Throwing errors for non-compliance issues supports strong governance practices by enforcing compliance checks as part of the CI/CD pipeline. This creates a clear audit trail for compliance efforts and helps in demonstrating adherence to regulatory standards during audits.

**Disadvantages:**

- **Potential for Deployment Delays:** Relying on error-throwing for compliance issues can lead to frequent deployment halts, slowing down development and deployment cycles, especially in environments where rapid iteration is critical.
- **Increased Developer Friction:** This strategy may increase friction among developers, who might find it frustrating to repeatedly address compliance errors before proceeding with deployments, potentially impacting morale and productivity.
- **Dependency on Comprehensive Rule Sets:** The effectiveness of this approach depends on having a comprehensive set of compliance rules. Developing and maintaining these rules requires continuous effort and expertise in compliance standards.

### 3. Internal, Customized Abstractions for Compliance

Creating internal, customized abstractions means developing higher-level constructs or libraries that encapsulate the organization's compliance logic, making it easier to build compliant infrastructure without deep expertise in specific compliance requirements.

**Implementation Steps:**

- **Develop Compliance Libraries:** Create libraries of custom constructs that encapsulate compliance requirements, such as specific encryption standards, networking configurations, or logging and monitoring setups.

- **Incorporate Compliance into Base Constructs:** For commonly used AWS services, create base constructs that include all necessary compliance configurations. These constructs serve as the foundation for any deployment of that service within the organization.

**Advantages:**

- **Streamlined Compliance Management:** Customized abstractions encapsulate compliance logic, making it easier for developers to deploy compliant resources without needing in-depth knowledge of specific compliance requirements. This streamlines the process of managing compliance across projects.
- **Consistency Across Deployments:** Using internal abstractions ensures that all deployments are consistent with the organization’s compliance policies. This uniformity helps in managing and auditing cloud resources more effectively.
- **Facilitates Reusability and Collaboration:** Developing a library of compliant constructs encourages reuse and collaboration within the organization. Teams can share and contribute to a centralized repository of compliance-focused resources, fostering a culture of security and compliance.

**Disadvantages:**

- **Limited Flexibility:** While internal abstractions can simplify compliance, they may also limit developers' flexibility in using AWS services, as they are constrained to the configurations defined within these abstractions.
- **Knowledge Silos:** Relying heavily on customized abstractions can create knowledge silos, where only a few team members understand the compliance logic encapsulated within these constructs. This can make troubleshooting and modifications challenging.
- **Initial Development Time:** Developing a comprehensive library of customized compliance constructs requires a significant upfront investment in time and resources. This initial development phase can delay other projects or lead to prioritization conflicts.

4. Standard Guidelines for Construct Function Signatures and Property Definitions in AWS CDK:

In the documentation for our custom CDK constructs, we adhere to a standardized approach for defining construct constructors and properties to ensure consistency, maintainability, and ease of use across our infrastructure as code (IaC) projects. Below is the guideline for writing constructors and property types for our custom constructs.

### Constructor Signature

For every construct where applicable, we aim to maintain a consistent function builder signature as follows:

```go
func NewFoo(scope constructs.Construct, id *string, props *CustomProps)
```

This signature comprises three parameters:

- `scope constructs.Construct`: The scope in which this construct is defined. This is a mandatory parameter and acts as the context in which the current construct will operate.
- `id *string`: A unique identifier for the construct, facilitating the differentiation of instances when multiple constructs of the same type are used within the same scope.
- `props *CustomProps`: A pointer to a custom properties `struct` that defines specific configuration options for the construct. This is the last parameter and is designed to be flexible to accommodate various customization needs.

### Custom Properties

When defining properties for a construct, the structure of `CustomProps` varies based on the type of construct being created:

#### Layer 2 Construct Properties

For a `layer 2` (default) construct, which typically enhances or modifies a single AWS resource, the properties should be defined as follows:

```go
type CustomProps struct {
    awscdkmodule.OriginalProps   // Embedding the original properties from the AWS CDK module
    CustomKey *type              // Define additional pointer to custom properties as needed
    // Add other custom fields here
}
```

This approach allows for the inclusion of all original properties from the AWS CDK module for the resource being wrapped or extended, alongside any additional custom properties specific to our organization's needs.

#### Layer 3 Construct Properties

For `layer 3` (pattern) constructs, which compose multiple AWS resources to implement a higher-level infrastructure pattern, the properties should be structured as:

```go
type CustomProps struct {
    Api        awsapigateway.IRestApi | *awsapigateway.LambdaRestApiProps // Interface or properties for API Gateway
    HostedZone awsroute53.IHostedZone | *awsroute53.HostedZoneProps       // Interface or properties for Route 53 Hosted Zone
    CustomKey  type                                                       // Additional custom properties as required
    // Include other necessary fields here
}
```

This structure supports the composition of various AWS resources, providing the flexibility to pass either direct references to existing resources (interfaces) or configurations for creating new ones (properties).

### Validation Method

If additional validation of properties is required, a validation method should be defined within the construct. This method ensures that the provided properties meet specific criteria before proceeding with the construct's deployment:

```go
func (props *CustomProps) Validation() error {
    // Implement validation logic here
    // Return nil if validation passes or an error if it fails
}
```

This method allows for the early detection and reporting of configuration errors, improving the reliability and predictability of deployments.

By following these guidelines for construct documentation, we aim to promote clarity, consistency, and best practices in the development and use of custom CDK constructs within our organization.

4. For the README.md file of each GoLang package within your AWS CDK project, consider structuring the documentation as follows to ensure clarity, comprehensiveness, and ease of use for developers. This structure aligns with best practices for documenting software projects, as exemplified in resources like "The CDK Book".

### Documentation Structure for README.md for each package

#### High-Level Overview

Begin with a concise summary that explains the purpose and functionality of the package. Describe what the package does, its target use case, and any specific problems it solves within the context of AWS CDK projects.

#### Architecture Diagram [CDK-Dia](https://github.com/pistazie/cdk-dia)

Include an architecture diagram that visually represents how the CDK constructs within the package interact with AWS services and with each other. Ensure the diagram is clear and annotated to help users understand the flow and structure at a glance. If the package documentation is in a format that supports images, embed the diagram directly; otherwise, provide a link to where the diagram can be viewed.

#### Detailed Feature List

Provide a detailed list of features offered by the package. This should include:

- Custom constructs and their purposes.
- Integration capabilities with AWS services.
- Any utilities or helper functions included in the package.
- Configuration options and customization points for users.

#### Sample Implementation Code

Include sample code snippets or examples that demonstrate how to use the constructs and features provided by the package. These examples should be practical and straightforward, showing how to instantiate constructs, configure properties, and integrate with other AWS services or constructs where applicable.
