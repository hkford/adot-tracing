#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { AppRunnerOpentelemetryStack } from '../lib/app-runner-opentelemetry-stack';

const app = new cdk.App();
new AppRunnerOpentelemetryStack(app, 'AppRunnerOpentelemetryStack', {
    env: { region: 'us-east-1' },
});
