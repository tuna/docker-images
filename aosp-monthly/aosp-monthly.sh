#!/bin/bash
set -e

REPO=${REPO:-"/usr/local/bin/repo"}
AOSP_DIR=${AOSP_DIR:-"/data/aosp"}
PKG_DIR=${PKG_DIR:-"/data/package"}
export REPO_URL="https://mirrors.tuna.tsinghua.edu.cn/git/git-repo/"

function repo_init() {
	mkdir -p $AOSP_DIR
	cd $AOSP_DIR
	$REPO init -u http://aosp.tuna.tsinghua.edu.cn/platform/manifest
}

function repo_sync() {
	cd $AOSP_DIR
	$REPO sync -f --network-only --no-clone-bundle
}

function package() {
	mkdir -p $PKG_DIR
	cd $PKG_DIR
	pkg="aosp-`date +%Y%m%d`.tar"
	base_dir=`dirname ${AOSP_DIR}`
	aosp_dir=`basename ${AOSP_DIR}`
	echo "Creating $pkg"
	tar cf $pkg -C $base_dir $aosp_dir
	echo "Calculating Checksum"
	md5sum $pkg > "${pkg}.md5"
	ln -sf $pkg aosp-latest.tar
	ln -sf "${pkg}.md5" aosp-latest.tar.md5
}

if [[ ! -d "$AOSP_DIR/.repo" ]]; then
	echo "Initializing AOSP working directory"
	repo_init
fi

repo_sync
package
