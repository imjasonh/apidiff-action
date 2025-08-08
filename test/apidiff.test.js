const path = require('path');
const exec = require('@actions/exec');

jest.mock('@actions/exec');

const { runApidiff } = require('../lib/apidiff');

describe('runApidiff', () => {
  test('exports all packages when comparing directories', async () => {
    exec.exec.mockImplementation((_cmd, _args, _options) => {
      // Simulate apidiff being installed and all commands succeeding
      return Promise.resolve(0);
    });

    const workingDirectory = '/work';
    const oldRef = 'old';
    const newRef = 'new';

    await runApidiff({ workingDirectory, oldRef, newRef });

    // Also run again with commit-like refs to exercise the git-ref path
    await runApidiff({ workingDirectory, oldRef: 'abc1234', newRef: 'def5678' });

    expect(exec.exec).toHaveBeenCalledWith(
      'apidiff',
      ['-w', 'apidiff.export', './...'],
      expect.objectContaining({ cwd: path.join(workingDirectory, oldRef) })
    );

    expect(exec.exec).toHaveBeenCalledWith(
      'apidiff',
      ['-w', 'apidiff.export', './...'],
      expect.objectContaining({ cwd: path.join(workingDirectory, newRef) })
    );
  });
});
