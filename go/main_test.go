package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

var appsFile = appsDir + "/apps.json"
var appsDir = "./test"

func setup() {
	err := os.Mkdir(appsDir, 0755)
	if err != nil {
		log.Fatalf("Error on setup %s", err.Error())
	}
}

func teardown() {
	err := os.RemoveAll(appsDir)
	if err != nil {
		log.Fatalf("Error on teardown %s", err.Error())
	}
}

func Test_parseLabels(t *testing.T) {

	type args struct {
		containers []types.Container
	}

	tests := []struct {
		name    string
		args    args
		want    []App
		wantErr bool
	}{
		{
			name: "Assert labels are parsed correctly into App struct",
			args: args{
				containers: []types.Container{
					{
						Names: []string{"MyApp", "MyApp"},
						Labels: map[string]string{
							"sui.app.url":  "MyApp.url.xyz",
							"sui.app.icon": "Mine",
							"sui.app.name": "MyApp",
						},
					},
				},
			},
			want: []App{
				{
					Name: "MyApp",
					Url:  "MyApp.url.xyz",
					Icon: "Mine",
				},
			},
			wantErr: false,
		},
		{
			name: "Assert unwanted labels are ignored",
			args: args{
				containers: []types.Container{
					{
						Names: []string{"NoLabelApp", "NoLabelApp"},
						Labels: map[string]string{
							"some.random.lable": "some.random.value",
						},
					},
				},
			},
			want:    []App{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLabels(tt.args.containers)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeAppsFile(t *testing.T) {

	setup()
	defer teardown()

	type args struct {
		apps       []App
		outputFile string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []App
	}{
		{
			name: "Assert file is created if it does not exist",
			args: args{
				apps: []App{
					{
						Name: "WhatsUpp",
						Url:  "WhatsUpp.xyz",
						Icon: "Phone",
					},
				},
				outputFile: appsFile,
			},
			want: []App{{
				Name: "WhatsUpp",
				Url:  "WhatsUpp.xyz",
				Icon: "Phone",
			},
			},
		},
		{
			name: "Assert file is overwritten if it already exists",
			args: args{
				apps: []App{
					{
						Name: "NothingsUpp",
						Url:  "NothingsUpp.xyz",
						Icon: "Nothing",
					},
				},
				outputFile: appsFile,
			},
			want: []App{{
				Name: "NothingsUpp",
				Url:  "NothingsUpp.xyz",
				Icon: "Nothing",
			},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeAppsFile(tt.args.apps, tt.args.outputFile)
			if err != nil {
				t.Error("Error should be nil")
			}

			appsParsed := readJson(appsFile)
			if !reflect.DeepEqual(tt.want, appsParsed) {
				fmt.Print("Result: ")
				fmt.Println(appsParsed)
				fmt.Print("Expected: ")
				fmt.Println(tt.want)
				t.Error("Is not equal")
			}
		})
	}
}

func readJson(filePath string) []App {

	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	a := struct {
		Apps []App `json:"apps"`
	}{
		[]App{},
	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'apps' which we defined above
	json.Unmarshal(byteValue, &a)
	return a.Apps

}
