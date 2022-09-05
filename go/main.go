package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type App struct {
	name string `json:"name"`
	url  string `json:"url"`
	icon string `json:"icon"`
}

func main() {

	outputFile := flag.String("apps-config", "/config/apps.json", "Location of apps.json file")
	flag.Parse()

	log.Println("Starting run")

	err := writeAppsFile(*outputFile)

	if err != nil{
		log.Fatal(err.Error())
	}

	log.Println("Stopping run")
}

func writeAppsFile(outputFile string) error {

	err := createFileIfNotExists(outputFile)
	if err != nil {
		return err
	}

	containers, err := getContainers()
	if err != nil{
		return err
	}
	
	apps := []App{}
	for _, container := range containers {
		app := newApp{container}

		// If app is empty, we do not append
		if app == App {
			log.Println("Container has no sui labels")
			continue
		}
		append(apps, app)
	}

	err := apps.toJson(outputFile)

	if err != nil{
		return err
	}

	return nil
}

func newApp(c container) app {
	log.Println("Parsing labels from container:", strings.Trim(fmt.Sprint(c.Names), "/[]"))
	app := App{}

	app.name = parseName(container.Labels)
	app.url = parseUrl(container.Labels)
	app.icon = parseIcon(container.Labels)

	return app

}

func parseName(labels container.Labels) string {
	for key, value := range labels {
		if key == "sui.app.name" {
			log.Printf("Container label sui.app.name: %s\n", value)
			return value
		}
	}
	return ""
}

func parseUrl(labels container.Labels) string {
	for key, value := range labels {
		if key == "sui.app.url" {
			log.Printf("Container label sui.app.url: %s\n", value)
			return value
		}
	}
	return ""
}

func parseIcon(labels container.Labels) string {
	for key, value := range labels {
		if key == "sui.app.icon" {
			log.Printf("Container label sui.app.icon: %s\n", value)
			return value
		}
	}
	return ""
}

func getContainers() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []types.Container, err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return []types.Container, err
	}

	return containers, nil

}

func createFileIfNotExists(file string) error {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		_, err := os.Create(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (apps Apps) toJson(file string) error {
	dat, err := json.MarshalIndent(apps, "", "    ")
	if err != nil {
		retur err
	}

	err = ioutil.WriteFile(file, dat, 0644)
	if err != nil{
		return err
	}
}
