const core = require('@actions/core');
const exec = require('@actions/exec');
const path = require('path');

/**
 * Runs the apidiff tool and returns the output
 * @param {Object} options
 * @param {string} options.workingDirectory - Directory to run apidiff in
 * @param {string} options.oldRef - Old git ref (base)
 * @param {string} options.newRef - New git ref (head)
 * @returns {Promise<string>} The apidiff output
 */
async function runApidiff({ workingDirectory, oldRef, newRef }) {
  core.info(`Running apidiff between ${oldRef} and ${newRef}`);

  // First, ensure apidiff is installed
  await installApidiff();

  // Create export files if comparing directories
  const isDirectory = !oldRef.match(/^[0-9a-f]{7,40}$/i);
  
  if (isDirectory) {
    // For directories, we need to create export files first
    const oldDir = path.join(workingDirectory, oldRef);
    const newDir = path.join(workingDirectory, newRef);
    const oldExportPath = path.join(oldDir, 'apidiff.export');
    const newExportPath = path.join(newDir, 'apidiff.export');
    
    // Create export for old version
    await exec.exec('apidiff', ['-w', 'apidiff.export', '.'], {
      cwd: oldDir
    });
    
    // Create export for new version
    await exec.exec('apidiff', ['-w', 'apidiff.export', '.'], {
      cwd: newDir
    });
    
    // Now compare the export files
    let output = '';
    let errorOutput = '';

    const options = {
      cwd: workingDirectory,
      listeners: {
        stdout: (data) => {
          output += data.toString();
        },
        stderr: (data) => {
          errorOutput += data.toString();
        }
      },
      ignoreReturnCode: true
    };

    const exitCode = await exec.exec('apidiff', [oldExportPath, newExportPath], options);
    
    if (errorOutput) {
      core.warning(`apidiff stderr: ${errorOutput}`);
    }

    // apidiff returns exit code 1 if there are incompatible changes
    // This is expected behavior, not an error
    if (exitCode !== 0 && exitCode !== 1) {
      throw new Error(`apidiff failed with exit code ${exitCode}: ${errorOutput}`);
    }

    return output;
  } else {
    // For git refs, use apidiff directly
    let output = '';
    let errorOutput = '';

    const options = {
      cwd: workingDirectory,
      listeners: {
        stdout: (data) => {
          output += data.toString();
        },
        stderr: (data) => {
          errorOutput += data.toString();
        }
      },
      ignoreReturnCode: true
    };

    const exitCode = await exec.exec('apidiff', [oldRef, newRef], options);

    if (errorOutput) {
      core.warning(`apidiff stderr: ${errorOutput}`);
    }

    // apidiff returns exit code 1 if there are incompatible changes
    // This is expected behavior, not an error
    if (exitCode !== 0 && exitCode !== 1) {
      throw new Error(`apidiff failed with exit code ${exitCode}: ${errorOutput}`);
    }

    return output;
  }
}

/**
 * Installs the apidiff tool if not already installed
 */
async function installApidiff() {
  // Check if apidiff is already installed
  try {
    await exec.exec('which', ['apidiff'], { silent: true });
    core.info('apidiff is already installed');
    return;
  } catch {
    // apidiff not found, install it
    core.info('Installing apidiff...');
    await exec.exec('go', ['install', 'golang.org/x/exp/cmd/apidiff@latest']);
    
    // Add Go bin to PATH
    const goPath = process.env.GOPATH || path.join(process.env.HOME, 'go');
    const goBin = path.join(goPath, 'bin');
    core.addPath(goBin);
    
    core.info('apidiff installed successfully');
  }
}

module.exports = { runApidiff };