main:
	go build -ldflags "-s -w" -a -o douyacun main.go
build:
	sh ./deployments/build.sh
deploy:
	go run main.go deploy --conf configs/debug.ini --dir  ~/Documents/github.com/book
	sh ./deployments/sync_image.sh
