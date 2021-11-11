# Creating GKE cluster & deploying

```
# show me all my projects
gcloud projects list | yank

# switch from other project to this project
gcloud config set project ___

# wizard flow for creating a new config (configs / projects relationship slightly unclear to me)
gcloud init

# show me all my "configurations"
gcloud config configurations list
gcloud config list

# make a GKE cluster
gcloud container clusters create my-blog

# show all clusters/GKE
gcloud container clusters list

# switching between configurations
gcloud config configurations activate default
gcloud config configurations activate my-blog

# ---------------
# https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl
# configuring cluster access for kubectl from gcloud

gcloud config set project ___
gcloud config configurations set my-blog

gcloud container clusters list

# need to add permissions...
# first look at existing policy
gcloud projects get-iam-policy aqueous-tube-325907 --formal=yaml

# add needed permissions (role to user)
gcloud projects add-iam-policy-binding aqueous-tube-325907 --member=user:cflynn.us@gmail.com --role=roles/container.clusterViewer

gcloud container clusters get-credentials autopilot-cluster-1 --region us-central1

# Setup for github actions deployment
# Make a service account
# https://docs.github.com/en/actions/deployment/deploying-to-your-cloud-provider/deploying-to-google-kubernetes-engine
gcloud iam service-accounts create github-actions --display-name="github-actions"

gcloud projects add-iam-policy-binding aqueous-tube-325907 \
  --member="serviceAccount:github-actions@aqueous-tube-325907.iam.gserviceaccount.com" \
  --role=roles/container.admin \
  --role=roles/storage.admin \
  --role=roles/container.clusterViewer

# grab the SA key
gcloud iam service-accounts keys create key.json --iam-account="github-actions@aqueous-tube-325907.iam.gserviceaccount.com"

# Need ClusterRoleBinding
kubectl create clusterrolebinding github-actions --clusterrole=cluster-admin --group=system:serviceaccounts

# suggested online, looks useful
kubectl create clusterrolebinding cluster-admin-binding \
  --clusterrole cluster-admin \
  --user $(gcloud config get-value account)
```

After all of the above, I used helm to install nginx as an ingress-controller (pretty simple)
https://cloud.google.com/community/tutorials/nginx-ingress-gke

Then I reserved a static IP address and updated my DNS records
```
gcloud compute addresses create saigonbros-1 \
    --global \
    --ip-version IPV4
```
