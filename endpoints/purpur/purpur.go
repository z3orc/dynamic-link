package purpur

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/z3orc/dynamic-link/util"
)

type Versions struct {
	Versions []string
}

type Version struct {
	Builds Builds
	Project string
	Version string
}

type Builds struct {
	All []string
	Latest string
}


func GetVersions() (Versions, error) {
	var versions Versions

	resp, err := util.GetJson("https://api.purpurmc.org/v2/purpur")
	if err != nil {
		return versions, err
	}

	err = json.Unmarshal(resp, &versions)
	if err != nil {
		return versions, err
	}

	return versions, nil
}

func GetVersion(id string) (Version, error){
	var version Version
	var url string;

	versions, err := GetVersions()
	if err != nil {
		return version, err
	}

	length := len(versions.Versions)

	for i := 0; i < int(length); i++ {
		currentId := versions.Versions[i]

		if currentId == id {
			url = fmt.Sprint("https://api.purpurmc.org/v2/purpur/",currentId)
			break
		}
	}

	if url == ""{
		err := errors.New("404")
		return version, err
	}

	resp, err := util.GetJson(url)
	if err != nil {
		return version, err
	}

	err = json.Unmarshal(resp, &version)
	if err != nil {
		return version, err
	}

	return version, nil
}

func GetLatestBuild(id string) (string, error){
	version, err := GetVersion(id)
	if err != nil {
		return "", err
	}

	return version.Builds.Latest, nil
}

func GetDownloadUrl(id string) (string, error){
	latestBuild, err := GetLatestBuild(id)
	if err != nil {
		return "", err
	}

	url := fmt.Sprint("https://api.purpurmc.org/v2/purpur/", id, "/", latestBuild, "/download" )

	return url, nil
}