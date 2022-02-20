package main

import (
	"context"
  "encoding/json"
  "io/ioutil"
  "os"
	"fmt"
  "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
	"flag"
	"strings"
)

type App struct {
  Name  string `json:"name"`
  Url   string `json:"url"`
  Icon  string `json:"icon"`
}

type Apps struct {
  Apps []App `json:"apps"`
}

func main(){
	filename := flag.String("apps-config", "/config/apps.json", "Location of apps.json file")
	interval := flag.Int("check-interval", 30, "Interval in seconds for checking container labels")
	flag.Parse()

  checkFileExists(*filename)

	for {
		fmt.Println()
		fmt.Println("Starting run")
		fmt.Println("---------------------")
		fmt.Println()
    updateJson(*filename)
		fmt.Println()
		fmt.Println("---------------------")
		fmt.Println("Stopping run")
		fmt.Println()
		time.Sleep(time.Duration(*interval) * time.Second)
  }

}

func checkFileExists(filename string) error {
    _, err := os.Stat(filename)
        if os.IsNotExist(err) {
            _, err := os.Create(filename)
                if err != nil {
                    return err
                }
        }
        return nil
}


func updateJson(filename string){
	  containers := getContainers()

	  apps_empty := []App{}
		apps := Apps{apps_empty}

	  for _, container := range containers {
			fmt.Println("---")
			fmt.Println("Found container with name:", strings.Trim(fmt.Sprint(container.Names), "/[]"))
			app := App{}

	    for key, value := range container.Labels{
	      if key == "sui.app.name" {
					fmt.Printf("Container label sui.app.name: %s\n", value)
	        app.Name = value
	      }
	      if key == "sui.app.url" {
	        app.Url = value
					fmt.Printf("Container label sui.app.url: %s\n", value)
	      }
	      if key == "sui.app.icon" {
	        app.Icon = value
					fmt.Printf("Container label sui.app.icon: %s\n", value)
	      }
	    }
			if (App{}) != app  {
	    	  apps.AddItem(app)
			} else {
				fmt.Println("Container has no sui labels")
			}
			fmt.Println("---")
	  }

	  writeJson(filename, apps)

}

func getContainers() []types.Container {
  ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

  containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
  if err != nil {
    panic(err)
  }

  return containers

}

func (apps *Apps) AddItem(app App) []App {
	apps.Apps = append(apps.Apps, app)
	return apps.Apps
}

func writeJson(filename string, apps Apps){
  dat, err := json.MarshalIndent(apps, "", "    ")
  if err != nil {
    panic(err)
  }

  err = ioutil.WriteFile(filename, dat, 0644)

}
