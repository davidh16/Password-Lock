package service

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"time"
)

func (s Service) UploadIconToBucket(ctx *gin.Context, path string, file *multipart.FileHeader) error {

	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	// Initialize Firebase Storage client
	client, err := s.firebaseApp.Storage(context.Background())
	if err != nil {
		return err
	}

	// Create a storage reference
	storageRef, err := client.Bucket(s.Cfg.StorageBucket)
	if err != nil {
		return err
	}

	//ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	object := storageRef.Object(path)

	// Upload the file to Firebase Storage
	wc := object.NewWriter(context.Background())
	_, err = io.Copy(wc, uploadedFile)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetEntityIconSignedUrl(ctx *gin.Context, path string) (string, error) {

	// Initialize Firebase Storage client
	client, err := s.firebaseApp.Storage(context.Background())
	if err != nil {
		return "", err
	}

	// Create a storage reference
	storageRef, err := client.Bucket(s.Cfg.StorageBucket)
	if err != nil {
		return "", err
	}

	signedUrl, err := storageRef.SignedURL(path, &storage.SignedURLOptions{
		Expires: time.Now().AddDate(100, 0, 0),
		Method:  "GET",
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return signedUrl, nil
}
