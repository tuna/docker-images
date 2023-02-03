#!/bin/bash
set -e

REPO=${REPO:-"/usr/local/bin/repo"}
AOSP_DIR=${AOSP_DIR:-"/data/aosp"}
PKG_DIR=${PKG_DIR:-"/data/package"}
GC_DAYS=${GC_DAYS:-"91"}
export REPO_URL="https://mirrors.tuna.tsinghua.edu.cn/git/git-repo/"

function repo_init() {
	# always start from init
	# as `REPO sync` does not GC well
	rm -fr $AOSP_DIR
	mkdir -p $AOSP_DIR
	cd $AOSP_DIR
	$REPO init -u http://aosp.tuna.tsinghua.edu.cn/platform/manifest
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

function gc() {
	rm -fr $AOSP_DIR
	cd $PKG_DIR
	find . -type f -mtime +${GC_DAYS} -delete
}

echo "Initializing AOSP working directory"
repo_init

echo "Packaging"
package

echo "Collecting old tarball"
gc
