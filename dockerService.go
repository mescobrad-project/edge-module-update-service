package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"strings"
	"fmt"
	"os"
	"os/exec"
	"io"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"crypto/tls"
	"gopkg.in/yaml.v3"
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
		return "The current version of the edge module is: " + version , nil
	}

	return "The edge module is not running", nil
}

func backupDirectory(containerImage, srcDirectory, dstDirectory string) error {

    // Create the Docker client.
    cli, err := client.NewEnvClient()
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

	containerName, _, _ := searchContainer(containerImage)

    // Get the plugin directory contents.
    reader, _, err := cli.CopyFromContainer(context.Background(), containerName,  srcDirectory)
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

    // Write the plugin directory contents to the destination folder.
    writer, err := os.Create(dstDirectory)
    if err != nil {
		fmt.Println(err.Error())
        return err
    }
    defer writer.Close()

    _, err = io.Copy(writer, reader)
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

    return nil
}

func restoreDirectory(containerImage, srcDirectory, dstDirectory string) error {

    // Create the Docker client.
    cli, err := client.NewEnvClient()
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

    // Get the plugin directory contents.
    reader, err := os.Open(srcDirectory)
    if err != nil {
		fmt.Println(err.Error())
        return err
    }
    defer reader.Close()

	containerName, _, _ := searchContainer(containerImage)

	// Create an exec instance.
    execConfig := types.ExecConfig{
        Cmd: []string{"touch", dstDirectory},
    }
    execID, err := cli.ContainerExecCreate(context.Background(), containerName, execConfig)
    if err != nil {
        return err
    }

    // Start the exec instance.
    err = cli.ContainerExecStart(context.Background(), execID.ID , types.ExecStartCheck{
        Tty: true,
    })
    if err != nil {
        return err
    }

    // Create the plugin directory in the container.
    err = cli.CopyToContainer(context.Background(), containerName, dstDirectory, reader, types.CopyToContainerOptions{})
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

    return nil
}

func findChildNode(service, field string, node *yaml.Node) *yaml.Node {
	compose := node.Content[0]
	var services *yaml.Node
	found := 0
	for _, out := range compose.Content {
		if found == 1 {
			services=out
			break
		}		
		if out.Value == "services" {
			found=1
		}
	}
	found=0
	for _, s := range services.Content {
		if found == 1 {
			for _, v := range s.Content {
				if found == 2 {
					fmt.Println(v.Value)
					return v
				}
				if v.Value == field {
					found=2
					fmt.Printf("Found child with value: ")
				}
			}
		}
		if s.Value == service {
			found=1
		}
	}
	return nil
}

func updateCompose(serviceName, fieldName, fieldValue string) error{
	
	fmt.Println("Updating the compose...")

	// Read the YAML file
	data, err := ioutil.ReadFile("docker/docker-compose.yaml")
	if err != nil {
		fmt.Println(err.Error())
        return err
    }
  
	// Unmarshal the YAML file
	/*var config interface{}
	err = yaml.Unmarshal(data, &config)*/
	var dockerCompose yaml.Node
	yaml.Unmarshal(data, &dockerCompose)
	if err != nil {
		fmt.Println(err.Error())
        return err
    }

	imageNode := findChildNode(serviceName, fieldName, &dockerCompose)
	if imageNode != nil {
		imageNode.SetString(fieldValue)
	}
	// Create a modified yaml file
	/*f, err := os.Create("modified.yaml")
	if err != nil {
		log.Fatalf("Problem creating file: %v", err)
	}
	defer f.Close()*/
	//f, err := os.Open("docker/docker-compose.yaml")
	//yaml.NewEncoder(f).Encode(dockerCompose.Content[0])
  
	// Update the version
	/*configMap := config.(map[string]interface{})
	services := configMap["services"].(map[string]interface{})
	selectedService := services[serviceName].(map[string]interface{})
	selectedService[fieldName] = fieldValue
  */
	// Marshal the modified YAML file
	/*output, err = yaml.Marshal(dockerCompose)
	if err != nil {
		fmt.Println(err.Error())
        return err
    }*/

	// Type-assert the config variable to a []byte slice
	/*data, ok := config.([]byte)

	if !ok {
		panic("Failed to type-assert config to []byte")
	}*/

	// Write the modified YAML file
	/*err = ioutil.WriteFile("docker/docker-compose.yaml", data, 0644)
	if err != nil {
		fmt.Println(err.Error())
        return err
    }*/
	// Create a modified yaml file
	f, err := os.Create("docker/docker-compose.yaml")
	if err != nil {
		fmt.Println("Problem creating file: %v", err)
	}
	defer f.Close()
	yaml.NewEncoder(f).Encode(dockerCompose.Content[0])

	fmt.Println("Running new compose")
	// run the new compose
	cmd := exec.Command("docker", "compose", "-f", "docker/docker-compose.yaml", "up", "-d")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Compose updated")
	return(nil)
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

	fmt.Println("Backup the installed plugins . . .")
	
	backupDirectory(containerName, "/var/lib/docker/volumes/plugins", "backup.tar")

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

	updateCompose("mescobrad_edge", "image", image)

	fmt.Printf("Container %s created with image %s\n", containerName, newImageName)

	fmt.Printf("Restoring plugin . . .")

	restoreDirectory(containerName, "backup.tar", "/var/lib/docker/volumes/plugins")

	return "Edge module updated\n", nil
}

func updateEdgeModuleLatest(containerName, imageName string) (string, error) {

	fmt.Println("Backup the installed plugins . . .")
	
	backupDirectory(containerName, "/usr/src/app/mescobrad_edge/plugins", "backup.tar")

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

	err = updateCompose("mescobrad_edge", "image", image)
	if err != nil {
		return "", err
	}

	fmt.Printf("Container %s created with image %s\n", containerName, newImageName)

	fmt.Printf("Restoring plugin . . .")

	restoreDirectory(containerName, "backup.tar", "/var/lib/docker/volumes/plugins")

	fmt.Printf("Plugin restored")

	return "Edge module updated\n", nil
}

func installPlugin(containerName, pluginName string) error{
	// Create the HTTP client.
    transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

    // Create the request.
    //req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://mescobrad-edge:8080/api/v1/plugins", nil)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://localhost:8080/api/v1/plugins", nil)
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

    // Set the request headers.
    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", "abc")

    // Create the JSON body.
    body := map[string]string{
        "url": fmt.Sprintf("https://github.com/mescobrad-project/%s.git", pluginName),
    }

    // Marshal the JSON body to the request body.
    jsonBody, err := json.Marshal(body)
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

    req.Body = ioutil.NopCloser(bytes.NewReader(jsonBody))

    // Send the request.
    resp, err := client.Do(req)
    if err != nil {
		fmt.Println(err.Error())
        return err
    }

    // Close the response body.
    defer resp.Body.Close()

    // Check the response status code.
    if resp.StatusCode != http.StatusOK {
		fmt.Println(err.Error())
        return err
    }

	return nil
}

func updatePluginImpl(containerName, pluginName string) error {
    // Backup the plugin configuration file.
	fmt.Println("Making the backup of the current configuration...")
    err := backupDirectory(containerName, fmt.Sprintf("/usr/src/app/mescobrad_edge/plugins/%s/plugin.config", pluginName), "plugin.config")
    if err != nil {
        return err
    }
	fmt.Println("Backup completed successfully")
	
	//Backup mideface software
	if pluginName == "mri_anonymisation_plugin" {
		fmt.Println("Making the backup of the mideface software...")
		err := backupDirectory(containerName, fmt.Sprintf("/usr/src/app/mescobrad_edge/plugins/%s/mideface", pluginName), "mideface")
		if err != nil {
			return err
		}
		fmt.Println("Backup of mideface completed successfully")
	}

	// Install the plugin from GitHub.
	fmt.Println("Installing the plugin...")
    err = installPlugin(containerName, pluginName)
    if err != nil {
        return err
    }
	fmt.Println("Plugin installed successfully")

    // Restore the plugin configuration file.
	fmt.Println("Restoring the backup of the previous configuration...")
    err = restoreDirectory(containerName, "plugin.config", fmt.Sprintf("/usr/src/app/mescobrad_edge/plugins/%s", pluginName))
    if err != nil {
        return err
    }
	fmt.Println("Configurations restored successfully")

	//Restore mideface software
	if pluginName == "mri_anonymisation_plugin" {
		fmt.Println("Restoring the backup of the mideface software...")
		err = restoreDirectory(containerName, "mideface", fmt.Sprintf("/usr/src/app/mescobrad_edge/plugins/%s", pluginName))
		if err != nil {
			return err
		}
		fmt.Println("Backup of mideface restored successfully")
	}

    return nil
}