@echo off
docker image prune -f 
docker volume prune -f
cd Backend
swag init
cd ..
docker-compose --env-file ".\Backend\.env" %*
exit /b