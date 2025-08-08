    // Create export for old version. Use './...' to include all packages
    // within the module, not just the root package. Without this, modules
    // that only contain packages in subdirectories (a common layout) would
    // fail with "no Go files".
    await exec.exec('apidiff', ['-w', 'apidiff.export', './...'], {
    await exec.exec('apidiff', ['-w', 'apidiff.export', './...'], {
