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
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type ChangelogItemResult struct {
	models.JiraChangelogItem
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string
	AuthorDisplayName string
	Created           time.Time
}

func ConvertChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	sprintIssueConverter, err := NewSprintIssueConverter(taskCtx)
	if err != nil {
		logger.Info(err.Error())
		return err
	}
	db := taskCtx.GetDb()
	logger.Info("covert changelog")
	// select all changelogs belongs to the board
	cursor, err := db.Table("_tool_jira_changelog_items").
		Joins(`left join _tool_jira_changelogs on (
			_tool_jira_changelogs.connection_id = _tool_jira_changelog_items.connection_id
			AND _tool_jira_changelogs.changelog_id = _tool_jira_changelog_items.changelog_id
		)`).
		Joins(`left join _tool_jira_board_issues on (
			_tool_jira_board_issues.connection_id = _tool_jira_changelogs.connection_id
			AND _tool_jira_board_issues.issue_id = _tool_jira_changelogs.issue_id
		)`).
		Select("_tool_jira_changelog_items.*, _tool_jira_changelogs.issue_id, author_account_id, author_display_name, created").
		Where("_tool_jira_changelog_items.connection_id = ? AND _tool_jira_board_issues.board_id = ?", connectionId, boardId).
		Rows()
	if err != nil {
		logger.Info(err.Error())
		return err
	}
	defer cursor.Close()
	issueIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	changelogIdGenerator := didgen.NewDomainIdGenerator(&models.JiraChangelogItem{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		InputRowType: reflect.TypeOf(ChangelogItemResult{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			row := inputRow.(*ChangelogItemResult)
			changelog := &ticket.Changelog{
				DomainEntity: domainlayer.DomainEntity{Id: changelogIdGenerator.Generate(
					row.ConnectionId,
					row.ChangelogId,
					row.Field,
				)},
				IssueId:     issueIdGenerator.Generate(row.ConnectionId, row.IssueId),
				AuthorId:    row.AuthorAccountId,
				AuthorName:  row.AuthorDisplayName,
				FieldId:     row.FieldId,
				FieldName:   row.Field,
				From:        row.FromString,
				To:          row.ToString,
				CreatedDate: row.Created,
			}
			sprintIssueConverter.FeedIn(connectionId, *row)
			return []interface{}{changelog}, nil
		},
	})
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	err = converter.Execute()
	if err != nil {
		return err
	}
	err = sprintIssueConverter.CreateSprintIssue()
	if err != nil {
		return err
	}
	return sprintIssueConverter.SaveAssigneeHistory()
}
