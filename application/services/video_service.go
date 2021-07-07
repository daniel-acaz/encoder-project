package services

import (
	"cloud.google.com/go/storage"
	"context"
	"encoder-project/application/repositories"
	"encoder-project/domain"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type VideoService struct {
	Video 				*domain.Video
	VideoRepository 	repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(v.Video.FilePath)

	reader, err := object.NewReader(ctx)

	if err != nil {
		return  err
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	file, err := os.Create(os.Getenv("localStorageEnv") + "/" + v.Video.ID + ".mp4")
	if err != nil{
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	log.Printf("vídeo %v has been stored", v.Video.ID)

	return nil

}

func (v *VideoService) Fragment() error {

	err := os.Mkdir(os.Getenv("localStorageEnv") + "/" + v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStorageEnv") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStorageEnv") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("localStorageEnv") + "/" + v.Video.ID + ".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStorageEnv") + "/" + v.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {

	err := os.Remove(os.Getenv("localStorageEnv") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", v.Video.ID, ".mp4")
	}

	err = os.Remove(os.Getenv("localStorageEnv") + "/" + v.Video.ID + ".frag")
	if err != nil {
		log.Println("error removing mp4 ", v.Video.ID, ".frag")
	}

	err = os.RemoveAll(os.Getenv("localStorageEnv") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", v.Video.ID, ".mp4")
	}

	log.Println("files have been removed: ", v.Video.ID)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}
