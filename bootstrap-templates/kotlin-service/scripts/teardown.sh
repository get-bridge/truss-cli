#!/bin/false
# ^^ failsafe: use `bash teardown.sh` or change shebang to actually tear down the service

# Feel free to delete this script from the repository!

set -euo pipefail

BASEDIR="$(cd "$(dirname ${BASH_SOURCE[0]:-$0})/.."; pwd)"

DRY_RUN=false

confirm_deletion() {
    local answer

    echo "/----------------------------------------------------------------------------------------------\\"
    echo "|   _                                                                                          |"
    echo "|  | |   This script is designed to COMPLETELY DESTROY all resources related to this service!  |"
    echo "|  | |   It WILL DELETE all AWS resources that were created by the scripts in './deploy',      |"
    echo "|  |_|   including your databases with ALL YOUR PRECIOUS DATA!                                 |"
    echo "|   _    It WILL ALSO DELETE your application from Truss and makes it unavailable for users    |"
    echo "|  (_)   and other services!                                                                   |"
    echo "|                                                                                              |"
    echo "\\----------------------------------------------------------------------------------------------/"
    echo

    echo "To confirm that you are ready to destroy the service, type the following sentence"
    echo "without spaces:"
    echo
    echo "    I understand the consequences of destroying {{.Params.name}}"
    echo

    echo "Your answer:"
    read answer

    if [[ "$answer" != "Iunderstandtheconsequencesofdestroying{{.Params.name}}" ]]; then
        echo "You didn't type the correct sentence, exiting!"
        exit 1
    fi
}

dryrun_exec() {
    if $DRY_RUN; then
        echo "exec: $@"
    else
        "$@"
    fi
}

terraform_destroy() {
    local env

    cd "$BASEDIR/deploy"

    for env in $(terraform workspace list | grep -v default | cut -c3-); do
        echo "Destroying Terraform resources for environment $env"
        dryrun_exec terraform workspace select "$env"
        echo yes | dryrun_exec terraform destroy

        dryrun_exec terraform workspace select default
        dryrun_exec terraform workspace delete "$env"
    done
}

truss_delete_tenant() {
    local env region deployment namespace tenant_count

    truss get-kubeconfig

    for deployment in {staging,prod}-{cmh,dub,syd}; do
        IFS=- read env region <<<"$deployment"
        # config="$HOME/.kube/kubeconfig-truss-$([[ "$env" == "prod" ]] && echo "prod" || echo "nonprod")-${region}"
        # namespace=$(basename $config | sed 's/kubeconfig-truss-//')

        tenant_count=$(truss wrap -e "$deployment" -- kubectl get tenant "{{.Params.name}}-${env}" -o name | wc -l)

        if [[ "$tenant_count" -gt 0 ]]; then
            dryrun_exec truss wrap -e "$deployment" -- kubectl delete tenant "{{.Params.name}}-${env}"
        fi
    done
}

truss_delete_config() {
    dryrun_exec aws s3 rm "s3://bridge-truss-config-us-east-2/{{.Params.name}}.yaml"
}


while [[ "$#" -gt 0 ]]; do
    case "$1" in
        --dry-run) DRY_RUN=true ;;
        *)
            echo "Unknown argument '$1'!"
            exit 1
    esac
    shift
done


if ! $DRY_RUN; then
    confirm_deletion
fi

cd "$BASEDIR"

terraform_destroy
truss_delete_config
truss_delete_tenant
