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

package migrationscripts

import (
	"context"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"

	"gorm.io/gorm"
)

type UpdateSchemas20220505 struct{}

func (*UpdateSchemas20220505) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameTable(archived.JiraSource{}, models.JiraConnection{})
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraBoard{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraProject{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraUser{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssue{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraBoardIssue{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraChangelog{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraChangelogItem{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraRemotelink{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssueCommit{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssueTypeMapping{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssueStatusMapping{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraSprint{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraBoardSprint{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraSprintIssue{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraWorklog{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220505) Version() uint64 {
	return 20220505212344
}

func (*UpdateSchemas20220505) Owner() string {
	return "Jira"
}

func (*UpdateSchemas20220505) Name() string {
	return "Rename source to connection "
}
