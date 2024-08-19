#sops:
gcloud auth login
sops decrypt .\config.enc.yml > .\config.yml
sops encrypt --gcp-kms projects/secretmanager-433009/locations/global/keyRings/sops/cryptoKeys/sops-key config.yml > config.enc.yml
