package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type App struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Icon string `json:"icon"`
}

func main() {

	outputFile := flag.String("apps-config", "./config/apps.json", "Location of apps.json file")
	flag.Parse()

	log.Println("Run started")

	containers, err := getContainers()
	if err != nil {
		log.Fatal(err.Error())
	}

	apps, err := parseLabels(containers)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = writeAppsFile(apps, *outputFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Run finished")
}

// parseLabels retrieves Docker containers with sui labels
func parseLabels(containers []types.Container) ([]App, error) {

	apps := []App{}
	for _, container := range containers {

		app := newApp(container)

		// If app is empty, we do not append
		if app == (App{}) {
			log.Println("Container has no sui labels")
			continue
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// writeAppsFile writes apps to a json file. outputFile will be created if it does not exist.
// outputFile will always be overwritten.
func writeAppsFile(apps []App, outputFile string) error {

	err := createFileIfNotExists(outputFile)
	if err != nil {
		return err
	}

	a := struct {
		Apps []App `json:"apps"`
	}{
		apps,
	}

	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(outputFile, data, 0755)
	if err != nil {
		return err
	}
	return nil

}

func newApp(c types.Container) App {
	log.Printf("Parsing labels from container: %+q", c.Names)

	app := App{}

	app.Name = parseName(c.Labels)
	app.Url = parseUrl(c.Labels)
	app.Icon = parseIcon(c.Labels)

	return app

}

func parseName(m map[string]string) string {
	if val, ok := m["sui.app.name"]; ok {
		log.Printf("Container label sui.app.name: %s\n", val)
		return val
	}
	return ""
}

func parseUrl(m map[string]string) string {
	if val, ok := m["sui.app.url"]; ok {
		log.Printf("Container label sui.app.url: %s\n", val)
		return val
	}
	return ""
}

func parseIcon(m map[string]string) string {
	if val, ok := m["sui.app.icon"]; ok {
		log.Printf("Container label sui.app.icon: %s\n", val)
		return val
	}
	return ""
}

func getContainers() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []types.Container{}, err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return []types.Container{}, err
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
