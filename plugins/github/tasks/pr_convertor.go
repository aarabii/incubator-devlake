/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertPullRequestsMeta = core.SubTaskMeta{
	Name:             "convertPullRequests",
	EntryPoint:       ConvertPullRequests,
	EnabledByDefault: true,
	Description:      "ConvertPullRequests data from Github api",
}

func ConvertPullRequests(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Model(&models.GithubPullRequest{}).Where("repo_id = ?", repoId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	prIdGen := didgen.NewDomainIdGenerator(&models.GithubPullRequest{})
	repoIdGen := didgen.NewDomainIdGenerator(&models.GithubRepo{})
	userIdGen := didgen.NewDomainIdGenerator(&models.GithubUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.GithubPullRequest{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			pr := inputRow.(*models.GithubPullRequest)
			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: prIdGen.Generate(pr.GithubId),
				},
				BaseRepoId:     repoIdGen.Generate(pr.RepoId),
				Status:         pr.State,
				Title:          pr.Title,
				Url:            pr.Url,
				AuthorId:       userIdGen.Generate(pr.AuthorId),
				AuthorName:     pr.AuthorName,
				Description:    pr.Body,
				CreatedDate:    pr.GithubCreatedAt,
				MergedDate:     pr.MergedAt,
				ClosedDate:     pr.ClosedAt,
				Key:            pr.Number,
				Type:           pr.Type,
				Component:      pr.Component,
				MergeCommitSha: pr.MergeCommitSha,
				BaseRef:        pr.BaseRef,
				BaseCommitSha:  pr.BaseCommitSha,
				HeadRef:        pr.HeadRef,
				HeadCommitSha:  pr.HeadCommitSha,
			}
			return []interface{}{
				domainPr,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
