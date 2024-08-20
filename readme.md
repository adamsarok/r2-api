#sops:
gcloud auth login
#on new installation: gcloud auth application-default login
sops decrypt .\config.enc.yml > .\config.yml
sops encrypt --gcp-kms projects/secretmanager-433009/locations/global/keyRings/sops/cryptoKeys/sops-key config.yml > config.enc.yml

docker build -t fuzzydice555/r2-api-go .
docker push fuzzydice555/r2-api-go


Docker compose:
services:
    perfume-tracker:
        image: fuzzydice555/r2-api-go
        ports:
          - 9088:8080
        environment:
          - R2_ENDPOINT=...
          - R2_BUCKET=test
          - R2_REGION=auto
          - R2_ACCESS_KEY=...
          - R2_SECRET_KEY=...
          - R2_UPLOAD_EXPIRY_MINUTES=30
          - R2_DOWNLOAD_EXPIRY_MINUTES=30
        restart: unless-stopped