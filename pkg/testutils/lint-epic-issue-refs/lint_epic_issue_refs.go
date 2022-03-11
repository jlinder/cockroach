// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.
package main

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	fixIssueRefRE = regexp.MustCompile(`(?im)(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?):?\s+(?:(?:(#\d+)|([\w.-]+[/][\w.-]+#\d+))[,.;]?(?:[ \t]+|$))+`)
	informIssueRefRE = regexp.MustCompile(`(?im)(?:see also|informs):?\s+(?:(?:(#\d+)|([\w.-]+[/][\w.-]+#\d+))[,.;]?(?:[ \t]+|$))+`)
	epicRefRE = regexp.MustCompile(`(?im)epic:?\s+(?:([A-Z]+-[0-9]+)[,.;]?(?:[ \t]+|$))+`)
	githubIssueRefRE = regexp.MustCompile(`(#\d+)|([\w.-]+[/][\w.-]+#\d+)`)
	jiraIssueRefRE = regexp.MustCompile(`[A-Z]+-[0-9]+`)
	releaseNoteCategoryRE = regexp.MustCompile(`(?im)^\s*release notes?\s*\(([^)]+)\):\s+`)
	releaseNoteNoneRE = regexp.MustCompile(`(?im)^\s*release notes?\s*:\s+none\s*`)
)

// TODO(jamesl): rename this struct. It doesn't feel quite right.
type commitInfo struct {
	epicRefs map[string]int
	issueCloseRefs map[string]int
	issueInformRefs map[string]int
	releaseNoteCategories map[string]int
	sha string
}

// releaseNoteRequiresRef checks if any of the release notes in the commit require an epic or
// issue reference.
func (ci *commitInfo) releaseNoteRequiresRef() bool {
	if ci.releaseNoteCategories == nil {
		return false
	}
	releaseNoteRequiresRef := false
	for _, category := range reflect.ValueOf(ci.releaseNoteCategories).MapKeys()  {
		// Commits with these categories don't need a ref
		if category.String() == "None" || category.String() == "bug fix" || category.String() == "bugfix" {
			continue
		}
		releaseNoteRequiresRef = true
	}
	return releaseNoteRequiresRef
}

type prBodyInfo struct {
	epicRefs map[string]int
	issueCloseRefs map[string]int
	issueInformRefs map[string]int
	prNumber int
}

// TODO(jamesl): needs a better name
type responseBuilder struct {
	failures []string
}

// TODO(jamesl): Implement me!
func (rb *responseBuilder) addCommitCheckMessage(message string, prId int, sha string) {
}

// TODO(jamesl): Implement me!
func (rb *responseBuilder) addPrBodyCheckMessage(message string, prId int) {
}

// TODO(jamesl): Implement me!
// constructMessage generates the message for the status check reported back to the PR
func (rb *responseBuilder) constructMessage() string {
	return ""
}

func extractStringsFromMessage(message string, firstMatch, secondMatch *regexp.Regexp) map[string]int {
	ids := map[string]int{}
	if allMatches := firstMatch.FindAllString(message, -1); len(allMatches) > 0 {
		for _, x := range allMatches {
			matches := secondMatch.FindAllString(x, -1)
			for _, match := range matches {
				ids[match] += 1
			}
		}
	}
	return ids
}

func extractFixIssueIDs(message string) map[string]int {
	return extractStringsFromMessage(message, fixIssueRefRE, githubIssueRefRE)
}

func extractInformIssueIDs(message string) map[string]int {
	return extractStringsFromMessage(message, informIssueRefRE, githubIssueRefRE)
}

func extractEpicIDs(message string) map[string]int {
	return extractStringsFromMessage(message, epicRefRE, jiraIssueRefRE)
}

func extractReleaseNoteCategories(message string) map[string]int {
	categories := map[string]int{}
	if allMatches := releaseNoteCategoryRE.FindAllStringSubmatch(message, -1); len(allMatches) > 0 {
		for _, x := range allMatches {
			categories[x[1]] += 1
		}
	}
	if allMatches := releaseNoteNoneRE.FindAllString(message, -1); len(allMatches) > 0 {
		categories["None"] = 1
	}
	return categories
}

func scanCommitsForEpicAndIssueRefs(
	ctx context.Context, ghClient *github.Client, params *Parameters,
) ([]commitInfo, error) {
	commits, _, err := ghClient.PullRequests.ListCommits(
		ctx, params.Org, params.Repo, params.PrNumber, &github.ListOptions{PerPage: 100},
	)
	if err != nil {
		return nil, err
	}

	// TODO(jamesl): add capturing of `Epic: None`
	var infos []commitInfo
	for _, commit := range commits {
		message := commit.GetCommit().GetMessage()
		var info = commitInfo{
			epicRefs:              extractEpicIDs(message),
			issueCloseRefs:        extractFixIssueIDs(message),
			issueInformRefs:       extractInformIssueIDs(message),
			releaseNoteCategories: extractReleaseNoteCategories(message),
			sha:                   commit.GetSHA(),
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// TODO(jamesl): add tests for this function
// TODO(jamesl): update to write into the response builder
func checkCommitsMissingReleaseNoteCategory(commitInfos []commitInfo, builder responseBuilder) []string {
	commitsMissingReleaseNote := []string{}
	for _, info := range commitInfos {
		if info.releaseNoteCategories == nil || len(info.releaseNoteCategories) == 0 {
			commitsMissingReleaseNote = append(commitsMissingReleaseNote, info.sha)
			continue
		}
	}
	return commitsMissingReleaseNote
}

// TODO(jamesl): add tests for this function
// TODO(jamesl): update to write into the response builder
// checkForMissingRefs determines if the PR and its commits has the needed refs.
// When the PR body is missing a ref and one or more commits are missing a ref, tell the
// caller which commits are missing refs.
//
// Returns:
// - list of commit SHAs that are missing refs
func checkForMissingRefs(prInfo prBodyInfo, commitInfos []commitInfo, builder responseBuilder) []string {
	// When the PR body has a ref, no refs are needed in individual commits
	if len(prInfo.epicRefs) + len(prInfo.issueInformRefs) + len(prInfo.issueCloseRefs) > 0 {
		return []string{}
	}

	commitsWithoutRefs := []string{}
	for _, info := range commitInfos {
		if !info.releaseNoteRequiresRef() {
			continue
		}
		if len(info.epicRefs) + len(info.issueInformRefs) + len(info.issueCloseRefs) == 0 {
			commitsWithoutRefs = append(commitsWithoutRefs, info.sha)
		}
	}
	return commitsWithoutRefs
}

// TODO(jamesl): update to write into the response builder
// multipleEpicsInPRBody checks that all epic references in the PR body are used in at least one
// commit message and that the individual commits note which epic they are associated with.
//
// Returns:
// - unusedPrEpicRefs: list of epic refs from the PR body that were not found in a commit
// - commitsMissingEpicRef: list of commit SHAs that don't contain an epic ref.
func checkPrEpicsUsedInCommits(prInfo prBodyInfo, commitInfos []commitInfo, builder responseBuilder) {
	if len(prInfo.epicRefs) < 2 {
		return
	}

	// Check PR body epics are all referenced from commit messages
	for _, prEpic := range reflect.ValueOf(prInfo.epicRefs).MapKeys() {
		found := false
		for _, ci := range commitInfos {
			if ci.epicRefs[prEpic.String()] > 0 {
				found = true
			}
		}
		if !found {
			// TODO(jamesl): add error message to the builder
		}
	}

	// Check all commit messages with a release note reference one of the PR epics
	for _, ci := range commitInfos {
		if !ci.releaseNoteRequiresRef() {
			continue
		}
		if len(ci.epicRefs) == 0 {
			// expected a ref. none set.
			// TODO(jamesl): add error message to the builder
		} else if len(ci.epicRefs) > 1 {
			// expected one ref. multiple set
			// TODO(jamesl): add error message to the builder
		} else {
			for _, epicRef := range reflect.ValueOf(ci.epicRefs).MapKeys() {
				if _, ok := prInfo.epicRefs[epicRef.String()]; !ok {
					// expected the ref to be in the PR body epic refs.
					// TODO(jamesl): add error message to the builder
				}
			}
		}
	}

	return
}

// TODO(jamesl): implement me
// checkMultipleEpicsInCommits checks for commits that contain multiple epic references and adds a
// warning that it is not a common case and to check it is intended.
func checkMultipleEpicsInCommits(commitInfos []commitInfo, builder responseBuilder) {

}

func lintEpicAndIssueRefs(ctx context.Context, ghClient *github.Client, params *Parameters) error {
	// TODO(jamesl): handle backport PRs containing multiple source PRs differently
	// How might they be handled?
	// - check that all source PR bodies have epic or issue refs?
	// - figure out which commits in this PR correspond to the commits in the source PRs and ...
	// - check all source PRs have all the expected references a PR should have. Do something
	//   with that info, like maybe link it back to the commits in this PR and report on it
	//   somehow??? (This seems really heavy / not useful.)
	// - don't check backports with multiple source PRs for the complex cases like
	//   "checkPrEpicsUsedInCommits"?
	// - these checks are really for making it possible for the automated docs issue generation to
	//   get the right epic(s) applied to them. Maybe make this a super simple check and then add
	//   info to the generated docs issue stating it needs to be reviewed more. This case shouldn't
	//   happen too often. The majority of backports are for a single source PR.

	commitInfos, err := scanCommitsForEpicAndIssueRefs(ctx, ghClient, params)
	if err != nil {
		// TODO: handle err.
	}

	pr, _, err := ghClient.PullRequests.Get(
		ctx, params.Org, params.Repo, params.PrNumber,
	)
	if err != nil {
		return fmt.Errorf("Error getting pull requests from GitHub: %v", err)
}
	prBody := pr.GetBody()
	var prInfo = prBodyInfo{
		epicRefs:   extractEpicIDs(prBody),
		issueCloseRefs: extractFixIssueIDs(prBody),
		issueInformRefs: extractInformIssueIDs(prBody),
		prNumber:  params.PrNumber,
	}

	builder := responseBuilder{}

	checkCommitsMissingReleaseNoteCategory(commitInfos, builder)
	checkForMissingRefs(prInfo, commitInfos, builder)
	checkPrEpicsUsedInCommits(prInfo, commitInfos, builder)
	checkMultipleEpicsInCommits(commitInfos, builder)

	// TODO:
	// - generate the status report message
	// - double check that directions for how to fix the reported issue are sufficiently clear
	// - write the status report back to the PR
	return nil
}

func lintPR(params *Parameters) error {
	ctx := context.Background()

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: params.Token},
	)))

	return lintEpicAndIssueRefs(ctx, client, params)
}
