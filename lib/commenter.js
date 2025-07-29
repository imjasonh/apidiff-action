const core = require('@actions/core');
const github = require('@actions/github');

const COMMENT_TAG = '<!-- apidiff-action -->';

/**
 * Creates or updates a comment on the PR with API diff results
 * @param {Object} options
 * @param {string} options.token - GitHub token
 * @param {string} options.body - Comment body (markdown)
 * @returns {Promise<void>}
 */
async function createOrUpdateComment({ token, body }) {
  const context = github.context;

  // Only run on pull requests
  if (!context.payload.pull_request) {
    core.info('Not running on a pull request, skipping comment');
    return;
  }

  const octokit = github.getOctokit(token);
  const { owner, repo } = context.repo;
  const prNumber = context.payload.pull_request.number;

  // Add our hidden tag to identify our comments
  const commentBody = `${COMMENT_TAG}\n${body}`;

  try {
    // Find existing comment
    const { data: comments } = await octokit.rest.issues.listComments({
      owner,
      repo,
      issue_number: prNumber,
    });

    const existingComment = comments.find((comment) => comment.body?.includes(COMMENT_TAG));

    if (existingComment) {
      // Update existing comment
      await octokit.rest.issues.updateComment({
        owner,
        repo,
        comment_id: existingComment.id,
        body: commentBody,
      });
      core.info(`Updated existing comment: ${existingComment.html_url}`);
    } else {
      // Create new comment
      const { data: newComment } = await octokit.rest.issues.createComment({
        owner,
        repo,
        issue_number: prNumber,
        body: commentBody,
      });
      core.info(`Created new comment: ${newComment.html_url}`);
    }
  } catch (error) {
    core.warning(`Failed to create/update PR comment: ${error.message}`);
    // Don't fail the action if commenting fails
  }
}

module.exports = { createOrUpdateComment };
