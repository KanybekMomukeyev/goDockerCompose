// Install Minio library.
//
// $ go get -u github.com/minio/minio-go
//
package main

import (
	"log"
	"github.com/minio/minio-go" // Import Minio library.
	"fmt"
)

func main() {
	// Use a secure connection.
	ssl := false

	// Initialize minio client object.
	minioClient, err := minio.New("138.68.84.55:9000",
		"AKIAIOSFODNN7EXAMPLE",
		"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", ssl)

	if err != nil {
		log.Fatalln(err)
	}

	// Creates bucket with name mybucket.
	err = minioClient.MakeBucket("mybucket2", "us-east-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully created mybucket.")

	// Upload an object 'myobject.txt' with contents from '/home/joe/myLocalFilename.txt'
	n, err := minioClient.FPutObject("mybucket2",
		"myobject.txt",
		"api/minio/myLocalFilename.txt",
		"application/text")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(n)

}