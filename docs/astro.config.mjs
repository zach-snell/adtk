// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightLlmsTxt from 'starlight-llms-txt';

export default defineConfig({
	site: 'https://zach-snell.github.io',
	base: '/adtk',
	integrations: [
		starlight({
			title: 'adtk',
			description: 'A dual-mode Go CLI & MCP Server for Azure DevOps',
			plugins: [
				starlightLlmsTxt({
					projectName: 'adtk (Azure DevOps Toolkit)',
					description: 'A dual-mode Go CLI and MCP Server for Azure DevOps. Provides 13 MCP tools with 82 actions for work items, repos, PRs, pipelines, wiki, boards, iterations, search, test plans, advanced security, projects, users, and attachments. Features multi-base-URL routing, PAT auth, WIQL 2-step pattern, response flattening, and rate limiting.',
					customSets: [
						{
							label: 'MCP Tools',
							description: 'All MCP tool reference documentation',
							paths: ['mcp/**'],
						},
						{
							label: 'CLI',
							description: 'CLI command reference',
							paths: ['cli/**'],
						},
					],
				}),
			],
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/zach-snell/adtk' },
			],
			editLink: {
				baseUrl: 'https://github.com/zach-snell/adtk/edit/main/docs/',
			},
			customCss: ['./src/styles/custom.css'],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Configuration', slug: 'getting-started/configuration' },
						{ label: 'Quick Start', slug: 'getting-started/quickstart' },
					],
				},
				{
					label: 'CLI Commands',
					items: [
						{ label: 'Overview', slug: 'cli/overview' },
						{ label: 'adtk auth', slug: 'cli/auth' },
						{ label: 'adtk mcp', slug: 'cli/mcp' },
						{ label: 'adtk work-items', slug: 'cli/work-items' },
						{ label: 'adtk repos', slug: 'cli/repos' },
						{ label: 'adtk pull-requests', slug: 'cli/pull-requests' },
						{ label: 'adtk pipelines', slug: 'cli/pipelines' },
						{ label: 'adtk wiki', slug: 'cli/wiki' },
						{ label: 'adtk search', slug: 'cli/search' },
						{ label: 'adtk iterations', slug: 'cli/iterations' },
						{ label: 'adtk boards', slug: 'cli/boards' },
						{ label: 'adtk projects', slug: 'cli/projects' },
						{ label: 'adtk attachments', slug: 'cli/attachments' },
						{ label: 'adtk test-plans', slug: 'cli/test-plans' },
						{ label: 'adtk security', slug: 'cli/security' },
					],
				},
				{
					label: 'MCP Tool Reference',
					items: [
						{ label: 'Overview', slug: 'mcp/overview' },
						{ label: 'manage_work_items', slug: 'mcp/manage-work-items' },
						{ label: 'manage_search', slug: 'mcp/manage-search' },
						{ label: 'manage_repos', slug: 'mcp/manage-repos' },
						{ label: 'manage_pull_requests', slug: 'mcp/manage-pull-requests' },
						{ label: 'manage_projects', slug: 'mcp/manage-projects' },
						{ label: 'manage_users', slug: 'mcp/manage-users' },
						{ label: 'manage_iterations', slug: 'mcp/manage-iterations' },
						{ label: 'manage_boards', slug: 'mcp/manage-boards' },
						{ label: 'manage_wiki', slug: 'mcp/manage-wiki' },
						{ label: 'manage_pipelines', slug: 'mcp/manage-pipelines' },
						{ label: 'manage_attachments', slug: 'mcp/manage-attachments' },
						{ label: 'manage_test_plans', slug: 'mcp/manage-test-plans' },
						{ label: 'manage_advanced_security', slug: 'mcp/manage-advanced-security' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Usage Examples', slug: 'guides/examples' },
						{ label: 'WIQL Guide', slug: 'guides/wiql-guide' },
						{ label: 'Microsoft Parity', slug: 'guides/microsoft-parity' },
					],
				},
				{
					label: 'Advanced',
					items: [
						{ label: 'Architecture', slug: 'advanced/architecture' },
						{ label: 'Security', slug: 'advanced/security' },
						{ label: 'Docker Deployment', slug: 'advanced/docker' },
						{ label: 'Development', slug: 'advanced/development' },
					],
				},
			],
		}),
	],
});
