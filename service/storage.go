package service

import (
	"bytes"
	"context"
	"firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"io"
	"log"
	"password-lock/config"
	"strings"
	"time"
)

func (s Service) UploadIconToBucket(ctx *gin.Context, entityUuid string) (string, error) {

	file, err := ctx.FormFile("icon") // "file" is the name of the file input field in your HTML form
	if err != nil {
		return "", err
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	cfg := config.GetConfig()

	opt := option.WithCredentialsFile("password-lock-486ee-firebase-adminsdk-xtd5c-cc43257771.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return "", err
	}

	// Initialize Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		return "", err
	}

	// Create a storage reference
	storageRef, err := client.Bucket(cfg.StorageBucket)
	if err != nil {
		return "", err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	path := strings.Join([]string{s.Me(ctx), entityUuid}, "/")

	object := storageRef.Object(path)

	// Upload the file to Firebase Storage
	wc := object.NewWriter(ctxWithTimeout)
	_, err = io.Copy(wc, uploadedFile)
	if err != nil {
		return "", err
	}
	err = wc.Close()
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s Service) DownloadEntityIcon(ctx *gin.Context, entityUuid string) ([]byte, error) {

	cfg := config.GetConfig()

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
	storageRef, err := client.Bucket(cfg.StorageBucket)
	if err != nil {
		return nil, err
	}

	path := strings.Join([]string{s.Me(ctx), entityUuid}, "/")

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
