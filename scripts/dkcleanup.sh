echo "Cleaning up old docker containers..."
docker ps -a -q | xargs docker rm 
echo "Done"
 
