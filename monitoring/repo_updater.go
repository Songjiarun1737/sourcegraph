package main

func RepoUpdater() *Container {
	return &Container{
		Name:        "repo-updater",
		Title:       "Repo Updater",
		Description: "Manages interaction with code hosts, instructs Gitserver to update repositories.",
		Groups: []Group{
			{
				Title: "General",
				Rows: []Row{
					{
						sharedFrontendInternalAPIErrorResponses("repo-updater"),
					},
				},
			},
			{
				Title:  "Container monitoring (not available on server)",
				Hidden: true,
				Rows: []Row{
					{
						sharedContainerRestarts("repo-updater"),
						sharedContainerMemoryUsage("repo-updater"),
						sharedContainerCPUUsage("repo-updater"),
					},
				},
			},
		},
	}
}
