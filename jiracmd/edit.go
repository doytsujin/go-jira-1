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

type EditOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.IssueUpdate  `yaml:",inline" json:",inline" figtree:",inline"`
	jira.SearchOptions    `yaml:",inline" json:",inline" figtree:",inline"`
	Overrides             map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Issue                 string            `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdEditRegistry(o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := EditOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("edit"),
		},
		Overrides: map[string]string{},
	}

	return &jiracli.CommandRegistryEntry{
		"Edit issue details",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEditUsage(cmd, &opts)
		},
		func(globals *jiracli.GlobalOptions) error {
			return CmdEdit(o, globals, &opts)
		},
	}
}

func CmdEditUsage(cmd *kingpin.CmdClause, opts *EditOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("query", "Jira Query Language (JQL) expression for the search to edit multiple issues").Short('q').StringVar(&opts.Query)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = jiracli.FlagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	cmd.Arg("ISSUE", "issue id to edit").StringVar(&opts.Issue)
	return nil
}

// Edit will get issue data and send to "edit" template
func CmdEdit(o *oreo.Client, globals *jiracli.GlobalOptions, opts *EditOptions) error {
	type templateInput struct {
		*jiradata.Issue `yaml:",inline"`
		Meta            *jiradata.EditMeta `yaml:"meta" json:"meta"`
		Overrides       map[string]string  `yaml:"overrides" json:"overrides"`
	}
	if opts.Issue != "" {
		issueData, err := jira.GetIssue(o, globals.Endpoint.Value, opts.Issue, nil)
		if err != nil {
			return err
		}
		editMeta, err := jira.GetIssueEditMeta(o, globals.Endpoint.Value, opts.Issue)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue:     issueData,
			Meta:      editMeta,
			Overrides: opts.Overrides,
		}
		err = jiracli.EditLoop(&opts.CommonOptions, &input, &issueUpdate, func() error {
			return jira.EditIssue(o, globals.Endpoint.Value, opts.Issue, &issueUpdate)
		})
		if err != nil {
			return err
		}
		fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, globals.Endpoint.Value, opts.Issue)

		if opts.Browse.Value {
			return CmdBrowse(globals, opts.Issue)
		}
	}
	results, err := jira.Search(o, globals.Endpoint.Value, opts)
	if err != nil {
		return err
	}
	for _, issueData := range results.Issues {
		editMeta, err := jira.GetIssueEditMeta(o, globals.Endpoint.Value, issueData.Key)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue: issueData,
			Meta:  editMeta,
		}
		err = jiracli.EditLoop(&opts.CommonOptions, &input, &issueUpdate, func() error {
			return jira.EditIssue(o, globals.Endpoint.Value, issueData.Key, &issueUpdate)
		})
		if err != nil {
			return err
		}
		fmt.Printf("OK %s %s/browse/%s\n", issueData.Key, globals.Endpoint.Value, issueData.Key)

		if opts.Browse.Value {
			return CmdBrowse(globals, issueData.Key)
		}
	}
	return nil
}
