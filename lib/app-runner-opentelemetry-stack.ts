import { Stack, StackProps, RemovalPolicy, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as path from 'path';
import {
    aws_iam as iam,
    aws_dynamodb as dynamodb,
    aws_apprunner as apprunner,
} from 'aws-cdk-lib';
import { DockerImageAsset } from 'aws-cdk-lib/aws-ecr-assets';

export class AppRunnerOpentelemetryStack extends Stack {
    constructor(scope: Construct, id: string, props?: StackProps) {
        super(scope, id, props);

        const instanceRole = new iam.Role(this, 'AppRunnerInstanceRole', {
            assumedBy: new iam.ServicePrincipal(
                'tasks.apprunner.amazonaws.com'
            ),
            managedPolicies: [
                {
                    managedPolicyArn:
                        'arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess',
                },
            ],
        });

        const accessRole = new iam.Role(this, 'AppRunnerAccessRole', {
            assumedBy: new iam.ServicePrincipal(
                'build.apprunner.amazonaws.com'
            ),
            managedPolicies: [
                {
                    managedPolicyArn:
                        'arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess',
                },
            ],
        });

        const table = new dynamodb.Table(this, 'Table', {
            partitionKey: { name: 'id', type: dynamodb.AttributeType.NUMBER },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: RemovalPolicy.DESTROY,
        });
        table.grantReadWriteData(instanceRole);

        const asset = new DockerImageAsset(this, 'ContainerImage', {
            directory: path.join(__dirname, '../container'),
        });

        const observabilityConfiguration =
            new apprunner.CfnObservabilityConfiguration(
                this,
                'ObservabilityConfiguration',
                {
                    traceConfiguration: {
                        // https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-observabilityconfiguration-traceconfiguration.html
                        vendor: 'AWSXRAY',
                    },
                }
            );
        const appPort = '8080';
        const serviceAppRunner = new apprunner.CfnService(
            this,
            'AppRunnerService',
            {
                sourceConfiguration: {
                    autoDeploymentsEnabled: true,
                    authenticationConfiguration: {
                        accessRoleArn: accessRole.roleArn,
                    },
                    imageRepository: {
                        imageIdentifier: asset.imageUri,
                        imageRepositoryType: 'ECR',

                        // the properties below are optional
                        imageConfiguration: {
                            port: appPort,
                            runtimeEnvironmentVariables: [
                                {
                                    name: 'TABLE_NAME',
                                    value: table.tableName,
                                },
                                {
                                    name: 'APP_PORT',
                                    value: appPort,
                                },
                                {
                                    name: 'BACKEND',
                                    value: 'XRAY',
                                },
                            ],
                        },
                    },
                },

                healthCheckConfiguration: {
                    healthyThreshold: 1,
                    interval: 10,
                    path: '/',
                    protocol: 'HTTP',
                    timeout: 5,
                    unhealthyThreshold: 5,
                },
                instanceConfiguration: {
                    cpu: '1024',
                    instanceRoleArn: instanceRole.roleArn,
                    memory: '2048',
                },
                observabilityConfiguration: {
                    observabilityEnabled: true,
                    observabilityConfigurationArn:
                        observabilityConfiguration.attrObservabilityConfigurationArn,
                },
            }
        );

        new CfnOutput(this, 'App Runner URL', {
            value: `https://${serviceAppRunner.attrServiceUrl}`,
        });

        new CfnOutput(this, 'DynamoDB Table Name', {
            value: `${table.tableName}`,
        });
    }
}
