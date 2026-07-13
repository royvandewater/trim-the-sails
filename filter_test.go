package main

import "testing"

func TestFilterReposEmptyNamesReturnsAll(t *testing.T) {
	repos := []string{"/x/alpha.git", "/x/beta.git"}

	got := filterRepos(repos, nil)

	if !equalStrings(got, repos) {
		t.Fatalf("filterRepos(%v, nil) = %v, want %v", repos, got, repos)
	}
}

func TestFilterReposKeepsOnlyNamedRepos(t *testing.T) {
	repos := []string{"/x/alpha.git", "/x/beta.git", "/x/gamma.git"}

	got := filterRepos(repos, []string{"beta.git"})

	want := []string{"/x/beta.git"}
	if !equalStrings(got, want) {
		t.Fatalf("filterRepos = %v, want %v", got, want)
	}
}

func TestFilterReposMatchesNameWithoutDotGitSuffix(t *testing.T) {
	repos := []string{"/x/alpha.git", "/x/beta.git"}

	got := filterRepos(repos, []string{"beta"})

	want := []string{"/x/beta.git"}
	if !equalStrings(got, want) {
		t.Fatalf("filterRepos = %v, want %v", got, want)
	}
}
