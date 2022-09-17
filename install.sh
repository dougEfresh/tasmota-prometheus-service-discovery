#!/usr/bin/env bash

set -euo pipefail

read -rp "GitHub Username: " user
read -rp "Projectname: " projectname

mkdir ../"$projectname"
rsync -a --progress --exclude .git --exclude .idea --exclude .vscode   ./ ../"$projectname"/
cd ../"$projectname"
mv README_TEMPLATE.md README.md
find . -type f -exec sed -i "s/tasmota-prometheus-service-discovery/$projectname/g" {} +
find . -type f -exec sed -i "s/dougEfresh/$user/g" {} +
git init
git add .
git commit -m "initial commit"
git remote add origin "git@github.com:$user/$projectname.git"

echo "template successfully installed."
