# Repo URL comes from:
# terraform-base$ terraform output -json|jq .ecr_repository.value.repository_url
repo="${repo:-465509910691.dkr.ecr.eu-central-1.amazonaws.com/echoserver}"
aws --profile priv-master ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin "$repo"
#docker build -t echoserver .
docker tag echoserver:latest "$repo":latest
docker push "$repo":latest
