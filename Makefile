all: sok

sok:
	go build -o sok main.go

clean:
	rm -f sok

.PHONY: all clean sok
