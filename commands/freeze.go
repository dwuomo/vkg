package commands

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/pepegar/vkg/config"
	"github.com/pepegar/vkg/config/vkgrc"
	"github.com/pepegar/vkg/utils"
)

func action() {
	config := config.GetVkgGonfig()
	files, _ := ioutil.ReadDir(config.PluginsPath)

	var wg sync.WaitGroup

	//errors := make(chan error)
	pluginsChan := make(chan vkgrc.VkgrcPlugin, len(files))

	for _, file := range files {
		fullyQualifiedPath := config.PluginsPath + file.Name()

		wg.Add(1)

		go func(path string) {
			defer wg.Done()

			branch, _ := utils.Git.GetBranchName(fullyQualifiedPath)
			repo, _ := utils.Git.GetRepository(fullyQualifiedPath)
			name, _ := utils.Git.GetRepoName(fullyQualifiedPath)

			plugin := vkgrc.VkgrcPlugin{
				Branch:     branch,
				Repository: repo,
				Name:       name,
			}

			pluginsChan <- plugin
		}(fullyQualifiedPath)
	}

	// wait for all goroutines to finish
	wg.Wait()
	close(pluginsChan)

	var plugins []vkgrc.VkgrcPlugin

	for plugin := range pluginsChan {
		plugins = append(plugins, plugin)
	}

	vkgrcFile := vkgrc.VkgrcJSON{
		Plugins: plugins,
	}

	a, _ := json.MarshalIndent(vkgrcFile, "", "  ")

	println(string(a))
}

var FreezeCommand = Command{
	Name:        "freeze",
	Usage:       "freeze",
	Description: "output installed plugins in .vkgrc format",
	Action:      action,
}
