# SDK release guide

This document explains how to publish RLAAS SDKs from release tags.

## Release strategy

- Source tag format: `vX.Y.Z` (example: `v1.1.1`)
- SDK package versions: `X.Y.Z` (without leading `v`)
- Publish all SDKs from the same tagged commit to keep versions aligned.

## Automated (recommended)

Use manual workflow:

- `.github/workflows/sdk-release.yml`

Access lock:

- The workflow includes a runtime authorization gate that allows only users listed in `.github/CODEOWNERS`.
- Optional: configure a protected GitHub Environment with required reviewers for an additional approval gate.

Inputs:

- `tag` (required): release tag like `v1.1.1`
- `publish_python`
- `publish_typescript`
- `publish_java`
- `publish_dotnet`
- `dry_run`

### Required repository secrets

- `PYPI_API_TOKEN` for Python publish
- `NPM_TOKEN` for npm publish
- `MAVEN_USERNAME`, `MAVEN_PASSWORD` for GitHub Packages Maven publish (optional override; defaults can use `github.token`)
- `NUGET_API_KEY` for NuGet publish

## Manual commands

### Python (`sdk/python`)

```bash
python -m pip install --upgrade pip build twine
python -m build
twine check dist/*
twine upload dist/*
```

### TypeScript (`sdk/typescript`)

```bash
npm ci
npm version 1.1.1 --no-git-tag-version
npm run build
npm publish --access public
```

### Java (`sdk/java`)

```bash
mvn -DskipTests versions:set -DnewVersion=1.1.1 -DgenerateBackupPoms=false
mvn -DskipTests package
mvn -DskipTests deploy -DaltDeploymentRepository=github::default::https://maven.pkg.github.com/rlaas-io/rlaas
```

### .NET (`sdk/dotnet/Rlaas.Sdk`)

```bash
dotnet restore
dotnet pack -c Release -p:PackageVersion=1.1.1 -p:Version=1.1.1
dotnet nuget push bin/Release/*.nupkg --source https://api.nuget.org/v3/index.json --api-key <NUGET_API_KEY>
```

## Validation after publish

Run the manual validation workflow in `rlaas-testing` against the same tag with SDK smoke enabled.

## Notes

- The Java workflow publishes to GitHub Packages Maven by default.
- If you need Maven Central release, add signing + OSSRH staging configuration and credentials.
