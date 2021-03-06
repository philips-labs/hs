/*
Copyright © 2020 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/philips-software/go-hsdp-api/has"

	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

// hasResourcesCreateCmd represents the create command
var hasResourcesCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c", "n", "new"},
	Short:   "Create a HAS resource",
	Long:    `Creates a HAS resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := cmd.Flags().GetInt("count")
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		images, _, err := client.Images.GetImages()
		if err != nil {
			fmt.Printf("error retrieving image list: %v\n", err)
			return
		}
		if len(*images) == 0 {
			fmt.Printf("No images found\n")
			return
		}
		hasImages := make([]hasImage, 0)
		for _, i := range *images {
			if !contains(i.Regions, currentWorkspace.HASRegion) { // Skip if no region matches
				continue
			}
			hasImages = append(hasImages, hasImage{
				Name:    i.Name,
				ID:      i.ID,
				Regions: strings.Join(i.Regions, ","),
			})
		}
		if len(hasImages) == 0 {
			fmt.Printf("No images for region %s found\n", currentWorkspace.HASRegion)
			return
		}
		prompt := promptui.Select{
			Label:     "Select image to use",
			Items:     hasImages,
			HideHelp:  true,
			Templates: imageSelectTemplate,
			IsVimMode: false,
			Stdout:    &bellSkipper{},
		}
		i, _, err := prompt.Run()
		if err != nil {
			return
		}
		image := hasImages[i].ID

		var resourceTypes = []struct {
			Name string
		}{
			{"g3s.xlarge"},
			{"g3.4xlarge"},
			{"g3.8xlarge"},
			{"g3.16xlarge"},
		}

		prompt = promptui.Select{
			Label:     "Select resource type to create",
			Items:     resourceTypes,
			Templates: resourceTypeSelectTemplate,
			HideHelp:  true,
			IsVimMode: false,
			Stdout:    &bellSkipper{},
		}
		i, _, err = prompt.Run()
		if err != nil {
			return
		}
		resourceType := resourceTypes[i].Name
		resources, resp, err := client.Resources.CreateResource(has.Resource{
			ImageID:      image,
			ResourceType: resourceType,
			Region:       currentWorkspace.HASRegion,
			Count:        count,
			ClusterTag:   "created-with-hs",
			EBS: has.EBS{
				DeleteOnTermination: true,
				VolumeSize:          50,
				VolumeType:          "standard",
			},
		})
		if err != nil {
			fmt.Printf("failed to create resources: %v\n", err)
			return
		}
		if resources == nil {
			fmt.Printf("failed to create resource: %v\n", resp)
			return
		}
		for _, r := range *resources {
			fmt.Printf("resource %s created, state: %s\n", r.ResourceID, r.State)
		}
	},
}

func init() {
	hasResourcesCmd.AddCommand(hasResourcesCreateCmd)

	hasResourcesCreateCmd.Flags().IntP("count", "c", 1, "Number of resources to create")
}
