package vanilla

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/z3orc/dynamic-link/util"
)

type Versions struct {
	Latest   map[string]string
	Versions []VersionInfo
}

type VersionInfo struct {
	Id string
	Type string
	Url string
}

type Version struct {
	Downloads VersionDownloads
}

type VersionDownloads struct {
	Server struct {
		Sha1 string
		Size int32
		Url string
	}
}

func GetVersions() (Versions, error){
	versions := Versions{}

	resp, err := util.GetJson("https://piston-meta.mojang.com/mc/game/version_manifest_v2.json")
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
		currentId := versions.Versions[i].Id

		if currentId == id {
			url = versions.Versions[i].Url
			break
		}
	}

	if url == ""{
		err := errors.New("404")
		return version, err
	}

	resp, err := util.GetJson(url)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resp, &version)
	if err != nil {
		log.Fatal(err)
	}

	return version, nil
}

func GetDownloadUrl(id string) (string, error){
	version, err := GetVersion(id)
	if err != nil {
		return "", err
	}

	url := version.Downloads.Server.Url

	return url, nil
}