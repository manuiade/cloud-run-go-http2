# Bypass Cloud Run 32MB upload/download limit using Go web server with HTTP2

Go HTTP2 web server implementation to deploy to Cloud Run to bypass its 32 MB upload/download limit

## Env vars
```
PROJECT_ID=<PROJECT_ID>
REPO=go-http2
REGION=europe-west1
gcloud config set project $PROJECT_ID
```

## Create a docker repository on artifact registry

```
gcloud artifacts repositories create $REPO \
	--repository-format=docker \
	--location=$REGION

gcloud auth configure-docker $REGION-docker.pkg.dev
```

## Build repository locally with docker (tag with container registry image) and push image

```
docker build --tag $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$REPO:v1.0 src/
docker push $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$REPO:v1.0
```

## Deploy Cloud Run instance

```
gcloud run deploy $REPO \
	--region $REGION \
	--allow-unauthenticated \
    --port 3000 \
    --use-http2 \
	--cpu 4 \
	--memory 16Gi \
    --image $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$REPO:v1.0
```

## Generate large file to upload/download

```
dd if=/dev/urandom of=./large.txt bs=1000MB count=1
```

## Test on Cloud Run
```
RUN_URL=$(gcloud run services describe $REPO --format json | jq -r '.status.url')

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" --http2 "$RUN_URL"
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" --http2 -F "file=@large.txt" "$RUN_URL" --output output.txt
```


## Cleanup

```
rm large.txt output.txt

gcloud run services delete $REPO --region $REGION --quiet

gcloud artifacts repositories delete $REPO \
	--location=$REGION --quiet
```