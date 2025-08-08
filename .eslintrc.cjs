module.exports = {
  root: true,
  ignorePatterns: ['dist/**', 'node_modules/**', 'coverage/**'],
  env: {
    node: true,
    jest: true,
    es2021: true,
  },
  extends: ['eslint:recommended', 'prettier'],
  parserOptions: {
    ecmaVersion: 2021,
    sourceType: 'script',
  },
  rules: {
    'no-console': 'off',
    'no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
  },
};
