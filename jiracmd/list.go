package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ListOptions struct {
	jiracli.GlobalOptions      `yaml:",inline" json:",inline" figtree:",inline"`
	jira.SearchOptions `yaml:",inline" json:",inline" figtree:",inline"`
}

func CmdListRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := ListOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("list"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Prints list of issues for given search criteria",
		func() error {
			return CmdList(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			if opts.MaxResults == 0 {
				opts.MaxResults = 500
			}
			if opts.QueryFields == "" {
				opts.QueryFields = "assignee,created,priority,reporter,status,summary,updated"
			}
			if opts.Sort == "" {
				opts.Sort = "priority asc, key"
			}
			return CmdListUsage(cmd, &opts)
		},
	}
}

func CmdListUsage(cmd *kingpin.CmdClause, opts *ListOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("assignee", "User assigned the issue").Short('a').StringVar(&opts.Assignee)
	cmd.Flag("component", "Component to search for").Short('c').StringVar(&opts.Component)
	cmd.Flag("issuetype", "Issue type to search for").Short('i').StringVar(&opts.IssueType)
	cmd.Flag("limit", "Maximum number of results to return in search").Short('l').IntVar(&opts.MaxResults)
	cmd.Flag("project", "Project to search for").Short('p').StringVar(&opts.Project)
	cmd.Flag("query", "Jira Query Language (JQL) expression for the search").Short('q').StringVar(&opts.Query)
	cmd.Flag("queryfields", "Fields that are used in \"list\" template").Short('f').StringVar(&opts.QueryFields)
	cmd.Flag("reporter", "Reporter to search for").Short('r').StringVar(&opts.Reporter)
	cmd.Flag("sort", "Sort order to return").Short('s').StringVar(&opts.Sort)
	cmd.Flag("watcher", "Watcher to search for").Short('w').StringVar(&opts.Watcher)
	return nil
}

// List will query jira and send data to "list" template
func CmdList(o *oreo.Client, opts *ListOptions) error {
	data, err := jira.Search(o, opts.Endpoint.Value, opts)
	if err != nil {
		return err
	}
	return jiracli.RunTemplate(opts.Template.Value, data, nil)
}