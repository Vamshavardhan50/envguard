# distributions/pip/envguard/main.py
# Python wrapper for envguard that downloads the native binary on first run and executes it.

import os
import sys
import platform
import urllib.request
import tarfile
import zipfile
import subprocess

VERSION = "1.0.2"
OWNER = "Vamshavardhan50"
REPO = "envguard"

def get_platform_details():
    system = platform.system().lower()
    machine = platform.machine().lower()

    if system == "darwin":
        os_name = "darwin"
    elif system == "linux":
        os_name = "linux"
    elif system == "windows":
        os_name = "windows"
    else:
        raise OSError(f"Unsupported operating system: {system}")

    if machine in ["amd64", "x86_64"]:
        arch_name = "amd64"
    elif machine in ["arm64", "aarch64"]:
        arch_name = "arm64"
    elif machine in ["386", "i386", "i686"]:
        arch_name = "386"
    else:
        raise OSError(f"Unsupported architecture: {machine}")

    # goreleaser ignores windows arm64 and darwin 386
    if os_name == "windows" and arch_name == "arm64":
        raise OSError("envguard binary is not built for Windows ARM64")
    if os_name == "darwin" and arch_name == "386":
        raise OSError("envguard binary is not built for macOS 386")

    return os_name, arch_name

def download_binary(dest_dir, os_name, arch_name):
    ext = "zip" if os_name == "windows" else "tar.gz"
    archive_name = f"envguard_{os_name}_{arch_name}.{ext}"
    url = f"https://github.com/{OWNER}/{REPO}/releases/download/v{VERSION}/{archive_name}"
    archive_path = os.path.join(dest_dir, archive_name)

    print(f"Downloading envguard v{VERSION} from {url}...")
    try:
        urllib.request.urlretrieve(url, archive_path)
    except Exception as e:
        print(f"Failed to download envguard release: {e}", file=sys.stderr)
        raise

    print("Extracting archive...")
    try:
        if ext == "zip":
            with zipfile.ZipFile(archive_path, 'r') as zip_ref:
                zip_ref.extractall(dest_dir)
        else:
            with tarfile.open(archive_path, 'r:gz') as tar_ref:
                tar_ref.extractall(dest_dir)
    except Exception as e:
        print(f"Failed to extract envguard archive: {e}", file=sys.stderr)
        raise
    finally:
        if os.path.exists(archive_path):
            os.remove(archive_path)

    binary_name = "envguard.exe" if os_name == "windows" else "envguard"
    binary_path = os.path.join(dest_dir, binary_name)
    if os_name != "windows" and os.path.exists(binary_path):
        os.chmod(binary_path, 0o755)

    print("envguard installation complete.")

def run():
    current_dir = os.path.dirname(os.path.abspath(__file__))
    os_name, arch_name = get_platform_details()
    binary_name = "envguard.exe" if os_name == "windows" else "envguard"
    binary_path = os.path.join(current_dir, binary_name)

    if not os.path.exists(binary_path):
        download_binary(current_dir, os_name, arch_name)

    args = sys.argv[1:]
    try:
        res = subprocess.run([binary_path] + args)
        sys.exit(res.returncode)
    except Exception as e:
        print(f"Failed to execute envguard: {e}", file=sys.stderr)
        sys.exit(2)

if __name__ == "__main__":
    run()
