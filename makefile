main:
	sh ./deployments/build.sh
deploy:
	go run main.go deploy --conf configs/debug.ini --dir  ~/Documents/github.com/book
	sh ./deployments/sync_image.sh
dev:
	go run main.go start --conf configs/debug.ini
	
