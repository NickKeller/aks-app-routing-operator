
# Release

Releasing a new version of the operator has been automated. You can see the release workflow in [release.yaml](../.github/workflows/release.yaml)

## Steps

A prerequisite to creating a release is updating the [CHANGELOG.md](../CHANGELOG.md) with the changes that have been made in the release. This should either be done by a PM or reviewed by a PM. PMs must be involved in this process.

After the CHANGELOG has been updated, you can start a release by going to the `Actions` tab and selecting `Release` on the left. Then click `Run workflow` and input the required parameters. It's very important that the SHA used is one that matches the changes detailed in the CHANGELOG exactly.

You can ensure the release workflow ran successfully by watching the workflow then verifying that both the image push and GitHub release were successful. 

## Hotfix

In the unlikely event that a hotfix is needed, you can create a hotfix release through the same steps detailed above. The semantic version should be bumped at the minor version level for the hotfix. For example, a hotfix for `1.0.0` would be released as `1.0.1`. You can note that this is a hotfix in the `CHANGELOG.md`.
