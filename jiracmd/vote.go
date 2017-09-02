package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type VoteAction int

const (
	VoteUP VoteAction = iota
	VoteDown
)

type VoteOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue                 string     `yaml:"issue,omitempty" json:"issue,omitempty"`
	Action                VoteAction `yaml:"-" json:"-"`
}

func CmdVoteRegistry(o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := VoteOptions{
		CommonOptions: jiracli.CommonOptions{},
		Action:        VoteUP,
	}

	return &jiracli.CommandRegistryEntry{
		"Vote up/down an issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdVoteUsage(cmd, &opts)
		},
		func(globals *jiracli.GlobalOptions) error {
			return CmdVote(o, globals, &opts)
		},
	}
}

func CmdVoteUsage(cmd *kingpin.CmdClause, opts *VoteOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Flag("down", "downvote the issue").Short('d').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Action = VoteDown
		return nil
	}).Bool()
	cmd.Arg("ISSUE", "issue id to vote").StringVar(&opts.Issue)
	return nil
}

// Vote will up/down vote an issue
func CmdVote(o *oreo.Client, globals *jiracli.GlobalOptions, opts *VoteOptions) error {
	if opts.Action == VoteUP {
		if err := jira.IssueAddVote(o, globals.Endpoint.Value, opts.Issue); err != nil {
			return err
		}
	} else {
		if err := jira.IssueRemoveVote(o, globals.Endpoint.Value, opts.Issue); err != nil {
			return err
		}
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, globals.Endpoint.Value, opts.Issue)

	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}
	return nil
}
