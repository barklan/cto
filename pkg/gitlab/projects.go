package gitlab

type Project struct {
	Links struct {
		Events        string `json:"events"`
		Issues        string `json:"issues"`
		Labels        string `json:"labels"`
		Members       string `json:"members"`
		MergeRequests string `json:"merge_requests"`
		RepoBranches  string `json:"repo_branches"`
		Self          string `json:"self"`
	} `json:"_links"`
	AllowMergeOnSkippedPipeline bool          `json:"allow_merge_on_skipped_pipeline"`
	AnalyticsAccessLevel        string        `json:"analytics_access_level"`
	Archived                    bool          `json:"archived"`
	AutoCancelPendingPipelines  string        `json:"auto_cancel_pending_pipelines"`
	AutoDevopsDeployStrategy    string        `json:"auto_devops_deploy_strategy"`
	AutoDevopsEnabled           bool          `json:"auto_devops_enabled"`
	AutocloseReferencedIssues   bool          `json:"autoclose_referenced_issues"`
	AvatarURL                   interface{}   `json:"avatar_url"`
	BuildCoverageRegex          string        `json:"build_coverage_regex"`
	BuildGitStrategy            string        `json:"build_git_strategy"`
	BuildTimeout                int64         `json:"build_timeout"`
	BuildsAccessLevel           string        `json:"builds_access_level"`
	CanCreateMergeRequestIn     bool          `json:"can_create_merge_request_in"`
	CiConfigPath                string        `json:"ci_config_path"`
	CiDefaultGitDepth           int64         `json:"ci_default_git_depth"`
	CiForwardDeploymentEnabled  bool          `json:"ci_forward_deployment_enabled"`
	CiJobTokenScopeEnabled      bool          `json:"ci_job_token_scope_enabled"`
	ComplianceFrameworks        []interface{} `json:"compliance_frameworks"`
	ContainerExpirationPolicy   struct {
		Cadence       string `json:"cadence"`
		Enabled       bool   `json:"enabled"`
		KeepN         int64  `json:"keep_n"`
		NameRegex     string `json:"name_regex"`
		NameRegexKeep string `json:"name_regex_keep"`
		NextRunAt     string `json:"next_run_at"`
		OlderThan     string `json:"older_than"`
	} `json:"container_expiration_policy"`
	ContainerRegistryAccessLevel             string      `json:"container_registry_access_level"`
	ContainerRegistryEnabled                 bool        `json:"container_registry_enabled"`
	ContainerRegistryImagePrefix             string      `json:"container_registry_image_prefix"`
	CreatedAt                                string      `json:"created_at"`
	CreatorID                                int64       `json:"creator_id"`
	DefaultBranch                            string      `json:"default_branch"`
	Description                              string      `json:"description"`
	EmailsDisabled                           interface{} `json:"emails_disabled"`
	EmptyRepo                                bool        `json:"empty_repo"`
	ExternalAuthorizationClassificationLabel string      `json:"external_authorization_classification_label"`
	ForkingAccessLevel                       string      `json:"forking_access_level"`
	ForksCount                               int64       `json:"forks_count"`
	HTTPURLToRepo                            string      `json:"http_url_to_repo"`
	ID                                       int64       `json:"id"`
	ImportError                              interface{} `json:"import_error"`
	ImportStatus                             string      `json:"import_status"`
	IssuesAccessLevel                        string      `json:"issues_access_level"`
	IssuesEnabled                            bool        `json:"issues_enabled"`
	JobsEnabled                              bool        `json:"jobs_enabled"`
	KeepLatestArtifact                       bool        `json:"keep_latest_artifact"`
	LastActivityAt                           string      `json:"last_activity_at"`
	LfsEnabled                               bool        `json:"lfs_enabled"`
	MergeMethod                              string      `json:"merge_method"`
	MergeRequestsAccessLevel                 string      `json:"merge_requests_access_level"`
	MergeRequestsEnabled                     bool        `json:"merge_requests_enabled"`
	Name                                     string      `json:"name"`
	NameWithNamespace                        string      `json:"name_with_namespace"`
	Namespace                                struct {
		AvatarURL interface{} `json:"avatar_url"`
		FullPath  string      `json:"full_path"`
		ID        int64       `json:"id"`
		Kind      string      `json:"kind"`
		Name      string      `json:"name"`
		ParentID  interface{} `json:"parent_id"`
		Path      string      `json:"path"`
		WebURL    string      `json:"web_url"`
	} `json:"namespace"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool   `json:"only_allow_merge_if_all_discussions_are_resolved"`
	OnlyAllowMergeIfPipelineSucceeds          bool   `json:"only_allow_merge_if_pipeline_succeeds"`
	OpenIssuesCount                           int64  `json:"open_issues_count"`
	OperationsAccessLevel                     string `json:"operations_access_level"`
	PackagesEnabled                           bool   `json:"packages_enabled"`
	PagesAccessLevel                          string `json:"pages_access_level"`
	Path                                      string `json:"path"`
	PathWithNamespace                         string `json:"path_with_namespace"`
	Permissions                               struct {
		GroupAccess struct {
			AccessLevel       int64 `json:"access_level"`
			NotificationLevel int64 `json:"notification_level"`
		} `json:"group_access"`
		ProjectAccess interface{} `json:"project_access"`
	} `json:"permissions"`
	PrintingMergeRequestLinkEnabled bool          `json:"printing_merge_request_link_enabled"`
	PublicJobs                      bool          `json:"public_jobs"`
	ReadmeURL                       string        `json:"readme_url"`
	RemoveSourceBranchAfterMerge    bool          `json:"remove_source_branch_after_merge"`
	RepositoryAccessLevel           string        `json:"repository_access_level"`
	RequestAccessEnabled            bool          `json:"request_access_enabled"`
	RequirementsEnabled             bool          `json:"requirements_enabled"`
	ResolveOutdatedDiffDiscussions  bool          `json:"resolve_outdated_diff_discussions"`
	RestrictUserDefinedVariables    bool          `json:"restrict_user_defined_variables"`
	RunnersToken                    string        `json:"runners_token"`
	SecurityAndComplianceEnabled    bool          `json:"security_and_compliance_enabled"`
	ServiceDeskAddress              string        `json:"service_desk_address"`
	ServiceDeskEnabled              bool          `json:"service_desk_enabled"`
	SharedRunnersEnabled            bool          `json:"shared_runners_enabled"`
	SharedWithGroups                []interface{} `json:"shared_with_groups"`
	SnippetsAccessLevel             string        `json:"snippets_access_level"`
	SnippetsEnabled                 bool          `json:"snippets_enabled"`
	SquashOption                    string        `json:"squash_option"`
	SSHURLToRepo                    string        `json:"ssh_url_to_repo"`
	StarCount                       int64         `json:"star_count"`
	SuggestionCommitMessage         string        `json:"suggestion_commit_message"`
	TagList                         []interface{} `json:"tag_list"`
	Topics                          []interface{} `json:"topics"`
	Visibility                      string        `json:"visibility"`
	WebURL                          string        `json:"web_url"`
	WikiAccessLevel                 string        `json:"wiki_access_level"`
	WikiEnabled                     bool          `json:"wiki_enabled"`
}

// func GetMainProjectName() (string, error) {
// 	project := Project{}
// 	dump, _ := request("")
// 	err := json.Unmarshal(dump, &project)
// 	if err != nil {
// 		return "", err
// 	}
// 	return project.Name, nil
// }
