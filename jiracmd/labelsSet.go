package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type LabelsSetOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue         string   `yaml:"issue,omitempty" json:"issue,omitempty"`
	Labels        []string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

func CmdLabelsSetRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := LabelsSetOptions{}
	return &jiracli.CommandRegistryEntry{
		"Set labels on an issue",
		func() error {
			return CmdLabelsSet(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdLabelsSetUsage(cmd, &opts)
		},
	}
}

func CmdLabelsSetUsage(cmd *kingpin.CmdClause, opts *LabelsSetOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to set on issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabels will set labels on a given issue
func CmdLabelsSet(o *oreo.Client, opts *LabelsSetOptions) error {
	issueUpdate := jiradata.IssueUpdate{
		Update: jiradata.FieldOperationsMap{
			"labels": jiradata.FieldOperations{
				jiradata.FieldOperation{
					"set": opts.Labels,
				},
			},
		},
	}

	if err := jira.EditIssue(o, opts.Endpoint.Value, opts.Issue, &issueUpdate); err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}