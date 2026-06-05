#!/usr/bin/env node
// distributions/npm/bin/envguard.js
// Executable entrypoint that spawns the downloaded native envguard binary.

const path = require('path');
const { spawn } = require('child_process');

const binaryName = process.platform === 'win32' ? 'envguard.exe' : 'envguard';
const binaryPath = path.join(__dirname, binaryName);

const args = process.argv.slice(2);

const proc = spawn(binaryPath, args, { stdio: 'inherit' });

proc.on('close', (code) => {
  process.exit(code);
});

proc.on('error', (err) => {
  console.error(`Failed to execute envguard: ${err.message}`);
  process.exit(2);
});
