const { parseApidiffOutput, formatChangesAsMarkdown } = require('../lib/parser');

describe('parseApidiffOutput', () => {
  test('parses empty output', () => {
    const result = parseApidiffOutput('');
    expect(result.hasBreakingChanges).toBe(false);
    expect(result.breakingCount).toBe(0);
    expect(result.compatibleCount).toBe(0);
    expect(result.packages).toEqual([]);
  });

  test('parses output with no changes', () => {
    // When there are no changes, apidiff produces no output
    const output = '';

    const result = parseApidiffOutput(output);
    expect(result.hasBreakingChanges).toBe(false);
    expect(result.breakingCount).toBe(0);
    expect(result.compatibleCount).toBe(0);
    expect(result.packages).toEqual([]);
  });

  test('parses output with breaking changes only', () => {
    const output = `Incompatible changes:
- (*Client).DoSomething: changed from func(string) error to func(context.Context, string) error
- MaxRetries: removed`;

    const result = parseApidiffOutput(output);
    expect(result.hasBreakingChanges).toBe(true);
    expect(result.breakingCount).toBe(2);
    expect(result.compatibleCount).toBe(0);
    expect(result.packages).toHaveLength(1);
    expect(result.packages[0].name).toBe('default');
    expect(result.packages[0].breaking).toHaveLength(2);
    expect(result.packages[0].compatible).toHaveLength(0);
  });

  test('parses output with compatible changes only', () => {
    const output = `Compatible changes:
- NewOption: added
- WithTimeout: added
- DefaultTimeout: added`;

    const result = parseApidiffOutput(output);
    expect(result.hasBreakingChanges).toBe(false);
    expect(result.breakingCount).toBe(0);
    expect(result.compatibleCount).toBe(3);
    expect(result.packages).toHaveLength(1);
    expect(result.packages[0].compatible).toHaveLength(3);
  });

  test('parses output with both breaking and compatible changes', () => {
    const output = `Incompatible changes:
- (*Client).DoSomething: changed from func(string) error to func(context.Context, string) error
Compatible changes:
- NewOption: added
- WithTimeout: added`;

    const result = parseApidiffOutput(output);
    expect(result.hasBreakingChanges).toBe(true);
    expect(result.breakingCount).toBe(1);
    expect(result.compatibleCount).toBe(2);
    expect(result.packages).toHaveLength(1);
    expect(result.packages[0].breaking).toHaveLength(1);
    expect(result.packages[0].compatible).toHaveLength(2);
  });

  test('parses real apidiff output', () => {
    const output = `Incompatible changes:
- (*Greeter).Greet: changed from func() string to func(bool) string
- MaxRetries: removed
- Process: changed from func(string) (string, error) to func(context.Context, string) (string, error)
Compatible changes:
- Greeter.Language: added`;

    const result = parseApidiffOutput(output);
    expect(result.hasBreakingChanges).toBe(true);
    expect(result.breakingCount).toBe(3);
    expect(result.compatibleCount).toBe(1);
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
    expect(markdown).toContain('No API changes detected');
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
    expect(markdown).toContain('Breaking changes | 1');
    expect(markdown).toContain('⚠️ **This PR contains breaking API changes!**');
    expect(markdown).toContain('❌ Breaking changes');
    expect(markdown).toContain('- Foo: removed');
    expect(markdown).not.toContain('`default`'); // Should not show default package name
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
    expect(markdown).toContain('Compatible changes | 2');
    expect(markdown).toContain('✅ Compatible changes');
    expect(markdown).toContain('- Bar: added');
    expect(markdown).toContain('- Baz: added');
  });
});
