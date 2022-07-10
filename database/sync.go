package database

import (
	"context"
	"fmt"
	"strings"

	"example.com/m/v2/util"
	"github.com/tidwall/gjson"
)

func SyncVanilla() {

	ctx := context.Background()
	client := Connect()

	jsonFromWeb, _ := util.GetJson("https://launchermeta.mojang.com/mc/game/version_manifest_v2.json")

	length := gjson.Get(jsonFromWeb, "versions.#").Int()

	for i := 0; i < int(length); i++ {
		path := fmt.Sprint("versions.", i, ".url")

		id := gjson.Get(jsonFromWeb, (fmt.Sprint("versions.", i, ".id"))).String()
		id = strings.ReplaceAll(id, " ", "-")

		build := gjson.Get(jsonFromWeb, path)

		buildJson, _ := util.GetJson(build.String())

		download_url := gjson.Get(buildJson, "downloads.server.url").String()
		old_download_url, _ := client.HGet(ctx, "vanilla", id).Result()

		if(download_url == old_download_url){
			break
		}

		if download_url != "" {
			fmt.Println(id, download_url)
			client.HSet(ctx, "vanilla", id, download_url)
		} else {
			fmt.Println("No download url")
		}

	}
}