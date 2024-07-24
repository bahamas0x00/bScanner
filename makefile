BINARY_NAME=bScanner
SRC_FILE=main.go

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ${SRC_FILE}
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ${SRC_FILE}
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows ${SRC_FILE}
	GOARCH=arm64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ${SRC_FILE}
	GOARCH=arm64 GOOS=linux go build -o ${BINARY_NAME}-arm64 ${SRC_FILE}
	

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -f ${BINARY_NAME}-darwin
	rm -f ${BINARY_NAME}-linux
	rm -f ${BINARY_NAME}-windows