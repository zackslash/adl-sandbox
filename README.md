# seed-application
Seed Application

#For local development
source export.sh


git remote add seed git@github.com:cubex/seed-application.git

git pull seed master:master

git reset $(git commit-tree HEAD^{tree} -m "seed commit")
