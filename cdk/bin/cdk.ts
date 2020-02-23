#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { ECSStack } from '../lib/ecs-stack';

const app = new cdk.App();
new ECSStack(app, 'ECSBoilerplateSampleStack', {
    /*
    // If vpcAttributes is not specified, new VPC is created.
    vpcAttributes: {
        vpcId: '',
        availabilityZones: [],
        publicSubnetIds: [],
        privateSubnetIds: [],
    },

    // DNS record. Even if this is not specified, you can access with ELB domain (***.elb.amazonaws.com)
    route53: {
        zoneId: '',
        zoneName: 'example.com',
        recordName: 'foo',
    },
    // Certificate Manager ARN. Required if accessing with HTTPS
    acmArn: 'arn:aws:acm:****'

    // default values
    containerPort: 8080,
    cpu: 256,
    memoryLimitMiB: 512,
    minCapacity: 1,
    maxCapacity: 5,
    scaleCPUPercent: 80
    */
});
