## References

- https://github.com/uma-arai/sbcntr-resources/blob/main/cloudformations/network_step1.yml
- https://dev.classmethod.jp/articles/interface-vpc-endpoint-cloudformation/
- https://qiita.com/pottava/items/970d7b5cda565b995fe7
- https://qiita.com/tomokyu/items/d341ba1f4a1ad1149fe4
- https://dev.classmethod.jp/articles/create-codecommit-with-code-by-cfn/
- https://dev.classmethod.jp/articles/cloudformation-supports-codebuild/
- https://gist.github.com/atheiman/cef4493821639df8192dfee7edde13af
- https://dlim716.medium.com/how-to-codebuild-with-docker-image-ca5d4389b486
- https://docs.aws.amazon.com/ja_jp/AWSCloudFormation/latest/UserGuide/aws-resource-iam-policy.html
- https://aws.amazon.com/premiumsupport/knowledge-center/cloudformation-attach-managed-policy/
- https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/creating-resources-with-cloudformation.html
- https://dev.classmethod.jp/articles/cloudformation-fargate/
- https://dev.classmethod.jp/articles/interface-vpc-endpoint-cloudformation/
- https://docs.aws.amazon.com/ja_jp/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html


## Commands

```bash
AWS_REGION_NAME=ap-northeast-1
ECR_REPOSITORY_NAME=gs-base
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
REPOSITORY_URI=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION_NAME}.amazonaws.com/${ECR_REPOSITORY_NAME}

docker image pull golang:1.19.3-alpine3.16
docker image tag golang:1.19.3-alpine3.16 ${REPOSITORY_URI}:golang1.19.3-alpine3.16

docker image pull alpine:3.16.3
docker image tag alpine:3.16.3 ${REPOSITORY_URI}:alpine3.16.3

aws ecr --region ap-northeast-1 get-login-password | docker login --username AWS --password-stdin https://${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/${ECR_REPOSITORY_NAME}

docker image push ${REPOSITORY_URI}:golang1.19.3-alpine3.16
docker image push ${REPOSITORY_URI}:alpine3.16.3

```

