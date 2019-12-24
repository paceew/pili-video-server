cd api
go install
cd ..
cd scheduler
go install
cd ..
cd streamserver
go install
cd ..
cd web
go install
cp -R ~/go/bin/api ~/go/src/github.com/pili-video-server/api_p
cp -R ~/go/bin/scheduler ~/go/src/github.com/pili-video-server/scheduler_p
cp -R ~/go/bin/streamserver ~/go/src/github.com/pili-video-server/streamserver_p
cp -R ~/go/bin/web ~/go/src/github.com/pili-video-server/web_p