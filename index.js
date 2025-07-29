const core = require('@actions/core');
const github = require('@actions/github');
const { runApidiff } = require('./lib/apidiff');
const { parseApidiffOutput, formatChangesAsMarkdown } = require('./lib/parser');
const { createOrUpdateComment } = require('./lib/commenter');

async function run() {
  try {
    // Get inputs
    const workingDirectory = core.getInput('working-directory');
    const failOnBreaking = core.getInput('fail-on-breaking') === 'true';
    const commentOnPr = core.getInput('comment-on-pr') === 'true';
    const token = core.getInput('token');

    // Get context
    const context = github.context;

    // Determine refs to compare
    let oldRef, newRef;

    // Check if INPUT_OLD and INPUT_NEW env vars are set (for testing)
    if (process.env.INPUT_OLD && process.env.INPUT_NEW) {
      oldRef = process.env.INPUT_OLD;
      newRef = process.env.INPUT_NEW;
      core.info(`Using provided refs: old=${oldRef}, new=${newRef}`);
    } else if (context.payload.pull_request) {
      // Running on a pull request
      oldRef = context.payload.pull_request.base.sha;
      newRef = context.payload.pull_request.head.sha;
      core.info(`Comparing PR base (${oldRef}) with head (${newRef})`);
    } else if (context.payload.before && context.payload.after) {
      // Running on a push event
      oldRef = context.payload.before;
      newRef = context.payload.after;
      core.info(`Comparing commits ${oldRef} with ${newRef}`);
    } else {
      throw new Error(
        'Unable to determine commits to compare. This action should be run on pull_request or push events.'
      );
    }

    // Run apidiff
    const output = await runApidiff({ workingDirectory, oldRef, newRef });

    // Parse output
    const parsedChanges = parseApidiffOutput(output);

    // Set outputs
    core.setOutput('has-breaking-changes', parsedChanges.hasBreakingChanges.toString());
    core.setOutput('breaking-count', parsedChanges.breakingCount.toString());
    core.setOutput('compatible-count', parsedChanges.compatibleCount.toString());

    // Create comment if requested and on PR
    if (commentOnPr && context.payload.pull_request) {
      const commentBody = formatChangesAsMarkdown(parsedChanges);
      await createOrUpdateComment({ token, body: commentBody });
    }

    // Fail if breaking changes detected and configured to do so
    if (failOnBreaking && parsedChanges.hasBreakingChanges) {
      core.setFailed(`Found ${parsedChanges.breakingCount} breaking API changes`);
    } else if (parsedChanges.hasBreakingChanges) {
      core.warning(`Found ${parsedChanges.breakingCount} breaking API changes`);
    } else {
      core.info('No breaking API changes detected');
    }
  } catch (error) {
    core.setFailed(`Action failed: ${error.message}`);
  }
}

// Run the action
run();
