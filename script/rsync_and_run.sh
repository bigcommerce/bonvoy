#!/usr/bin/env bash

set -eo pipefail

die () {
    echo >&2 "$@"
    exit 1
}

CLUSTER=$1
TARGET=${2:-"auth-grpc"}
ZONE=${3:-"us-central1-f"}
PROJECT=${4:-"bc-cloud-vm"}

# Input validation
[[ -n $CLUSTER ]] || die "You must provide a cluster name"
[[ $PROJECT =~ ^(bc-cloud-vm|bc-cloud-vm-int) ]] || die "Valid project types are: bc-cloud-vm or bc-cloud-vm-int"

GCE_CLUSTER_NAME="${CLUSTER}-cloud-dev-cdvm"

# Get the IP address of a target virtual machine
IP=$(gcloud compute instances list --project="${PROJECT}" --filter="name:${GCE_CLUSTER_NAME} AND zone:${ZONE}" --format="value(INTERNAL_IP)" )
[[ -n $IP ]] || die "Unable to get the IP address for the specified machine."

# rsync the files to the target destination. Using `sudo rsync` since the code is
# owned by root on the VM.
rsync -zarh \
    -e "ssh -i ~/.ssh/google_compute_engine -o CheckHostIP=no -o StrictHostKeyChecking=no -o UserKnownHostsFile=~/.ssh/google_compute_known_hosts" \
    --rsync-path="sudo rsync" \
    --exclude '.DS_Store' \
    --exclude '.idea' \
    --exclude '.bundle' \
    --exclude '.circleci' \
    --exclude '.env' \
    --exclude '.tmp' \
    --exclude '.git' \
    --exclude 'tmp' \
    --exclude 'bonvoy' \
    ./ \
    "${USER}@${IP}:/opt/bonvoy"

gcloud compute ssh "${GCE_CLUSTER_NAME}" --zone "${ZONE}" --project "${PROJECT}" --internal-ip --ssh-flag="-A" --ssh-flag="-t" --command "sudo /opt/bonvoy/script/build-and-test.sh $TARGET"
