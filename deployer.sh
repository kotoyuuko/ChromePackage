#!/bin/bash

mkdir ./deploy
cd ./deploy

git init
git config --global push.default matching
git config --global user.email "${GitHubEMail}"
git config --global user.name "${GitHubUser}"
git remote add origin https://${GitHubKEY}@github.com/kotoyuuko/ChromePackage.git
git pull origin gh-pages

rm -rf ./*
mv ../chrome.json ../deploy/

git add --all .
git commit -m "Daily check by Travis CI"
git push --quiet --force origin HEAD:gh-pages

cd ..
rm -rf ./deploy
