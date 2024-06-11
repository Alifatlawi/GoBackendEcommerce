package controllers

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"mime/multipart"
	"net/url"
	"os"
)

func uploadToAzure(file *multipart.FileHeader, filename string) (string, error) {

	var accountName = os.Getenv("AZURE_STORAGE_ACCOUNT")
	var accountKey = os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	var containerName = os.Getenv("AZURE_STORAGE_CONTAINER")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Invalid credentials with error: %s", err.Error())
		return "", err
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)
	blobURL := containerURL.NewBlockBlobURL(filename)

	fileContent, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return "", err
	}
	defer fileContent.Close()

	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), fileContent, blobURL,
		azblob.UploadStreamToBlockBlobOptions{
			BufferSize: 4 * 1024 * 1024,
			MaxBuffers: 20,
		})

	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return "", err
	}

	return blobURL.String(), nil
}
