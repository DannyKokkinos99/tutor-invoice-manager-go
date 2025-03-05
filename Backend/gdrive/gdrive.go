package gdrive

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GDrive struct {
	service *drive.Service
}

func NewGDrive(serviceAccountFile string) (*GDrive, error) {
	b, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to read service account key file: %v", err)
	}

	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to create Drive client: %v", err)
	}

	return &GDrive{service: srv}, nil
}

func (g *GDrive) CreateFolder(name, parentID string) (*drive.File, error) {
	// Check if a folder with the same name already exists in the parent directory
	query := fmt.Sprintf("'%s' in parents and mimeType='application/vnd.google-apps.folder' and name='%s'", parentID, name)

	r, err := g.service.Files.List().
		Q(query).
		Fields("files(id, name)"). // Only retrieve the id and name of existing folders
		PageSize(1).               // We only need to check if at least one exists
		Do()

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to check for existing folder: %v", err)
	}

	if len(r.Files) > 0 {
		return r.Files[0], nil
	}

	// Create the new folder since no folder with the same name exists
	folder := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}

	file, err := g.service.Files.Create(folder).Do()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to create folder: %v", err)
	}

	return file, nil
}

func (g *GDrive) CreateFile(name, parentID, mimeType string, content []byte) (*drive.File, error) {
	// Create a new file instance with the specified MIME type
	file := &drive.File{
		Name:     name,
		Parents:  []string{parentID},
		MimeType: mimeType, // Specify the MIME type here
	}

	// Use an io.Reader to provide the content for the upload
	fileRes, err := g.service.Files.Create(file).Media(bytes.NewReader(content)).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create file in Drive: %v", err)
	}

	return fileRes, nil
}

func (g *GDrive) PrintFolderContents(folderID string, pageSize int64, printIDs bool) error {
	var pageToken string
	for {
		r, err := g.service.Files.List().
			Q(fmt.Sprintf("'%s' in parents", folderID)).
			PageSize(pageSize).
			PageToken(pageToken).
			Fields("nextPageToken, files(id, name, mimeType)").
			Do()

		if err != nil {
			return fmt.Errorf("unable to retrieve folder contents: %v", err)
		}

		if len(r.Files) == 0 {
			fmt.Println("No files found.")
			return nil
		}

		for _, file := range r.Files {

			if file.MimeType == "application/vnd.google-apps.folder" {
				if printIDs {
					fmt.Printf("  --> Entering folder: %s (%s)\n", file.Name, file.Id)
				} else {
					fmt.Printf("  --> Entering folder: %s\n", file.Name)
				}

				if err := g.PrintFolderContents(file.Id, pageSize, printIDs); err != nil {
					return err
				}
			} else {
				if printIDs {
					fmt.Printf("- %s (%s)\n", file.Name, file.Id)
				} else {
					fmt.Printf("- %s\n", file.Name)
				}

			}
		}

		// Check if there's a next page
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return nil
}

func (g *GDrive) UploadFile(localFilePath, parentID, desiredFileName string) (*drive.File, error) {
	// Open the local file
	file, err := os.Open(localFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open local file: %v", err)
	}
	defer file.Close()

	// Determine the MIME type based on the file extension
	ext := filepath.Ext(localFilePath)
	mimeType := mime.TypeByExtension(ext)

	// If the MIME type could not be determined, set a default type
	if mimeType == "" {
		mimeType = "application/octet-stream" // Default to binary stream
	}

	// Create a new file instance with the specified properties
	driveFile := &drive.File{
		Name:     desiredFileName,
		Parents:  []string{parentID},
		MimeType: mimeType,
	}

	// Upload the file to Google Drive
	fileRes, err := g.service.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to upload file to Drive: %v", err)
	}

	return fileRes, nil
}
