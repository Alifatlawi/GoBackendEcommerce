package controllers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func uploadToAzure(file *multipart.FileHeader, filename string) (string, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER")

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

// deleteFromAzure deletes a blob from Azure Storage
func deleteFromAzure(imageURL string) error {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Invalid credentials with error: %s", err.Error())
		return err
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, err := url.Parse(imageURL)
	if err != nil {
		return err
	}

	blobURL := azblob.NewBlobURL(*URL, p)

	_, err = blobURL.Delete(context.Background(), azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		log.Printf("Failed to delete blob: %v", err)
		return err
	}

	return nil
}
