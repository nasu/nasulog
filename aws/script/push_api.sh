#/bin/sh
set -ex

#TODO: validate vars
AWS_ACCOUNT=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=$(aws configure get region)
ECR_HOST=${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com
REPOSITORY=nasulog
TAG=${1:-latest}

docker build -t $ECR_HOST/$REPOSITORY:$TAG ../api/
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_HOST
docker push $ECR_HOST/$REPOSITORY:$TAG