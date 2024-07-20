package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/manifoldco/promptui"
)

type branch struct {
	RefName   string
	ShortName string
	ShortHash string
	FullHash  string
}

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		color.Red(fmt.Sprintf("Error Retrieving Path: %v\n", err))

		return
	}

	branches, err := getBranches(currentDir)
	if err != nil {
		color.Red(fmt.Sprintf("Error Retriving .git Data: %v\n", err))

		return
	}

	option, err := listBranchesAndSelectTarget(branches)
	if err != nil {
		color.Red(fmt.Sprintf("Error Selecting Branch: %v\n", err))

		return
	}

	err = performSwitch(currentDir, option)
	if err != nil {
		color.Red(fmt.Sprintf("Error Switching To Selected Branch(%s): %v\n", option, err))

		return
	}

	color.Green(fmt.Sprintf("Switched to Branch %s\n", option.ShortName))
}

func getBranches(path string) ([]branch, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	iter, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	var branches []branch

	err = iter.ForEach(func(c *plumbing.Reference) error {
		branches = append(branches, branch{
			string(c.Name()),
			c.Name().Short(),
			c.Hash().String()[:7],
			c.Hash().String(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func listBranchesAndSelectTarget(options []branch) (branch, error) {
	templates := &promptui.SelectTemplates{
		Label: "{{ . }}",

		Active:   "\U0001F33F {{ .ShortName | cyan }} ({{ .FullHash | red }})",
		Inactive: "  {{ .ShortName | cyan }} ({{ .ShortHash | red }})",
		Selected: "Switching to \U0001F33F {{ .RefName | green}}",
	}

	prompt := promptui.Select{
		Label:     "Select Target Branch",
		Items:     options,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return branch{}, err
	}

	selectedBranch := options[i]

	return selectedBranch, nil
}

func performSwitch(path string, selectedBranch branch) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(selectedBranch.RefName),
		Force:  true,
	})
	if err != nil {
		return err
	}

	return nil
}
