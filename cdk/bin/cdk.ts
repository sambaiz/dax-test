#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { ECSStack } from '../lib/ecs-stack';

const app = new cdk.App();
new ECSStack(app, 'DaxTestStack', {
    // If vpcAttributes is not specified, new VPC is created.
    vpcAttributes: {
        vpcId: 'vpc-95d62ef0',
        availabilityZones: ['ap-northeast-1d', 'ap-northeast-1a', 'ap-northeast-1c'],
        publicSubnetIds: ['subnet-0aeba822', 'subnet-1ffb3a68', 'subnet-e8b356b1'],
        privateSubnetIds: ['subnet-0aeba822', 'subnet-1ffb3a68', 'subnet-e8b356b1'],
    },
    /*
    // DNS record. Even if this is not specified, you can access with ELB domain (***.elb.amazonaws.com)
    route53: {
        zoneId: '',
        zoneName: 'example.com',
        recordName: 'foo',
    },
    // Certificate Manager ARN. Required if accessing with HTTPS
    acmArn: 'arn:aws:acm:****'
    */
    // default values
    containerPort: 8080,
    cpu: 256,
    memoryLimitMiB: 2048,
    minCapacity: 1,
    maxCapacity: 10,
    scaleCPUPercent: 80
});
