package database

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"example.com/m/v2/util"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
)

func HeartBeat() {
	func(){
		for range time.Tick(time.Hour * 12) {
			Sync()
		}
	}()
    fmt.Scanln()
}

func Sync(){

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		err := syncPurpur()
		if err != "" {
			fmt.Println("[PURPUR]:  ",err)
		}
		wg.Done()
	}()

	go func() {
		err := syncVanilla()
		if err != "" {
			fmt.Println("[VANILLA]: ",err)
		}
		wg.Done()
	}()

	go func() {
		err := syncPaper()
		if err != "" {
			fmt.Println("[PAPER]:   ",err)
		}
		wg.Done()
	}()

	wg.Wait()
}

func initSync(url string) (context.Context, *redis.Client, string,string) {
	ctx := context.Background()
	client := Connect()
	state := Check(client)

	if(!state){
		return ctx, client, "", "Could not connect to redis"
	}

	jsonFromWeb, err := util.GetJson(url)

	if err != nil {
		return ctx, client, jsonFromWeb, "Could not get json from web"
	}

	return ctx, client, jsonFromWeb, ""
}

func syncVanilla() string{
 	ctx, client, jsonFromWeb, err := initSync("https://launchermeta.mojang.com/mc/game/version_manifest_v2.json")

	if err != ""{
		return err
	}	

	length := gjson.Get(jsonFromWeb, "versions.#").Int()

	for i := 0; i < int(length); i++ {
		path := fmt.Sprint("versions.", i, ".url")

		id := gjson.Get(jsonFromWeb, (fmt.Sprint("versions.", i, ".id"))).String()
		id = strings.ReplaceAll(id, " ", "-")
		isSnapshot := strings.Contains(id, "w")
		isRC := strings.Contains(id, "rc")
		isPre := strings.Contains(id, "pre")

		if isSnapshot || isPre || isRC {
			continue
		}

		build := gjson.Get(jsonFromWeb, path)

		buildJson, _ := util.GetJson(build.String())

		download_url := gjson.Get(buildJson, "downloads.server.url").String()
		old_download_url, _ := client.HGet(ctx, "vanilla", id).Result()

		if(download_url == old_download_url){
			break
		}

		if download_url != "" {
			fmt.Println("[VANILLA]: " , id, download_url)
			client.HSet(ctx, "vanilla", id, download_url)
		} else {
			fmt.Println("No download url")
		}

	}
	client.Close()
	return ""
}

func syncPaper() string{
	ctx, client, jsonFromWeb, err := initSync("https://api.papermc.io/v2/projects/paper")

	if err != ""{
		return err
	}	

	length := gjson.Get(jsonFromWeb, "versions.#").Int()

	for i := 0; i < int(length); i++ {
		id := gjson.Get(jsonFromWeb, fmt.Sprint("versions.", i)).String()

		versionJsonPath := fmt.Sprint("https://api.papermc.io/v2/projects/paper/versions/", id)
		versionJson, _ := util.GetJson(versionJsonPath)
		buildLength := gjson.Get(versionJson, "builds.#").Int()
		build := gjson.Get(versionJson, (fmt.Sprint("builds.", buildLength - 1))).Int()

		buildJsonPath := fmt.Sprint(versionJsonPath, "/builds/", build)
		buildJson, _ := util.GetJson(buildJsonPath)

		fileName := gjson.Get(buildJson, "downloads.application.name")

		downloadUrl := fmt.Sprint(buildJsonPath, "/downloads/", fileName)

		oldDownloadUrl, _ := client.HGet(ctx, "paper", id).Result()

		if(downloadUrl == oldDownloadUrl){
			break
		}

		if downloadUrl != "" {
			fmt.Println("[PAPER]:   ", id, downloadUrl)
			client.HSet(ctx, "paper", id, downloadUrl)
		} else {
			fmt.Println("No download url")
		}
	}
	client.Close()
 	return ""
}

func syncPurpur() string{
	ctx, client, jsonFromWeb, err := initSync("https://api.purpurmc.org/v2/purpur")

	if err != ""{
		return err
	}	

	length := gjson.Get(jsonFromWeb, "versions.#").Int()

	for i := 0; i < int(length); i++ {
		id := gjson.Get(jsonFromWeb, fmt.Sprint("versions.", i)).String()

		versionJsonPath := fmt.Sprint("https://api.purpurmc.org/v2/purpur/", id)
		versionJson, _ := util.GetJson(versionJsonPath)
		buildLength := gjson.Get(versionJson, "builds.all.#").Int()
		build := gjson.Get(versionJson, (fmt.Sprint("builds.all.", buildLength - 1))).Int()

		buildPath := fmt.Sprint(versionJsonPath, "/", build)

		downloadUrl := fmt.Sprint(buildPath, "/download")

		oldDownloadUrl, _ := client.HGet(ctx, "purpur", id).Result()

		if(downloadUrl == oldDownloadUrl){
			break
		}

		if downloadUrl != "" {
			fmt.Println("[PURPUR]:  ", id, downloadUrl)
			client.HSet(ctx, "purpur", id, downloadUrl)
		} else {
			fmt.Println("No download url")
		}
	}
	client.Close()
	return ""
}