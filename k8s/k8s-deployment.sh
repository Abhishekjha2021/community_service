#!/bin/sh
kubectl apply -f $NAMESPACE/$BITBUCKET_REPO_SLUG/eks/deployment.yaml
echo $IMAGE_TAG
kubectl set image deployment/$DEPLOYMENT_NAME $K8S_CONTAINER=$ECR_REPO_URI/$AWS_ECR_REPOSITORY:$IMAGE_TAG -n $NAMESPACE;
kubectl rollout restart deployment $DEPLOYMENT_NAME -n $NAMESPACE;
kubectl set env deployment/$DEPLOYMENT_NAME --env="CONFIG_FILE=$CONFIG_FILE" -n $NAMESPACE;