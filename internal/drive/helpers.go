package drive

import (
	"encoding/json"
	"log"

	"github.com/omaressameldin/lunch-roulette/internal/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	drive "google.golang.org/api/drive/v3"
)

func canUseDrive() bool {
	return env.DoesHaveCredentials()
}

func serviceAccount() *drive.Service {
	credentials := env.GetDriveCredentials()

	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal([]byte(credentials), &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}

	client := config.Client(oauth2.NoContext)
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalln(err)
	}

	return srv
}

func getFile(service *drive.Service, fileName string) *drive.File {
	list, err := service.Files.List().Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range list.Files {
		if f.Name == fileName {
			return f
		}
	}

	return nil
}

func deleteFile(service *drive.Service, f *drive.File) {
	if f == nil {
		return
	}

	err := service.Files.Delete(f.Id).Do()
	if err != nil {
		log.Fatal(err)
	}
}
