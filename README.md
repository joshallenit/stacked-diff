# Developer Scripts for Stacked Diff Workflow

These scripts make it easier to build from the command line and to create and update PR's with Github. They facilitates a [stacked diff workflow](https://newsletter.pragmaticengineer.com/p/stacked-diffs), where you always commit on `main` branch and have can have multiple streams of work all on `main`.

## TL;DR

Using a stacked diff workflow like this allows you to work on separate streams of work without changing branches. This project is a Command Line Interface that manages git commits and branches to allow you to quickly use a stacked diff workflow. It uses the Github CLI to create pull requests and add reviewers once PR checks have passed.

## Installation

Clone the repository or download the [latest release](https://github.com/tinyspeck/stacked-diff-workflow/releases), and then:

### Mac

**Optional:** As this is a CLI, do yourself a favor and install [iTerm](https://iterm2.com/) and [zsh](https://ohmyz.sh/), as it makes working for the command line more pleasant.

```bash
# Install Github CLI
brew install gh 
# Setup login for Github CLI
gh auth login 
# Add the /bin directory to your PATH. 
# Replace the directory below to wherever you cloned the repository or unzipped the release.
# For example if using zsh and cloned in your home directory:
echo "export PATH=\$PATH:\$HOME/stacked-diff-workflow/bin" >> ~/.zshrc
source ~/.zshrc
```

### Windows

Install [Git and Git Bash](https://gitforwindows.org/)

## Stacked Diff Workflow CLI

```bash
usage: sd [top-level-flags] <command> [<args>]`

Possible commands are:

   add-reviewers    Add reviewers to Pull Request on Github once its checks have passed
   branch-name      Outputs branch name of commit
   checkout         Checks out branch associated with commit indicator
   code-owners      Outputs code owners for all of the changes in branch
   log              Displays git log of your changes
   new              Create a new pull request from a commit on main
   rebase-main      Bring your main branch up to date with remote
   replace-commit   Replaces a commit on main branch with the contents its associated branch
   update           Add commits from main to an existing PR
   wait-for-merge   Waits for a pull request to be merged

To learn more about a command use: sd <command> --help

flags:

  -log-level string
        Possible log levels:
           debug
           info
           warn
           error
        Default is info, except on commands that are for output purposes,
        (namely branch-name and log), which have a default of error.
```

### Basic Commands

#### log

Displays summary of the git commits on current branch that are not in the remote branch.

Useful to view list indexes, or copy commit hashes, to use for the commitIndicator required by other commands.

A ✅ means that there is a PR associated with the commit (actually it means there is a branch, but having a branch means there is a PR when using this workflow). If there is more than one commit on the associated branch, those commits are also listed (indented under the their associated commit summary).

```bash
usage: sd log
```

<img width="663" alt="image" src="https://user-images.githubusercontent.com/79605685/210386995-9c3e7179-24ed-4d59-9b3e-2b3b34aa6ccc.png">

#### new

Create a new PR with a cherry-pick of the given commit indicator.

This command first creates an associated branch, (with a name based on the commit summary), and then uses Github CLI to create a PR.

Can also add reviewers once PR checks have passed, see "--reviewers" flag.

```bash
usage: sd new [flags] [commitIndicator (default is HEAD commit on main)]

Ticket Number:

If you prefix a (Jira-like formatted) ticket number to the git commit
summary then the "Ticket" section of the PR description will be
populated with it.

For example:

"CONV-9999 Add new feature"

Templates:

The Pull Request Title, Body (aka Description), and Branch Name are
created from golang templates.

The default templates are:

   branch-name.template:      src/go/config/branch-name.template
   pr-description.template:   src/go/config/pr-description.template
   pr-title.template:         src/go/config/pr-title.template

To change a template, copy the default from src/go/config/ into
~/.stacked-diff-workflow/ and modify contents.

The possible values for the templates are:

   CommitBody                   Body of the commit message
   CommitSummary                Summary line of the commit message
   CommitSummaryCleaned         Summary line of the commit message without
                                spaces or special characters
   CommitSummaryWithoutTicket   Summary line of the commit message without
                                the prefix of the ticket number
   FeatureFlag                  Value passed to feature-flag flag
   TicketNumber                 Jira ticket as parsed from the commit summary
   Username                     Name as parsed from git config email.
                                Note: any dots (.) in username are converted to
                                      dashes (-) before being used in
                                      branch-name.template.


flags:

  -base string
        Base branch for Pull Request (default "main")
  -draft
        Whether to create the PR as draft (default true)
  -feature-flag string
        Value for FEATURE_FLAG in PR description
  -indicator string
        Indicator type to use to interpret commitIndicator:
           commit   a commit hash, can be abbreviated,
           pr       a github Pull Request number,
           list     the order of commit listed in the git log, as indicated
                    by "sd log"
           guess    the command will guess the indicator type:
              Number between 0 and 99:       list
              Number between 100 and 999999: pr
              Otherwise:                     commit
         (default "guess")
  -min-checks int
        Minimum number of checks to wait for before verifying that checks
        have passed before adding reviewers. It takes some time for checks
        to be added to a PR by Github, and if you add-reviewers too soon it
        will think that they have all passed. (default 4)
  -reviewers string
        Comma-separated list of Github usernames to add as reviewers once
        checks have passed.
```

<img width="938" alt="image" src="https://user-images.githubusercontent.com/79605685/210406914-9b43f0e0-ac11-498f-bdd7-5a48e07dcbc0.png">

###### Note on Commit Messages

Keep your commit summary to a [reasonable length](https://www.midori-global.com/blog/2018/04/02/git-50-72-rule). The commit summary is used as the branch name. To add more detail use the [commit description](https://stackoverflow.com/questions/40505643/how-to-do-a-git-commit-with-a-subject-line-and-message-body/40506149#40506149). The
created branch name is truncated to 120 chars as Github has problems with very long
branch names.


#### update

Add commits from local main branch to an existing PR.

Can also add reviewers once PR checks have passed, see "--reviewers" flag.

```bash
usage: sd update [flags] <commitIndicator> [fixup commitIndicator (defaults to head commit) [fixup commitIndicator...]]

flags:

  -indicator string
        Indicator type to use to interpret commitIndicator:
           commit   a commit hash, can be abbreviated,
           pr       a github Pull Request number,
           list     the order of commit listed in the git log, as indicated
                    by "sd log"
           guess    the command will guess the indicator type:
              Number between 0 and 99:       list
              Number between 100 and 999999: pr
              Otherwise:                     commit
         (default "guess")
  -min-checks int
        Minimum number of checks to wait for before verifying that checks
        have passed before adding reviewers. It takes some time for checks
        to be added to a PR by Github, and if you add-reviewers too soon it
        will think that they have all passed. (default 4)
  -reviewers string
        Comma-separated list of Github usernames to add as reviewers once
        checks have passed.
```

#### add-reviewers

Add reviewers to Pull Request on Github once its checks have passed.

If PR is marked as a Draft, it is first marked as "Ready for Review".

```bash
usage: sd add-reviewers [flags] [commitIndicator [commitIndicator]...]

flags:

  -indicator string
        Indicator type to use to interpret commitIndicator:
           commit   a commit hash, can be abbreviated,
           pr       a github Pull Request number,
           list     the order of commit listed in the git log, as indicated
                    by "sd log"
           guess    the command will guess the indicator type:
              Number between 0 and 99:       list
              Number between 100 and 999999: pr
              Otherwise:                     commit
         (default "guess")
  -min-checks int
        Minimum number of checks to wait for before verifying that checks
        have passed before adding reviewers. It takes some time for checks
        to be added to a PR by Github, and if you add-reviewers too soon it
        will think that they have all passed. (default 4)
  -poll-frequency duration
        Frequency which to poll checks. For valid formats see https://pkg.go.dev/time#ParseDuration (default 30s)
  -reviewers string
        Comma-separated list of Github usernames to add as reviewers once
        checks have passed.
        Falls back to PR_REVIEWERS environment variable.
  -when-checks-pass
        Poll until all checks pass before adding reviewers (default true)
```

<img width="904" alt="image" src="https://user-images.githubusercontent.com/79605685/210428712-bcea3ce7-e70f-4982-aa54-48e166221a1d.png">

###### Reviewers

You can specify more than one reviewer using a comma-delimited string.

To use the environment variable instead of the "--reviewers" flag:

```bash
export PR_REVIEWERS=first-user,second-user,third-user
```

Add this to your shell rc file (`~/.zshrc` or `~/.bashrc`) and run `source <rc-file>`

### Commands for Rebasing and Fixing Merge Conflicts

#### rebase-main

#### checkout

Checks out the branch associated with commit indicator.

For when you want to merge only the branch with with origin/main, rather than your entire local main branch, verify why CI is failing on that particular branch, or for any other reason.

After modifying the branch you can use "sd replace-commit" to sync local main.

```bash
usage: sd checkout [flags] <commitIndicator>

flags:

  -indicator string
        Indicator type to use to interpret commitIndicator:
           commit   a commit hash, can be abbreviated,
           pr       a github Pull Request number,
           list     the order of commit listed in the git log, as indicated by "sd log"
           guess    the command will guess the indicator type:
              Number between 0 and 99:       list
              Number between 100 and 999999: pr
              Otherwise:                     commit
         (default "guess")
```

#### replace-commit

Replaces a commit on main branch with the squashed contents of its associated branch.

This is useful when you make changes within a branch, for example to fix a problem found on CI, and want to bring the changes over to your local main branch.

```bash
usage: sd replace-commit [flags] <commitIndicator>

flags:

  -indicator string
        Indicator type to use to interpret commitIndicator:
           commit   a commit hash, can be abbreviated,
           pr       a github Pull Request number,
           list     the order of commit listed in the git log, as indicated
                    by "sd log"
           guess    the command will guess the indicator type:
              Number between 0 and 99:       list
              Number between 100 and 999999: pr
              Otherwise:                     commit
         (default "guess")
```

### Commands for Custom Scripting

#### branch-name
#### wait-for-merge

### Other Commands

#### code-owners

#### git-prs

*Note this is a stand-alone script and not part of the "sd" cli.*

Lists all of your open PRs. Useful for copying PR numbers.

```bash
usage: git-prs
```

<img width="904" alt="image" src="https://github.com/tinyspeck/stacked-diff-workflow/assets/79605685/7e7a5708-58dc-4060-96b9-89615a86c009">


### To Help You Build

*Note: Only [Android](https://github.com/tinyspeck/slack-android-ng) build scripts are currently supported.*

#### Script: assemble-app

`assemble-app`

Calls `./gradlew assembleInternalDebug` and build tests. Use "-s" (silent) flag to not use voice (`say`) to announce success/failure. Any options after, or in-lieu of, "-s" will be passed to the `./gradle` command, for example `--rerun-tasks`.

#### Script: install-app

`install-app`

Calls `./gradlew assembleInternalDebug` and install on real device. Use "-s" (silent) flag to not use voice (`say`) to announce success/failure. Any options after, or in-lieu of, "-s" will be passed to the `./gradle` command, for example `--rerun-tasks`.

#### Script: install-apk

`install-apk`

Installs the already compiled APK on a real devices. Useful for after you have run `install-app` but forgot to plugin in your phone 😄. It's faster than running `install-app` again as it doesn't run gradle.

## Example Workflow

### Creating and Updating PRs

Use **sd new** and **sd update** to create and update PR's while always staying on `main` branch.

### To Update Main

Once a PR has been merged, just rebase main normally. The local PR commit will be replaced by the one that Github created when squasing and merging.

```bash
git fetch && git rebase origin/main
```

If you run into conflicts with a commit that has already been merged you can just ignore it. This can happen, for example, if a change was made on github.com and it is not reflected in your local commit. Obviously, only do this if the PR has actually already been merged into main! The error message from rebase will let you know which commit has conflicts.

```bash
git reset --hard head && git rebase --continue
```

This process has been automated by the `sd rebase-main` command.

#### To Fix Merge Conflicts

##### Easy Flow

If you just are rebasing with `main` and the commit with merge conflict has already been **merged**, then the process is simpler.

1. Fix Merge Conflict

```bash
# switch to feature branch that has a merge conflict
sd checkout <commitIndicator> 
git fetch && git merge origin/main
# ... and address any merge conflicts
# Update your PR
git push origin/xxx 
```

2\. Merge PR via Github

3\. [Update your Main Branch](#to-update-main)

##### Longer Running Flow

If you want to update your main branch *before* you merge your PR, you can use **replace-head** to keep your local `main` up to date.

```bash
# switch to feature branch that has a merge conflict
sd checkout <commitIndicator> 
# rebase or merge
git fetch && git merge origin/main
# ... and address any merge conflicts
# Update your PR
git push origin/xxx 
# Rebase your local main branch.
git switch main
git rebase origin/main
# hit same merge conflicts, use replace-head to copy the fixes you just made
replace-head <commitIndicator>
# continue with the rebase
git add . && git rebase --continue
# All done... now both the feature branch and your local main are rebased with main, 
# and the merge conflicts only had to be fixed once
```

# Building Source and Contributing

See the [Contributing Developer Guide](DEVELOPER_GUIDE.md), which includes instructions on how to build the source, as well as an overview of the code.

## Stacked Pull Requests?

Note: these scripts do *not* facilitate Stacked *Pull Requests*. Github does some things that add friction to using Stacked PR's, even with support from third party software. For example, after merging one of the PR's in the stack, the other PR's will require a re-review. Instead of Stacked PRs, it's recommended to organize your PR's, as much as reasonably possible, so that they can be all be rebased against main at the same time. When there are dependencies, wait for dependant PR to be merged before putting up the next one. You may find that often you are still working on the next commit while the other is being reviewed/merged.

## Contact

Join the discussion in [#devel-stacked-diff-workflow](https://slack-pde.slack.com/archives/C03V94N2A84)

## Acknowledgments

- Thanks to for publishing this article that inspired the first version of the scripts 
https://kastiglione.github.io/git/2020/09/11/git-stacked-commits.html

- Thanks to the Github team for creating a CLI that is leveraged here.


Good example

https://github.com/helm/helm

https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-github-profile/customizing-your-profile/managing-your-profile-readme


	// so I could do something like this:
	// https://github.com/Nutlope/aicommits
	// to create an AI git commit message

	/*
					https://www.hatica.io/blog/ai-commit-tools/
				https://github.com/kamushadenes/chloe/blob/main/.github/scripts/release-notes.py
		the commit messages don't look that great, probably not the best idea
		but splitting commits into smaller ones based on compilability could be good. Maybe figuring out which dependencies are declared in which file and then which dependencies are used in each file? That could work


				I could actually try these commit messages to see how well they work or don't work
				Other ideas:
					- show git log very fast and then check for branches async using ANSI codes
					- output the latest head commit so that a rollback can be done "git reset --hard xxx", but that doesn't keep track of what it does to branches which could screw up a PR, undo? that could be useful.
					- Include PR link instead of just a :green-check: beside PR's, there might be a way to http link from terminal. Too long to do it for every branch, would require the ANSI backspace idea. Or it could run in the background and update it next? Not sure if Go can launch background tasks and terminate, but probably? Or do it with a visual gui? wow, that would be something. https://github.com/charmbracelet/bubbletea
					- Next `gh pr view 83824 --json latestReviews` and ensure developer is already not approved so that the review is not dismissed
					- show the output of git commands in a tabbed window that uses ANSI escape codes to move around the screen
					- better error handling so that it reverted on error rather then leaving in an indeterminite state... but wouldn't this mean that I have to save error codes so they can be reported upstream?
					- aliases for reviewer names, so ideally what



POLISH
- update README with latest changes
- code review, review/organize/document code
- standardize all test so they are using assert* and named properly
- test for error conditions (ugh not fun) + but could lead to better error handling (rollback)
- functions only need to be UpperCase if exported from PACKAGE, not file
					

create a PR summary message based on commit messages and comments in code? would that work... probably not.

sending a diff would be too exact

migrate git-merge to use the new merge queue -- mention that it was deleted to Evan in case he wants to do it

NOT WORTH DOING

It's easier to just set envionment variables

sd config set reviewer-alias ankit search
sd config show reviewer-alias




						*/
