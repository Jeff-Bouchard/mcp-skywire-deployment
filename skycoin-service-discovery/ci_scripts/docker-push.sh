#!/usr/bin/env bash

function print_usage() {
  echo "Use: $0 [-t <docker_image_tag_name>] [-p | -b]"
  echo "use -p for push (it builds and push the image)"
  echo "use -b for build image locally"
}

if [[ $# -ne 3 ]]; then
  print_usage
  exit 0
fi

function docker_build() {
	# following 2 lines are needed for SD to function. these can't be run from within the container and need to be run on the host machine
  if [ $(uname -s) == "Linux" ]; then \
		sudo bash -c 'echo 1 > /proc/sys/net/ipv4/ip_forward' && \
		sudo bash -c 'echo 1 > /proc/sys/net/ipv6/conf/all/forwarding'; \
	fi
  docker image build \
    --tag=skycoinpro/service-discovery:"$tag" \
     .
	# revert the changes
	if [ $(uname -s) == "Linux" ]; then \
		sudo bash -c 'echo 0 > /proc/sys/net/ipv4/ip_forward' && \
		sudo bash -c 'echo 0 > /proc/sys/net/ipv6/conf/all/forwarding'; \
	fi
}

function docker_push() {
  docker tag skycoinpro/service-discovery:"$tag" skycoinpro/service-discovery:"$tag"
  docker image push skycoinpro/service-discovery:"$tag"
}

while getopts ":t:pb" o; do
  case "${o}" in
  t)
    tag="$(echo "${OPTARG}" | tr -d '[:space:]')"
    if [[ $tag == "develop" ]]; then
      tag="test"
    elif [[ $tag == "master" ]]; then
      tag="prod"
    else
      tag=$tag
    fi
    ;;
  p)
    docker_build
    docker_push
    ;;
  b)
    docker_build
    ;;
  *)
    print_usage
    ;;
  esac
done
