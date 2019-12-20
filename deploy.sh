DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
echo Deploying from $DIR

echo Installing all necessary packages
# Update and install the sub modules
cd ../BotSecrets
git pull
go install
cd ../ScriptureBot
git pull
go install

echo Installing app entry package
# Update and install the entry module
cd $DIR
git pull
go install

echo Deploying app
gcloud app deploy --quiet