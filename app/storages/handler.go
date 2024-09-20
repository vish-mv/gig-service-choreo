package storages

import (
	"GIG/app/storages/interfaces"
	"GIG/app/storages/minio"
	"github.com/lsflk/gig-sdk/libraries"
	"github.com/revel/revel"
	"log"
	"os"
)

var fileStorageHandler interfaces.StorageHandlerInterface

type FileStorageHandler struct {
	interfaces.StorageHandlerInterface
}

func LoadStorageHandler() {
	cacheDirectory, _ := revel.Config.String("file.cache")

	if cacheDirectory == "" { // default value
		cacheDirectory = "cache/"
	}

	if err := libraries.EnsureDirectory(cacheDirectory); err != nil {
		log.Fatal(err)
	}

	fileStorageHandler = minio.NewHandler(cacheDirectory) //change storage handler here
}

func (f FileStorageHandler) GetFile(directoryName string, filename string) (*os.File, error) {
	var localFile *os.File
	tempDir := fileStorageHandler.GetCacheDirectory() + directoryName + "/"
	sourcePath := tempDir + filename

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) { // if file is not cached
		localFile, err = fileStorageHandler.GetFile(directoryName, filename)
		if err != nil {
			log.Println(err)
			return localFile, err
		}

	} else { // if file is cached
		localFile, err = os.Open(sourcePath)
		if err != nil {
			log.Println(err)
			return localFile, err
		}
	}
	return localFile, nil
}

func (f FileStorageHandler) UploadFile(directoryName string, filePath string) error {
	return fileStorageHandler.UploadFile(directoryName, filePath)
}

func (f FileStorageHandler) GetCacheDirectory() string {
	return fileStorageHandler.GetCacheDirectory()
}
