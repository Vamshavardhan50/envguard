#!/usr/bin/env bash
# scripts/bump-version.sh
# Bumps project version across distributions and creates git tag.

set -euo pipefail

cd "$(dirname "$0")/.."

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <version> (e.g. 1.0.0)" >&2
    exit 1
fi

VERSION=$1
VERSION="${VERSION#v}"

if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
    echo "Error: Version must be in semver format x.y.z (e.g. 1.0.0)" >&2
    exit 1
fi

TAG="v$VERSION"

if ! git diff-index --quiet HEAD --; then
    echo "Error: Git working directory is not clean. Please commit or stash changes first." >&2
    exit 1
fi

BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH" != "main" ] && [ "$BRANCH" != "master" ]; then
    echo "Warning: You are on branch '$BRANCH'. Releases are typically bumped on 'main' or 'master'."
    read -p "Do you want to continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "Bumping version to $TAG..."

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' -E 's/"version": "[^"]+"/"version": "'"$VERSION"'"/' distributions/npm/package.json
    sed -i '' -E 's/version = "[^"]+"/version = "'"$VERSION"'"/' distributions/pip/pyproject.toml
    sed -i '' -E 's/VERSION = "[^"]+"/VERSION = "'"$VERSION"'"/' distributions/pip/envguard/main.py
else
    sed -i -E 's/"version": "[^"]+"/"version": "'"$VERSION"'"/' distributions/npm/package.json
    sed -i -E 's/version = "[^"]+"/version = "'"$VERSION"'"/' distributions/pip/pyproject.toml
    sed -i -E 's/VERSION = "[^"]+"/VERSION = "'"$VERSION"'"/' distributions/pip/envguard/main.py
fi

git add distributions/npm/package.json distributions/pip/pyproject.toml distributions/pip/envguard/main.py
git commit -m "chore: bump version to $VERSION"

git tag -a "$TAG" -m "Release $TAG"

echo "Successfully bumped version and created tag $TAG!"
echo "To publish, run:"
echo "  git push origin HEAD --tags"
