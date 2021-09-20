#!/bin/sh

BINARY_NAME_UNIX=run.sh
BINARY_NAME_WINDOWS=run.bat
BIN_FOLDER=bin


# Node-Build:
	mkdir -p $BIN_FOLDER
	cp -r src/* "$BIN_FOLDER"
	npm install --silent --no-progress --prefix "$BIN_FOLDER"

	#Unix
	{
	echo "#!/bin/sh"
	echo "node \$(dirname \"\$0\")/index.js"

	echo 'if [ "$GITLAB_REPOSITORY" != "" ] ; then'
	echo 'cd "$CURRENT_PWD"'
	echo "# Init Repo"
	echo 'git init'
	echo 'git remote add origin https://${RIT_GITLAB_USERNAME}:${RIT_GITLAB_TOKEN}@${GITLAB_REPOSITORY}-app'
	echo ''
	echo '# Commit Repo'
	echo 'git config user.email "$RIT_GITLAB_EMAIL" > /dev/null'
  echo 'git config user.name "$RIT_GITLAB_USERNAME" > /dev/null'
  echo 'git add . && git commit -m "Step Pipes" -s'
  echo '# Update Repo'
  echo 'git pull origin master --allow-unrelated-histories --no-edit'
  echo 'git push -u origin master'
  echo ''

  echo "fi"

	} >>  $BIN_FOLDER/$BINARY_NAME_UNIX
	chmod +x "$BIN_FOLDER/$BINARY_NAME_UNIX"

	#Windows
	echo "node index.js" > $BIN_FOLDER/$BINARY_NAME_WINDOWS

#Docker Files:
	cp Dockerfile set_umask.sh $BIN_FOLDER
