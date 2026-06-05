// distributions/npm/scripts/install.js
// Downloads the platform-specific native binary of envguard on postinstall.

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const packageJson = require('../package.json');
const version = packageJson.version;

const platformMap = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
};

const archMap = {
  x64: 'amd64',
  arm64: 'arm64',
  ia32: '386'
};

const os = platformMap[process.platform];
const arch = archMap[process.arch];

if (!os || !arch) {
  console.error(`Unsupported platform/architecture: ${process.platform}/${process.arch}`);
  process.exit(1);
}

const ext = os === 'windows' ? 'zip' : 'tar.gz';
const archiveName = `envguard_${os}_${arch}.${ext}`;
const url = `https://github.com/Vamshavardhan50/envguard/releases/download/v${version}/${archiveName}`;

const binDir = path.join(__dirname, '../bin');
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

const archivePath = path.join(binDir, archiveName);

console.log(`Downloading envguard v${version} from ${url}...`);

function download(url, dest, cb) {
  const file = fs.createWriteStream(dest);
  https.get(url, (response) => {
    if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
      return download(response.headers.location, dest, cb);
    }
    if (response.statusCode !== 200) {
      fs.unlink(dest, () => {});
      return cb(new Error(`Failed to download binary: HTTP ${response.statusCode}`));
    }
    response.pipe(file);
    file.on('finish', () => {
      file.close(cb);
    });
  }).on('error', (err) => {
    fs.unlink(dest, () => {});
    cb(err);
  });
}

download(url, archivePath, (err) => {
  if (err) {
    console.error(`Error downloading binary: ${err.message}`);
    process.exit(1);
  }

  console.log('Extracting archive...');
  try {
    if (os === 'windows') {
      execSync(`powershell -Command "Expand-Archive -Path '${archivePath}' -DestinationPath '${binDir}' -Force"`);
    } else {
      execSync(`tar -xzf "${archivePath}" -C "${binDir}"`);
    }
    console.log('Extraction complete!');
    fs.unlinkSync(archivePath);

    if (os !== 'windows') {
      const binaryPath = path.join(binDir, 'envguard');
      fs.chmodSync(binaryPath, 0o755);
    }
  } catch (extractErr) {
    console.error(`Failed to extract binary: ${extractErr.message}`);
    process.exit(1);
  }
});
