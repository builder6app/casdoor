// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package routers

import (
	"fmt"
	"strings"

	"github.com/beego/beego/context"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/casvisor/casvisor-go-sdk/casvisorsdk"
)

func getUser(ctx *context.Context) (username string) {
	defer func() {
		if r := recover(); r != nil {
			username = getUserByClientIdSecret(ctx)
		}
	}()

	username = ctx.Input.Session("username").(string)

	if username == "" {
		username = getUserByClientIdSecret(ctx)
	}

	return
}

func getUserByClientIdSecret(ctx *context.Context) string {
	clientId := ctx.Input.Query("clientId")
	clientSecret := ctx.Input.Query("clientSecret")
	if clientId == "" || clientSecret == "" {
		return ""
	}

	application, err := object.GetApplicationByClientId(clientId)
	if err != nil {
		panic(err)
	}

	if application == nil || application.ClientSecret != clientSecret {
		return ""
	}

	return util.GetId(application.Organization, application.Name)
}

func RecordMessage(ctx *context.Context) {
	if ctx.Request.URL.Path == "/api/login" || ctx.Request.URL.Path == "/api/signup" {
		return
	}

	userId := getUser(ctx)

	ctx.Input.SetParam("recordUserId", userId)
}

func AfterRecordMessage(ctx *context.Context) {
	fmt.Println("AfterRecordMessage====>")
	record, err := object.NewRecord(ctx)
	fmt.Println("AfterRecordMessage==NewRecord==>")
	if err != nil {
		fmt.Printf("AfterRecordMessage() error: %s\n", err.Error())
		return
	}

	urlPath := getUrlPath(ctx.Request.URL.Path)

	if strings.HasPrefix(urlPath, "/api/notify-payment") {
		urlPath = "/api/notify-payment"
		record.Action = "notify-payment"
	}

	if strings.HasPrefix(urlPath, "/api/invoice-payment") {
		urlPath = "/api/invoice-payment"
		record.Action = "invoice-payment"
	}

	userId := ctx.Input.Params()["recordUserId"]
	fmt.Printf("AfterRecordMessage==userId:%s==>\n", userId)
	if userId != "" {
		record.Organization, record.User = util.GetOwnerAndNameFromId(userId)
	}

	var record2 *casvisorsdk.Record
	recordSignup := ctx.Input.Params()["recordSignup"]
	fmt.Printf("AfterRecordMessage==recordSignup:%s==>\n", recordSignup)
	if recordSignup == "true" {
		record2 = object.CopyRecord(record)
		record2.Action = "new-user"

		var user *object.User
		user, err = object.GetUser(userId)
		if err != nil {
			fmt.Printf("AfterRecordMessage() error: %s\n", err.Error())
			return
		}
		if user == nil {
			err = fmt.Errorf("the user: %s is not found", userId)
			fmt.Printf("AfterRecordMessage() error: %s\n", err.Error())
			return
		}

		record2.Object = util.StructToJson(user)
	}
	fmt.Printf("AfterRecordMessage==record:%s==>\n", record.Action)
	util.SafeGoroutine(func() {
		object.AddRecord(record)
		fmt.Printf("AfterRecordMessage==AddRecord:%s==>\n", record.Action)
		if record2 != nil {
			fmt.Printf("AfterRecordMessage==AddRecord2:%s==>\n", record2.Action)
			object.AddRecord(record2)
		}
	})
}
