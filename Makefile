all: mac

macx86:
	GOOS=darwin GOARCH=amd64 go build -o sok main.go

mac:
	GOOS=darwin GOARCH=arm64 go build -o sok main.go

lin:
	GOOS=linux GOARCH=amd64 go build -o sok main.go

win:
	GOOS=windows GOARCH=amd64 go build -o sok.exe main.go

release:
	GOOS=darwin GOARCH=amd64 go build -o build/macx86/sok main.go
	GOOS=darwin GOARCH=amd64 go build -o build/macx86/sokru main.go
	GOOS=darwin GOARCH=arm64 go build -o build/macm1/sok main.go
	GOOS=darwin GOARCH=arm64 go build -o build/macm1/sokru main.go
	GOOS=linux GOARCH=amd64 go build -o build/linux/sok main.go
	GOOS=linux GOARCH=amd64 go build -o build/linux/sokru main.go
	GOOS=windows GOARCH=amd64 go build -o build/win/sok.exe main.go
	GOOS=windows GOARCH=amd64 go build -o build/win/sokru.exe main.go
	zip -r build/sok-macx86.zip build/macx86
	zip -r build/sok-macm1.zip build/macm1
	zip -r build/sok-linux.zip build/linux
	zip -r build/sok-win.zip build/win

clean:
	rm -f sok{.exe} build/*

.PHONY: all clean mac macx86 lin win
