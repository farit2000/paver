#!/bin/sh
YUM_CMD=$(which yum)
DNF_CMD=$(which dnf)
APT_GET_CMD=$(which apt-get)
PACKAGE=git
echo "GIT SETUP"
if [[ ! -z $YUM_CMD ]]; then
	sudo yum update
	sudo yum install $PACKAGE
    yum install $YUM_PACKAGE_NAME
elif [[ ! -z $DNF_CMD ]]; then
	sudo dnf update
	sudo dnf install $PACKAGE
elif [[ ! -z $APT_GET_CMD ]]; then
	sudo apt-get update
	sudo apt-get install $PACKAGE
else
    echo "error can't install package $PACKAGE"
    exit 1;
fi
git --version
echo "SUCCESS GIT SETUP"
