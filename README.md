
![](https://github.com/akiirobot/sign-storage-url/workflows/DeployToGoogleCloud/badge.svg)

# Sign Storage Url API

This is a micro service that implement the [V4 signing process with Cloud Storage tools](https://cloud.google.com/storage/docs/access-control/signing-urls-with-helpers#storage-signed-url-get-object-go) on Google Cloud Function. You don't need to write the signing process in you program. Just call this Sign Storage Url API after you deploy to google cloud function.

## Before Start

### Environment Variable

| Environment  Varibale | Required | Description |
|-----------------------|----------|---|
| SERVICE_JSON_FILE     | required | This is a key for signing storage. see [Get Signed Url Key](#Get Signed Url Key) |
| BUCKET_NAME           | optional | Add the multiple buckets separated by colon **:** to your white lists or leave blank for all buckets. |
| EXPIRE_TIME           | optional | How much time (minute) the sign url is valid. 15 minutes is default value. |

### Setup Google Cloud Project

Create a Google Cloud Project from [here](https://console.cloud.google.com/projectcreate) and set on terminal.

```shell
export ProjectID=<project_id>

# Set Project ID
gcloud config set project ${ProjectID}
```

### Create Google Cloud Storage If Needed

You can change the storage location and default storage class.

```shell
export BucketName=<bucket_name>

# Create Bucket
gsutil mb -b on -c Standard -p ${ProjectID} -l asia gs://${BucketName}
```

### Get Signed Url Key

After running the commands below, you will get the **signed-url-key.json** file under the current folder. You will need it on the next step.

```shell
# Create Service Account
gcloud iam service-accounts create "signed-url" --display-name "signed-url"

# Grant Service Account with storage object admin
gcloud projects add-iam-policy-binding ${ProjectID} \
  --member serviceAccount:signed-url@${ProjectID}.iam.gserviceaccount.com \
  --role roles/storage.objectAdmin

# Create Key
gcloud iam service-accounts keys create signed-url-key.json --iam-account signed-url@${ProjectID}.iam.gserviceaccount.com
```

## Deploy to Google cloud Function

### Allow All Storage Bucket

```shell
gcloud functions deploy sign --entry-point=SignedUrl --runtime=go111 --trigger-http --quiet \
  --set-env-vars SERVICE_JSON_FILE=signed-url-key.json
```

**Testing**

You can get the signed url from running the following command

```shell
curl -k -X POST -F "bucket=<bucket-name>" -F "method=POST" -F "object=hello.txt" https://<gcloud-function-url>/sign
```

### Limit Multiple Storage Bucket

If you want to contraint which bucket is allowed to access, adding the **BUCKET_NAME** variable spliting each bucket with colon **:**

```shell
# multiple buckets
gcloud functions deploy sign --entry-point=SignedUrl --runtime=go111 --trigger-http --quiet \
  --set-env-vars SERVICE_JSON_FILE=signed-url-key.json,BUCKET_NAME=<bucket-name1>:<bucket-name2>
```

**Testing**

You can get the signed url from running the following command

```shell
curl -k -X POST -F "bucket=<bucket-name>" -F "method=POST" -F "object=hello.txt" https://<gcloud-function-url>/sign
```

### Limit Single Storage Bucket

```shell
gcloud functions deploy sign --entry-point=SignedUrl --runtime=go111 --trigger-http --quiet \
  --set-env-vars SERVICE_JSON_FILE=signed-url-key.json,BUCKET_NAME=<bucket-name1>
```

**Testing**

You can ignore bucket form data if you depoly a single storage bucket

```shell
curl -k -X POST -F "method=POST" -F "object=hello.txt" https://<gcloud-function-url>/sign
```

## Create CI/CD Deploy Key (Optional)

This key can be used in CI/CD.

Grant the Required permission:
- **roles/iam.serviceAccountUser**: Can act as the service account to start a Google Cloud Function.
- **roles/cloudfunctions.developer**: Deploy to Google Cloud Function
- **roles/storage.admin**: Create or delete a storage, and sign the storage url

```shell
# Create Service Account
gcloud iam service-accounts create "deploy-app" --display-name "deploy-app"

# Grant Service Account with cloud functions developer for deploy
gcloud projects add-iam-policy-binding ${ProjectID} \
  --member serviceAccount:deploy-app@${ProjectID}.iam.gserviceaccount.com \
  --role roles/iam.serviceAccountUser
gcloud projects add-iam-policy-binding ${ProjectID} \
  --member serviceAccount:deploy-app@${ProjectID}.iam.gserviceaccount.com \
  --role roles/cloudfunctions.developer
gcloud projects add-iam-policy-binding ${ProjectID} \
  --member serviceAccount:deploy-app@${ProjectID}.iam.gserviceaccount.com \
  --role roles/storage.admin

# Create Key
gcloud iam service-accounts keys create deploy-key.json --iam-account deploy-app@${ProjectID}.iam.gserviceaccount.com
```

