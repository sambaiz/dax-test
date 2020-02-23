## ecs-boilerplate

Boilerplate to create ECS (Fargate) resources and deploy by AWS CDK

### Get started

1. Edit `cdk/bin/cdk.ts` 

```
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
```

2. Develop app
3. Deploy

```
$ AWS_PROFILE make deploy
...
Outputs:
ECSBoilerplateSampleStack.ServiceServiceURL*** = http://***.elb.amazonaws.com

$ curl http://***.elb.amazonaws.com
Hi this is ECS
```