#!/bin/bash

#
# Build and push the docker image
#

Error() {
    echo "error: ${1}"
    exit 1
}

script_path=$(dirname "${0}")
docker_namespace='kisiodigital'
image='stomptherabbit'

# handle options
# default values
dry_run='false'
while [ ${#} -gt 0 ]; do
    case "${1}" in
    '-d' | '--dry-run') dry_run='true' ;;
    *) Error "flag ${1} is unknown, the only recognized flag is '--dry-run (-d)'" ;;
    esac
    shift 1
done

[ "${dry_run}" = 'true' ] && echo "dry-run activated! The docker image will not be pushed."

# build only on a clean git status
[ -n "$(git status --untracked-files=no --porcelain)" ] && Error "git status is not clean, build aborted"

# step 0: prepare docker image name with registry and tag
# get the tag from git
tag=$(git describe --tags --abbrev=0 2> /dev/null)
has_tag=$(git tag --list --points-at HEAD)
[ -z "${has_tag}" ] && tag='latest'
[ -z "${tag}" ] && Error "impossible to get a tag"


# step 1: build docker images for git HEAD
image_fullname="${docker_namespace}/${image}:${tag}"
echo "Building $image_fullname"
docker build --force-rm -t "${image_fullname}" "${script_path}"/../..

# step 2: push the image to the registry
if [ "${dry_run}" = 'false' ]; then
    echo "pushing '${image_fullname}'"
    docker push "${image_fullname}" || Error "pushing ${image} failed"
fi

exit 0
