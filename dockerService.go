package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"strings"
	"fmt"
)

func searchContainer(containerImage string) (string, string, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}

	// Check if any container has the target image name
	for _, container := range containers {
		// Split l'immagine del contenitore utilizzando "/"
		imageParts := strings.Split(container.Image, "/")
	
		// Controlla se ci sono almeno 2 parti nell'URL dell'immagine
		if len(imageParts) < 2 {
			continue
		}
	
		// Estrai il nome dell'immagine con etichetta da imageParts[1]
		imageNameWithTag := imageParts[1]
	
		// Split il nome dell'immagine con etichetta utilizzando ":"
		nameAndTag := strings.Split(imageNameWithTag, ":")
	
		// Controlla se c'è almeno una parte (il nome dell'immagine)
		if len(nameAndTag) < 1 {
			continue
		}
	
		// Estrai il nome dell'immagine
		imageName := nameAndTag[0]
	
		// Controlla se il nome dell'immagine corrisponde a containerImage
		if imageName == containerImage {
			tag := "" // Inizializza la variabile per l'etichetta
			if len(nameAndTag) > 1 {
				tag = nameAndTag[1] // Estrai l'etichetta se presente
			}
			return container.ID, tag, nil
		}
	}
	

	return "", "", nil
}

func checkEdgeModulePresence(imageName string) (string, error) {
	ID, version, err := searchContainer(imageName)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	if ID != "" {
		return "The current version of the edge module is: " + version + "\n", nil
	}

	return "The edge module is not running\n", nil
}


func updateEdgeModule(containerName, imageName, imageVersion string) (string, error) {
	
	fmt.Println("Updating to version " + imageVersion)

	// Crea un client Docker
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}

	// Cerca se il container è running
	ID, _, err := searchContainer(imageName)
	if err != nil {
		return "",err
	}

	// Ferma il container se sta runnando
	if ID != "" {
		fmt.Println("Stopping the running container with ID " + ID)
		// Ferma il container (se in esecuzione)
		err = cli.ContainerStop(context.Background(), ID, container.StopOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return "",err
		}

		// Rimuovi il vecchio container (se esiste)
		err = cli.ContainerRemove(context.Background(), ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
		}

	newImageName := imageName+":"+imageVersion
	fmt.Println("Pulling the image " + newImageName)	
	// Tira la nuova immagine dal registro
	image, err := pullImage(imageName, imageVersion)
	if err != nil {
		return "", err
	}
	fmt.Println("Successfully pulled image")	

	fmt.Println("Creating new container")
	// Crea e avvia un nuovo container con l'immagine aggiornata
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		nil,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	fmt.Printf("Container %s created with image %s\n", containerName, newImageName)
	return "Edge module updated\n", nil
}

func updateEdgeModuleLatest(containerName, imageName string) (string, error) {
	versions, err := listVersions(imageName)
	if err != nil {
		return "", err
	}

	imageVersion := versions[len(versions)-1]

	fmt.Println("Updating to latest version : " + imageVersion)

	// Crea un client Docker
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}

	// Cerca se il container è running
	ID, _, err := searchContainer(imageName)
	if err != nil {
		return "",err
	}

	// Ferma il container se sta runnando
	if ID != "" {
		fmt.Println("Stopping the running container with ID " + ID)
		// Ferma il container (se in esecuzione)
		err = cli.ContainerStop(context.Background(), ID, container.StopOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return "",err
		}

		// Rimuovi il vecchio container (se esiste)
		err = cli.ContainerRemove(context.Background(), ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
		}

	newImageName := imageName+":"+imageVersion
	fmt.Println("Pulling the image " + newImageName)	
	// Tira la nuova immagine dal registro
	image, err := pullImage(imageName, imageVersion)
	if err != nil {
		return "", err
	}
	fmt.Println("Successfully pulled image")	

	fmt.Println("Creating new container")
	// Crea e avvia un nuovo container con l'immagine aggiornata
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		nil,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	fmt.Printf("Container %s created with image %s\n", containerName, newImageName)
	return "Edge module updated\n", nil
}