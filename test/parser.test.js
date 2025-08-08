const { describe, test } = require('node:test');
const assert = require('node:assert/strict');
const { parseApidiffOutput, formatChangesAsMarkdown } = require('../lib/parser');

describe('parseApidiffOutput', () => {
  test('parses empty output', () => {
    const result = parseApidiffOutput('');
    assert.strictEqual(result.hasBreakingChanges, false);
    assert.strictEqual(result.breakingCount, 0);
    assert.strictEqual(result.compatibleCount, 0);
    assert.deepStrictEqual(result.packages, []);
  });

  test('parses output with no changes', () => {
    // When there are no changes, apidiff produces no output
    const output = '';

    const result = parseApidiffOutput(output);
    assert.strictEqual(result.hasBreakingChanges, false);
    assert.strictEqual(result.breakingCount, 0);
    assert.strictEqual(result.compatibleCount, 0);
    assert.deepStrictEqual(result.packages, []);
  });

  test('parses output with breaking changes only', () => {
    const output = `Incompatible changes:\n- (*Client).DoSomething: changed from func(string) error to func(context.Context, string) error\n- MaxRetries: removed`;

    const result = parseApidiffOutput(output);
    assert.strictEqual(result.hasBreakingChanges, true);
    assert.strictEqual(result.breakingCount, 2);
    assert.strictEqual(result.compatibleCount, 0);
    assert.strictEqual(result.packages.length, 1);
    assert.strictEqual(result.packages[0].name, 'default');
    assert.strictEqual(result.packages[0].breaking.length, 2);
    assert.strictEqual(result.packages[0].compatible.length, 0);
  });

  test('parses output with compatible changes only', () => {
    const output = `Compatible changes:\n- NewOption: added\n- WithTimeout: added\n- DefaultTimeout: added`;

    const result = parseApidiffOutput(output);
    assert.strictEqual(result.hasBreakingChanges, false);
    assert.strictEqual(result.breakingCount, 0);
    assert.strictEqual(result.compatibleCount, 3);
    assert.strictEqual(result.packages.length, 1);
    assert.strictEqual(result.packages[0].compatible.length, 3);
  });

  test('parses output with both breaking and compatible changes', () => {
    const output = `Incompatible changes:\n- (*Client).DoSomething: changed from func(string) error to func(context.Context, string) error\nCompatible changes:\n- NewOption: added\n- WithTimeout: added`;

    const result = parseApidiffOutput(output);
    assert.strictEqual(result.hasBreakingChanges, true);
    assert.strictEqual(result.breakingCount, 1);
    assert.strictEqual(result.compatibleCount, 2);
    assert.strictEqual(result.packages.length, 1);
    assert.strictEqual(result.packages[0].breaking.length, 1);
    assert.strictEqual(result.packages[0].compatible.length, 2);
  });

  test('parses real apidiff output', () => {
    const output = `Incompatible changes:\n- (*Greeter).Greet: changed from func() string to func(bool) string\n- MaxRetries: removed\n- Process: changed from func(string) (string, error) to func(context.Context, string) (string, error)\nCompatible changes:\n- Greeter.Language: added`;

    const result = parseApidiffOutput(output);
    assert.strictEqual(result.hasBreakingChanges, true);
    assert.strictEqual(result.breakingCount, 3);
    assert.strictEqual(result.compatibleCount, 1);
  });
});

describe('formatChangesAsMarkdown', () => {
  test('formats no changes', () => {
    const parsedChanges = {
      hasBreakingChanges: false,
      breakingCount: 0,
      compatibleCount: 0,
      packages: [],
    };

    const markdown = formatChangesAsMarkdown(parsedChanges);
    assert.ok(markdown.includes('No API changes detected'));
  });

  test('formats breaking changes', () => {
    const parsedChanges = {
      hasBreakingChanges: true,
      breakingCount: 1,
      compatibleCount: 0,
      packages: [
        {
          name: 'default',
          breaking: [{ message: 'Foo: removed', compatible: false }],
          compatible: [],
        },
      ],
    };

    const markdown = formatChangesAsMarkdown(parsedChanges);
    assert.ok(markdown.includes('Breaking changes | 1'));
    assert.ok(markdown.includes('⚠️ **This PR contains breaking API changes!**'));
    assert.ok(markdown.includes('❌ Breaking changes'));
    assert.ok(markdown.includes('- Foo: removed'));
    assert.ok(!markdown.includes('`default`'));
  });

  test('formats compatible changes', () => {
    const parsedChanges = {
      hasBreakingChanges: false,
      breakingCount: 0,
      compatibleCount: 2,
      packages: [
        {
          name: 'default',
          breaking: [],
          compatible: [
            { message: 'Bar: added', compatible: true },
            { message: 'Baz: added', compatible: true },
          ],
        },
      ],
    };

    const markdown = formatChangesAsMarkdown(parsedChanges);
    assert.ok(markdown.includes('Compatible changes | 2'));
    assert.ok(markdown.includes('✅ Compatible changes'));
    assert.ok(markdown.includes('- Bar: added'));
    assert.ok(markdown.includes('- Baz: added'));
  });
});
