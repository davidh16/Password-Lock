package service

import (
	"bytes"
	"context"
	"firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"password-lock/config"
	"strings"
)

func (s Service) UploadIconToBucket(ctx *gin.Context, path string, file *multipart.FileHeader) error {

	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	cfg := config.GetConfig()

	opt := option.WithCredentialsFile("password-lock-486ee-firebase-adminsdk-xtd5c-cc43257771.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	// Initialize Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		return err
	}

	// Create a storage reference
	storageRef, err := client.Bucket(cfg.StorageBucket)
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

func (s Service) DownloadEntityIcon(ctx *gin.Context, entityUuid string) ([]byte, error) {

	Cfg := config.GetConfig()

	opt := option.WithCredentialsFile("password-lock-486ee-firebase-adminsdk-xtd5c-cc43257771.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	// Initialize Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	// Create a storage reference
	storageRef, err := client.Bucket(Cfg.StorageBucket)
	if err != nil {
		return nil, err
	}

	me := ctx.Value("me").(string)

	path := strings.Join([]string{me, entityUuid}, "/")

	object := storageRef.Object(path)
	if err != nil {
		return nil, err
	}

	reader, err := object.NewReader(context.Background())
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	// Copy the remote file's data to the buffer.
	_, err = io.Copy(&buffer, reader)
	if err != nil {
		log.Fatalf("Error reading file content: %v", err)
	}

	return buffer.Bytes(), nil
}
