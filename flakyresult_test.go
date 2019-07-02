// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flakyresult_test

import (
	"encoding/base64"
	"context"
	// "encoding/json"

	// "io/ioutil"
	// "strconv"
	// "strings"
	"testing"
	//"encoding/base64"
	"fmt"
	"time"

	"reflect"

	// "cloud.google.com/go/storage"
	"cloud.google.com/go/spanner"
	// "google.golang.org/api/iterator"

	//database "cloud.google.com/go/spanner/admin/database/apiv1"

	"istio.io/bots/policybot/pkg/storage"
	"istio.io/bots/policybot/pkg/testresults"
	"istio.io/bots/policybot/pkg/testflakes"
	"istio.io/bots/policybot/pkg/config"
	"istio.io/pkg/env"
	"istio.io/bots/policybot/pkg/gh"
	"istio.io/bots/policybot/pkg/storage/cache"
	span "istio.io/bots/policybot/pkg/storage/spanner"

	// "istio.io/bots/policybot/pkg/config"
	// span "istio.io/bots/policybot/pkg/storage/spanner"
	// "istio.io/pkg/log"
)

// func read() (client client.Client, ctx context.Context, []*storage.TestResult, error) {
// 	iter := client.Single().Read(ctx, "TestResults", spanner.AllKeys(),
// 		[]string{"OrgID", "RepoID", "TestName", "PrNum", "RunNum", "StartTime",
// 			"FinishTime", "TestPassed", "CloneFailed", "Sha", "Result", "BaseSha", "RunPath"})
// 	defer iter.Stop()
// 	testResults := []*storage.TestResult{}
// 	for {
// 		row, err := iter.Next()
// 		if err == iterator.Done {
// 			//scope.Infof("finished reading")
// 			return testResults, nil
// 		}
// 		if err != nil {
// 			fmt.Println("read err 1")
// 			fmt.Println(err)
// 			return nil, err
// 		}
// 		testResult := &storage.TestResult{}
// 		err = row.ToStruct(testResult)
// 		if err != nil {
// 			fmt.Println("read err 2")
// 			fmt.Println(err)
// 			return nil, err
// 		}
// 		// if err := row.Columns(&testResult.OrgID, &testResult.RepoID, &testResult.TestName, &testResult.PrNum, &testResult.RunNum,
// 		// 	&testResult.StartTime, &testResult.FinishTime, &testResult.TestPassed,
// 		// 	&testResult.CloneFailed, &testResult.Sha, &testResult.Result, &testResult.BaseSha, &testResult.RunPath); err != nil {
// 		// 	return nil, err
// 		// }
// 		testResults = append(testResults, testResult)
// 	}
// }


func TestResults(t *testing.T) {
	context := context.Background()
	const layout = "1/2/2006 15:04:05"
	time1, _ := time.Parse(layout, "11/16/2018 07:03:22")
	t1 := time1.Local()
	time2, _ := time.Parse(layout, "11/16/2018 07:15:44")
	t2 := time2.Local()
	var correctInfo = &storage.TestResult{
	OrgID:       "MDEyOk9yZ2FuaXphdGlvbjIzNTM0NjQ0",
	RepoID:      "MDEwOlJlcG9zaXRvcnk3NDE3NTgwNQ==",
	TestName:    "release-test",
	PrNum:       110,
	RunNum:      155,
	StartTime:   t1,
	FinishTime:  t2,
	TestPassed:  true,
	CloneFailed: false,
	Sha:         "fee4aae74eb4debaf621d653abe8bfcf0ce6a4ea",
	Result:      "SUCCESS",
	BaseSha:     "d995c19aefe6b5ff0748b783e8b69c59963bc8ae",
	RunPath:     "pr-logs/pull/istio_istio/110/release-test/155/",
	}
	orgID := "MDEyOk9yZ2FuaXphdGlvbjIzNTM0NjQ0"
	repoID := "MDEwOlJlcG9zaXRvcnk3NDE3NTgwNQ=="
	var prNum int64 = 110

	prResultTest, err := testresults.NewPrResultTester(context, "istio-flakey-test")
	if err != nil {
		fmt.Println(err)
	return
	}
	testResults, _ := prResultTest.CheckTestResultsForPr(prNum, orgID, repoID)
	test := testResults[0]
	if !reflect.DeepEqual(test, correctInfo) {
	t.Fail()
	}
}

func TestFlakes(t *testing.T) {
	const (
		githubWebhookSecret     = "Secret for the GitHub webhook"
		githubToken             = "Token to access the GitHub API"
		gcpCreds                = "Base64-encoded credentials to access GCP"
		configRepo              = "GitHub org/repo/branch where to fetch policybot config"
		configFile              = "Path to a configuration file"
		sendgridAPIKey          = "API Key for sendgrid.com"
		zenhubToken             = "Token to access the ZenHub API"
		port                    = "TCP port to listen to for incoming traffic"
		githubOAuthClientSecret = "Client secret for GitHub OAuth2 flow"
		githubOAuthClientID     = "Client ID for GitHub OAuth2 flow"
		httpsOnly               = "Send https redirect if x-forwarded-header is not set"
		spannerDatabase         = "the name of the database to sync to"
	)
	// ca := config.DefaultArgs()
	// ca.StartupOptions.GitHubWebhookSecret = env.RegisterStringVar("GITHUB_WEBHOOK_SECRET", ca.StartupOptions.GitHubWebhookSecret, githubWebhookSecret).Get()
	// ca.StartupOptions.GitHubToken = env.RegisterStringVar("GITHUB_TOKEN", ca.StartupOptions.GitHubToken, githubToken).Get()
	// ca.StartupOptions.ZenHubToken = env.RegisterStringVar("ZENHUB_TOKEN", ca.StartupOptions.ZenHubToken, zenhubToken).Get()
	// ca.StartupOptions.GCPCredentials = env.RegisterStringVar("GCP_CREDS", ca.StartupOptions.GCPCredentials, gcpCreds).Get()
	// ca.StartupOptions.ConfigRepo = env.RegisterStringVar("CONFIG_REPO", ca.StartupOptions.ConfigRepo, configRepo).Get()
	// ca.StartupOptions.ConfigFile = env.RegisterStringVar("CONFIG_FILE", ca.StartupOptions.ConfigFile, configFile).Get()
	// ca.StartupOptions.SendGridAPIKey = env.RegisterStringVar("SENDGRID_APIKEY", ca.StartupOptions.SendGridAPIKey, sendgridAPIKey).Get()
	// ca.StartupOptions.Port = env.RegisterIntVar("PORT", ca.StartupOptions.Port, port).Get()
	// ca.StartupOptions.GitHubOAuthClientSecret =
	// 	env.RegisterStringVar("GITHUB_OAUTH_CLIENT_SECRET", ca.StartupOptions.GitHubOAuthClientSecret, githubOAuthClientSecret).Get()
	// ca.StartupOptions.GitHubOAuthClientID =
	// 	env.RegisterStringVar("GITHUB_OAUTH_CLIENT_ID", ca.StartupOptions.GitHubOAuthClientID, githubOAuthClientID).Get()
	ctx := context.Background()
	// ght := gh.NewThrottledClient(ctx, ca.StartupOptions.GitHubToken)
	// fmt.Println("ght")
	// fmt.Println(ght)


	ca := config.DefaultArgs()

	ca.StartupOptions.GitHubToken = env.RegisterStringVar("GITHUB_TOKEN", ca.StartupOptions.GitHubToken, githubToken).Get()
	ca.StartupOptions.GCPCredentials = env.RegisterStringVar("GCP_CREDS", ca.StartupOptions.GCPCredentials, gcpCreds).Get()
	ca.StartupOptions.ConfigRepo = env.RegisterStringVar("CONFIG_REPO", ca.StartupOptions.ConfigRepo, configRepo).Get()
	ca.StartupOptions.ConfigFile = env.RegisterStringVar("CONFIG_FILE", ca.StartupOptions.ConfigFile, configFile).Get()
	ca.SpannerDatabase = env.RegisterStringVar("spannerDatabase", ca.SpannerDatabase, spannerDatabase).Get()
// syncerCmd.PersistentFlags().StringVarP(&ca.SpannerDatabase, "spannerDatabase", "", ca.SpannerDatabase, spannerDatabase)
	ca.SpannerDatabase = "projects/istio-testing/instances/istio-policy-bot/databases/gh"
	ca.StartupOptions.ConfigFile = "/clyu/Downloads/istio-testing-d0fd59c7878c.json"
	fmt.Println("ca")
	fmt.Println(ca)
	client, err := spanner.NewClient(ctx, "projects/istio-testing/instances/istio-policy-bot/databases/gh")
	
	fmt.Println("aaa")
	// client, err := spanner.NewClient(ctx, "projects/"+project+"/instances/"+instance+"/databases/"+database)
	if err != nil {
		fmt.Println("create new spanner client")
		fmt.Println(err)
	}
	fmt.Println("spanner client")
	fmt.Println(client)


	creds, err := base64.StdEncoding.DecodeString(ca.StartupOptions.GCPCredentials)
	if err != nil {
		fmt.Errorf("unable to decode GCP credentials: %v", err)
		t.Fail()
	}
	fmt.Println("creds")
	fmt.Println(creds)
	ght := gh.NewThrottledClient(context.Background(), ca.StartupOptions.GitHubToken)
	fmt.Println(ght)
		store, err := span.NewStore(context.Background(), ca.SpannerDatabase, creds)
	if err != nil {
		fmt.Errorf("unable to create storage layer: %v", err)
		t.Fail()
	}
	// defer store.close()
	fmt.Println("store store store")
	fmt.Println(store)

	cache := cache.New(store, ca.CacheTTL)
	fmt.Println(cache)
	
	// var testResults []*storage.TestResult

	// creds, err := base64.StdEncoding.DecodeString(ca.StartupOptions.GCPCredentials)
	// if err != nil {
	// 	fmt.Println("fail to create creds")
	// 	fmt.Println(err)
	// 	t.Fail()
	// }
	// fmt.Println("creds")
	// fmt.Println(creds)
	// ca.SpannerDatabase = "projects/istio-testing/instances/istio-policy-bot/databases/gh"
	// store, err := span.NewStore(ctx, ca.SpannerDatabase, creds)
	// if err != nil {
	// 	fmt.Println("fail to create store")
	// 	fmt.Println(err)
	// 	t.Fail()
	// }
	// fmt.Println("store")
	// fmt.Println(store)
	// cache := cache.New(store, ca.CacheTTL)

	flakeTester, err := testflakes.NewFlakeTester(ctx, nil, nil, client, "TestResults")
	if err != nil {
		fmt.Println("err 1")
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println("flake tester")
	fmt.Println(flakeTester)
	// // orgID := "MDEyOk9yZ2FuaXphdGlvbjIzNTM0NjQ0"
 //    // repoID := "MDEwOlJlcG9zaXRvcnk3NDE3NTgwNQ"

	// err = store.QueryTestResultByPrNumber(flakeTester.Ctx, 110, func(testResult *storage.TestResult) error {
	// 	testResults = append(testResults, testResult)
	// 	return nil
	// })
	// fmt.Println("err 1")
	// fmt.Println(err)
	// fmt.Println(testResults)

	// testResults, err = flakeTester.read

	// testResults, err = flakeTester.Store.ReadAllTestResults(flakeTester.Ctx)

	// if err != nil {
	// 	fmt.Println("err 2")
	// 	fmt.Println(err)
	// }
	// resultMap := flakeTester.ProcessResults(testResults)
	// flakyResults := flakeTester.CheckResults(resultMap)
	// fmt.Println(flakyResults)
}
