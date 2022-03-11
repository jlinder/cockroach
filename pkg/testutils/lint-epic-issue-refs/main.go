// Copyright 2022 The Cockroach Authors.
// // Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// Check that GitHub PR descriptions and commit messages contain the
// expected epic and issue references.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	//"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lint-epic-issue-refs",
	Short: "lint-epic-issue-refs checks PR body and commit messages for epic and issue refs",
	Args: func(cmd *cobra.Command, args []string) error {
		_, err := parseArgs(args)
		return err
	},
	RunE: func(_ *cobra.Command, args []string) error {
		params := defaultEnvParameters()
		prNumber, err := parseArgs(args)
		if err != nil {
			return err
		}
		params.PrNumber = prNumber
		return lintPR(params)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func parseArgs(args []string) (int, error) {
	if len(args) < 1 {
		return -1, errors.New("No PR number specified")
	}
	if len(args) > 1 {
		return -1, errors.New("One argument is required: a PR number")
	}
	prNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return -1, fmt.Errorf("Invalid PR number: %v", err)
	}

	return prNumber, nil
}

type Parameters struct {
	Token string // GitHub API token
	Org   string
	Repo  string
	PrNumber int
}

func defaultEnvParameters() *Parameters {
	const (
		githubAPITokenEnv      = "GITHUB_API_TOKEN"
		githubOrgEnv           = "GITHUB_ORG"
		githubRepoEnv          = "GITHUB_REPO"
	)

	return &Parameters{
		Token: maybeEnv(githubAPITokenEnv, ""),
		Org:   maybeEnv(githubOrgEnv, "cockroachdb"),
		Repo:  maybeEnv(githubRepoEnv, "cockroach"),
	}
}

func maybeEnv(envKey, defaultValue string) string {
	v := os.Getenv(envKey)
	if v == "" {
		return defaultValue
	}
	return v
}
