package paper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/z3orc/dynamic-link/util"
)

type Versions struct {
	Versions []string
}

type Version struct {
	Builds []int
}

type Build struct {
	Downloads struct {
		Application struct {
			Name string
		}
	}
}


func GetVersions() (Versions, error) {
	resp, err := util.GetJson("https://api.papermc.io/v2/projects/paper")
	if err != nil {
		log.Fatal(err)
	}

	versions := Versions{}

	err = json.Unmarshal(resp, &versions)
	if err != nil {
		log.Fatal(err)
	}

	return versions, nil
}

func GetVersion(version string) (Version, error){
	builds := Version{}
	url := "https://api.papermc.io/v2/projects/paper/versions/" + version

	err := util.CheckUrl(url)
	if err != nil {
		return builds, errors.New("404")
	}

	resp, err := util.GetJson("https://api.papermc.io/v2/projects/paper/versions/" + version)
	if err != nil {
		return builds, err
	}

	err = json.Unmarshal(resp, &builds)
	if err != nil {
		return builds, err
	}

	return builds, nil
}

func GetLatestBuild(id string) (string, error){
	version, err := GetVersion(id)
	if err != nil {
		return "", err
	}
	builds := version.Builds

	latest := builds[len(builds) - 1]
	latestAsString := fmt.Sprintf("%v", latest)

	return latestAsString, nil
}

func GetJarName(id string) (string, error){
	latestBuild, err := GetLatestBuild(id)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%s",id,latestBuild)

	resp, err := util.GetJson(url)
	if err != nil {
		return "", err
	}

	build := Build{}

	err = json.Unmarshal(resp, &build)
	if err != nil {
		return "", err
	}

	return build.Downloads.Application.Name, nil
}

func GetDownloadUrl(id string) (string, error){
	latestBuild, err := GetLatestBuild(id)
	if err != nil {
		return "", err
	}
	jarName, err := GetJarName(id)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%s/downloads/%s", id, latestBuild, jarName)

	return url, nil
}