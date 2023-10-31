package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"sort"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"fmt"
	"strings"
	"encoding/base64"
	"encoding/json"
)

func listVersions(repositoryName string) ([]string, error) {
	accessKeyID := ""
	secretAccessKey := ""

	credentialsProvider := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	// Load AWS configuration from environment variables, AWS credentials file, and shared config file
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentialsProvider),
	)
	if err != nil {
		return nil, err
	}

	// Create a new ECR client
	client := ecr.NewFromConfig(cfg)

	// Specify the repository name and get the image list
	input := &ecr.DescribeImagesInput{
		RepositoryName: &repositoryName,
	}
	output, err := client.DescribeImages(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	// Extract image tags from the response
	var imageTags []string
	for _, image := range output.ImageDetails {
		imageTags = append(imageTags, image.ImageTags...)
	}

	sort.Strings(imageTags)

	return imageTags, nil
}

func pullImage(repositoryName, imageTag string) (string,error){
	accessKeyID := ""
	secretAccessKey := ""

	credentialsProvider := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	// Load AWS configuration from environment variables, AWS credentials file, and shared config file
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-south-1"),
		config.WithCredentialsProvider(credentialsProvider),
	)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}
	
	// Create a new ECR client
	ecrClient := ecr.NewFromConfig(cfg)
	
	// Get the authentication token to authenticate with the ECR registry
	auth, err := ecrClient.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}
	
	// Extract the ECR registry URL and authentication token
	registryURL := *auth.AuthorizationData[0].ProxyEndpoint
	encodedToken := *auth.AuthorizationData[0].AuthorizationToken

	decodedToken, err := base64.URLEncoding.DecodeString(encodedToken)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}

	authenticationToken := strings.Split(string(decodedToken), ":")[1]

	// Initialize Docker client with the extracted authentication token
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}
	// Authenticate with the ECR registry using the extracted token
	authConfig := types.AuthConfig{
		Username:      "AWS",
		Password:      authenticationToken,
		ServerAddress: registryURL,
	}
	authConfigBytes, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}


 	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	registryWithoutProtocol := strings.TrimPrefix(registryURL, "https://")

	// Pull the image from the ECR repository
	fmt.Println("Pulling the image from ECR : " + registryWithoutProtocol+"/"+repositoryName+":"+imageTag)
	pullOptions := types.ImagePullOptions{RegistryAuth: authConfigEncoded,}
	_, err = dockerClient.ImagePull(context.TODO(), registryWithoutProtocol+"/"+repositoryName+":"+imageTag, pullOptions)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}

	return registryWithoutProtocol+"/"+repositoryName+":"+imageTag, nil

}