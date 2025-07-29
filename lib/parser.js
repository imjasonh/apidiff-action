const core = require('@actions/core');

/**
 * Parses apidiff output to extract changes
 * @param {string} output - Raw apidiff output
 * @returns {Object} Parsed changes with breaking and compatible counts
 */
function parseApidiffOutput(output) {
  const result = {
    hasBreakingChanges: false,
    breakingCount: 0,
    compatibleCount: 0,
    packages: []
  };

  if (!output || output.trim() === '') {
    return result;
  }

  // Split into lines
  const lines = output.split('\n');
  
  let currentPackage = {
    name: 'default',
    breaking: [],
    compatible: []
  };
  result.packages.push(currentPackage);
  
  let inCompatibleSection = false;
  let inIncompatibleSection = false;

  for (const line of lines) {
    const trimmedLine = line.trim();

    // Skip empty lines
    if (!trimmedLine) {
      continue;
    }

    // Check for section headers
    if (trimmedLine === 'Compatible changes:') {
      inCompatibleSection = true;
      inIncompatibleSection = false;
      continue;
    }

    if (trimmedLine === 'Incompatible changes:') {
      inIncompatibleSection = true;
      inCompatibleSection = false;
      result.hasBreakingChanges = true;
      continue;
    }

    // Parse changes (they start with "- ")
    if ((inCompatibleSection || inIncompatibleSection) && trimmedLine.startsWith('- ')) {
      const changeMessage = trimmedLine.substring(2).trim();
      
      const change = {
        message: changeMessage,
        compatible: inCompatibleSection
      };

      if (inCompatibleSection) {
        currentPackage.compatible.push(change);
        result.compatibleCount++;
      } else if (inIncompatibleSection) {
        currentPackage.breaking.push(change);
        result.breakingCount++;
      }
    }
  }

  // Clean up empty packages
  result.packages = result.packages.filter(pkg => 
    pkg.breaking.length > 0 || pkg.compatible.length > 0
  );

  // If no packages with changes, return empty packages array
  if (result.packages.length === 0) {
    result.packages = [];
  }

  core.info(`Parsed ${result.breakingCount} breaking changes and ${result.compatibleCount} compatible changes`);

  return result;
}

/**
 * Formats changes for display
 * @param {Object} parsedChanges - Parsed changes from parseApidiffOutput
 * @returns {string} Formatted markdown string
 */
function formatChangesAsMarkdown(parsedChanges) {
  let markdown = '# API Compatibility Check Results\n\n';

  if (parsedChanges.packages.length === 0) {
    markdown += '✅ **No API changes detected**\n';
    return markdown;
  }

  // Summary
  markdown += '## Summary\n\n';
  markdown += '| Type | Count |\n';
  markdown += '|------|-------|\n';
  markdown += `| Breaking changes | ${parsedChanges.breakingCount} |\n`;
  markdown += `| Compatible changes | ${parsedChanges.compatibleCount} |\n\n`;

  if (parsedChanges.hasBreakingChanges) {
    markdown += '⚠️ **This PR contains breaking API changes!**\n\n';
  }

  // Details by package
  markdown += '## Details\n\n';
  
  for (const pkg of parsedChanges.packages) {
    if (pkg.name !== 'default') {
      markdown += `### \`${pkg.name}\`\n\n`;
    }

    if (pkg.breaking.length > 0) {
      markdown += '#### ❌ Breaking changes\n\n';
      for (const change of pkg.breaking) {
        markdown += `- ${change.message}\n`;
      }
      markdown += '\n';
    }

    if (pkg.compatible.length > 0) {
      markdown += '#### ✅ Compatible changes\n\n';
      for (const change of pkg.compatible) {
        markdown += `- ${change.message}\n`;
      }
      markdown += '\n';
    }
  }

  return markdown;
}

module.exports = { parseApidiffOutput, formatChangesAsMarkdown };