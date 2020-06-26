/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/git"
	"github.com/wso2/product-apim-tooling/import-export-cli/specs/params"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"strconv"
)

var flagVCSStatusEnvName string           // name of the environment to be added

// push command related usage Info
const vcsStatusCmdLiteral = "status"
const vcsStatusCmdShortDesc = "Gets the status report of project changes of the specified environment"
const vcsStatusCmdLongDesc = ``

const vcsStatusCmdCmdExamples = utils.ProjectName + ` ` + vcsStatusCmdLiteral + ` `  + ` -e dev`

// pushCmd represents the push command
var VCSStatusCmd = &cobra.Command{
	Use:     vcsStatusCmdLiteral,
	Short:   vcsStatusCmdShortDesc,
	Long:    vcsStatusCmdLongDesc,
	Example: vcsStatusCmdCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + vcsStatusCmdLiteral + " called")
		totalProjectsToUpdate, updatedProjectsPerType := git.GetStatus(flagVCSStatusEnvName, git.FromRevTypeLastAttempted)

		if totalProjectsToUpdate == 0 {
			fmt.Println("Everything is up-to-date")
			return
		}

		fmt.Println("Projects to Update (" + strconv.Itoa(totalProjectsToUpdate) + ")");
		printProjectsToUpdate(utils.ProjectTypeApi, updatedProjectsPerType[utils.ProjectTypeApi])
		printProjectsToUpdate(utils.ProjectTypeApiProduct, updatedProjectsPerType[utils.ProjectTypeApiProduct])
		printProjectsToUpdate(utils.ProjectTypeApplication, updatedProjectsPerType[utils.ProjectTypeApplication])
	},
}

func printProjectsToUpdate(projectType string, projects []*params.ProjectParams) {
	if len(projects) != 0 {
		fmt.Println("\n" + projectType + "s (" + strconv.Itoa(len(projects)) + ") ...")
		for i, projectParam := range projects {
			var operation string
			var failed string
			if projectParam.Deleted {
				operation = "[delete]"
			} else {
				operation = "[save]"
			}
			if projectParam.FailedDuringPreviousPush {
				failed = "[failed]"
			}
			fmt.Println(strconv.Itoa(i+1) + ": " + operation + "\t" + failed + "\t" + projectParam.NickName +
				": (" + projectParam.RelativePath + ")")
		}
	}
}

func init() {
	VCSCmd.AddCommand(VCSStatusCmd)

	VCSStatusCmd.Flags().StringVarP(&flagVCSStatusEnvName, "environment", "e", "", "Name of the " +
		"environment to check the project(s) status")

	_ = VCSStatusCmd.MarkFlagRequired("environment")
}

/*
Projects to Update (5)...

APIs (3) ...
1: [+/-]	PizzaShackAPI-1.0.0:	(apis/PizzaShackAPI-1.0.0)
2: [+/-]	Sample-1.0.3:	(	apis/Sample-1.0.3)
3: [failed]	SwaggerPetstoreYAML-1.0.0:	(apis/SwaggerPetstoreYAML-1.0.0)



 */