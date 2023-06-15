#!/bin/bash
# set -exo pipefail

# Source: https://github.com/gotugboat/tugboat/blob/main/scripts/install.sh

# Default values
BINARY_NAME=${BINARY_NAME:-"tugboat"}
INSTALL_DIR=${INSTALL_DIR:-"/usr/local/bin"}
TUGBOAT_VERSION=${TUGBOAT_VERSION:-"latest"}
REPO=${REPO:-"gotugboat/tugboat"}
DEBUG=${DEBUG:-false}
DRY_RUN=${DRY_RUN:-false}
VERIFY_CHECKSUM=${VERIFY_CHECKSUM:-true}
TMP_ROOT="./tmp"

usage() {
  echo -e "An installation script for the tugboat cli \n\n"
  echo "Usage:"
  echo "${0} [OPTIONS]"
  echo ""
  echo "Options:"
  echo "-h, --help             Show this help message"
  echo "    --debug            Output more information about the execution"
  echo "    --dry-run          Print out what will happen, do not execute"
  echo "-v, --version          Choose the desired version to install (default: ${TUGBOAT_VERSION})"
  echo ""
  echo "Example:"
  echo "${0} --version 1.0.0"
  echo ""
  echo "For more information please visit the documentation at https://gotugboat.io/docs/getting-started/introduction/"
  echo ""
}

setup_color() {
	# Only use colors if connected to a terminal
  if [ -t 1 ]; then
    RED=$(printf '\033[31m')
    GREEN=$(printf '\033[32m')
    YELLOW=$(printf '\033[33m')
    BLUE=$(printf '\033[34m')
    BOLD=$(printf '\033[1m')
    DIM=$(printf '\033[2m')
    UNDER=$(printf '\033[4m')
    RESET=$(printf '\033[m')
  fi
}

# Logging functions
log_debug() {
  if [[ ${DEBUG} == true ]]; then
    printf "${DIM}DEBU: %s${RESET}\n" "$*"
  fi
}

log_info() {
  printf "${BLUE}INFO:${RESET} %s\n" "$*"
}

log_warn() {
  printf "${YELLOW}WARN:${RESET} %s\n" "$*"
}

log_err() {
  printf "${RED}ERRO:${RESET} %s\n" "$*"
}

# Check if the command exists in the system's list of commands
check_command() {
  if command -v "$1" >/dev/null 2>&1; then
    return 0
  else
    return 1
  fi
}

# Gets the canonical name for the system architecture
get_arch() {
  arch=""
  case "$(arch)" in
    amd64|x86_64)    arch='amd64' ;;
    aarch64|arm64)   arch='arm64' ;;
    armhf|armv7l)    arch='arm' ;;
    *) echo "Unsupported architecture $(arch)" ; exit 1 ;;
  esac
  echo "${arch}"
}

# Gets the operating system name (i.e. windows, linux, darwin)
get_os() {
  os=$(echo `uname` | tr '[:upper:]' '[:lower:]')
  case "${os}" in
    # Minimalist GNU for Windows
    mingw*|cygwin*) os='windows';;
  esac
  echo "${os}"
}

# Returns the latest release from GitHub
check_latest_release() {
  local git_project=$1
  if check_command curl; then
    local latest_release=$(curl -L -s -H 'Accept: application/json' https://github.com/${git_project}/releases/latest)
  fi
  local release_tag=$(echo $latest_release | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
  echo $(echo ${release_tag} 2> /dev/null | sed 's/^.//')
}

check_requirements() {
  log_info "Checking requirements"

  if ! check_command curl; then
    log_err "curl is required to proceed with the installation"
    exit 1
  fi
}

download_binary() {
  local version=$1
  local os="$(get_os)"
  local arch="$(get_arch)"

  log_debug "os: ${os}"
  log_debug "arch: ${arch}"
  if [[ "${arch}" == *"Unsupported architecture"* ]]; then
    log_err ${arch}
    exit 1
  fi

  local github_url="https://github.com/${REPO}/releases/download/v${version}"
  local download_url="${github_url}/${BINARY_NAME}-${os}-${arch}.tar.gz"
  local checksum_url="${download_url}.sha256sum"

  log_info "Downloading ${download_url}"

  if [[ "${DRY_RUN}" == "true" ]]; then
    return
  fi

  # create a temp directory (removed with fail_trap)
  TMP_ROOT="$(mktemp -dt ${BINARY_NAME}-installer-XXXXXX)"
  log_debug "tmp dir: ${TMP_ROOT}"
  local download_file="${TMP_ROOT}/${BINARY_NAME}.tar.gz"
  local checksum_file="${download_file}.sha256sum"
  
  if check_command curl ; then
    curl -f -SsLo "${download_file}" "${download_url}"
    if [[ "$?" == "22" ]]; then
      log_err "Download failed: there is likely not a release for ${os}-${arch}"
      exit 1
    fi
    curl -SsLo "${checksum_file}" "${checksum_url}"
  fi
}

verify_checksum() {
  local download_file="${TMP_ROOT}/${BINARY_NAME}.tar.gz"
  local checksum_file="${download_file}.sha256sum"
  
  log_info "Verifying the checksum"

  if [[ "${DRY_RUN}" == "true" ]]; then
    log_debug "echo ${checksum_file} ${download_file} | sha256sum --check"
    return
  fi

  if ! echo "$(<${checksum_file}) ${download_file}" | sha256sum --check >/dev/null 2>&1 ; then
    log_err "The checksum file did not match, aborting"
    exit 1
  fi
}

verify_file() {
  if [ "${VERIFY_CHECKSUM}" == "true" ]; then
    verify_checksum
  fi
}

install_file() {
  local binary_tmp_location="${TMP_ROOT}/${BINARY_NAME}"
  local binary_location="${binary_tmp_location}/${BINARY_NAME}"
  local binary_install_location="${INSTALL_DIR}/${BINARY_NAME}"

  log_info "Installing ${BINARY_NAME} to ${INSTALL_DIR}"
  log_debug "Extracting download to ${TMP_ROOT}"

  if [[ "${DRY_RUN}" == "true" ]]; then
    log_debug "tar -xzf ${binary_tmp_location}.tar.gz -C ${TMP_ROOT}"
    log_debug "cp ${binary_location} ${binary_install_location}"
    return
  fi
  
  if ! tar -xzf "${binary_tmp_location}.tar.gz" -C "${TMP_ROOT}"; then
    log_err "Extracting archive failed"
    exit 1
  fi

  log_debug "cp ${binary_location} ${binary_install_location}"
  cp "${binary_location}" "${binary_install_location}" > /dev/null 2>&1 || sudo cp "${binary_location}" "${binary_install_location}"
}

install_binary() {
  if [[ "${TUGBOAT_VERSION}" == "latest" ]]; then
    install_version=$(check_latest_release ${REPO})
  else
    install_version="${TUGBOAT_VERSION}"
  fi

  log_debug "Installing version: ${install_version}"
  download_binary "${install_version}"
  verify_file
  install_file
}

verify_installation() {
  log_info "Verifying the installation"

  if [[ "${DRY_RUN}" == "true" ]]; then
    return
  fi

  if ! check_command tugboat; then
    log_err "${BINARY_NAME} not found. Is ${INSTALL_DIR} on your "'$PATH?'
    exit 1
  fi
}

cleanup() {
  if [[ -d "${TMP_ROOT}" ]]; then
    log_debug "removing: ${TMP_ROOT}"
    rm -rf "$TMP_ROOT"
  fi
}

fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    log_err "Failed to install $BINARY_NAME"
  fi
  cleanup
  exit $result
}

main() {
  # Parse options
  while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
      -h|--help)
      usage
      exit 0
      ;;
      --debug)
      DEBUG=true
      ;;
      --dry-run)
      DRY_RUN=true
      ;;
      -v|--version)
      TUGBOAT_VERSION=$2
      shift
      ;;
      *)
      echo "Unknown option: $key"
      usage
      exit 1
      ;;
    esac
    shift
  done

  setup_color

  if [[ "${DRY_RUN}" == "true" ]]; then
    log_warn "Dry run in progress"
  fi

  check_requirements
  install_binary
  verify_installation

  if [[ "${DRY_RUN}" != "true" ]]; then
    printf "\n${GREEN}%s${RESET}\n\n" "${BINARY_NAME} has been installed successfully!"
  fi
}

# Stop execution on any error
trap "fail_trap" EXIT

main "$@"
