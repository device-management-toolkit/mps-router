<a name="v2.3.2"></a>
## [v2.3.2] - 2024-01-09
### Build
- **deps:** bump actions/upload-artifact from 3.1.3 to 4.0.0 (#09c5907) 
- **deps:** bump github/codeql-action from 2.22.10 to 3.22.11 (#7c4546e) 
- **deps:** bump github/codeql-action from 2.22.9 to 2.22.10 (#d735ed3) 
- **deps:** bump aquasecurity/trivy-action from 0.15.0 to 0.16.0 (#e186044) 
- **deps:** bump golang from `5c1cabd` to `feceecc` (#cd1251e) 
- **deps:** bump github/codeql-action from 2.22.8 to 2.22.9 (#f48d5af) 
- **deps:** bump actions/setup-go from 4.1.0 to 5.0.0 (#9ca7b66) 
- **deps:** bump go.mongodb.org/mongo-driver from 1.13.0 to 1.13.1 (#ee7f3c9) 
- **deps:** bump golang from `70afe55` to `5c1cabd` (#35f4a80) 
- **deps:** bump aquasecurity/trivy-action from 0.14.0 to 0.15.0 (#d0926b9) 

### Chore
- update build tasks and changelog (#380d241) 

### Ci
- add trivy results to github security tab (#bb967ce) 

<a name="v2.3.1"></a>
## [v2.3.1] - 2023-12-04
### Build
- **deps:** bump golang from `110b07a` to `70afe55` (#89c2cd4) 
- **deps:** bump actions/dependency-review-action from 3.1.3 to 3.1.4 (#116fd5b) 
- **deps:** bump github/codeql-action from 2.22.7 to 2.22.8 (#b0ba160) 
- **deps:** bump step-security/harden-runner from 2.6.0 to 2.6.1 (#e374019) 
- **deps:** bump github/codeql-action from 2.22.6 to 2.22.7 (#6895d39) 
- **deps:** bump actions/dependency-review-action from 3.1.2 to 3.1.3 (#2675357) 
- **deps:** bump github/codeql-action from 2.22.5 to 2.22.6 (#05296b1) 
- **deps:** bump aquasecurity/trivy-action from 0.13.1 to 0.14.0 (#393e491) 
- **deps:** bump go.mongodb.org/mongo-driver from 1.12.1 to 1.13.0 (#3d45cf0) 
- **deps:** bump actions/dependency-review-action from 3.1.1 to 3.1.2 (#51c2018) 
- **deps:** bump golang from `96a8a70` to `110b07a` (#ed3f0bc) 
- **deps:** bump wagoid/commitlint-github-action from 5.4.3 to 5.4.4 (#a772246) 
- **deps:** bump actions/dependency-review-action from 3.1.0 to 3.1.1 (#9247418) 
- **deps:** bump golang from `926f7f7` to `96a8a70` (#74660ae) 
- **deps:** bump aquasecurity/trivy-action from 0.13.0 to 0.13.1 (#63900ce) 
- **deps:** bump github/codeql-action from 2.22.4 to 2.22.5 (#a8da474) 
- **deps:** bump aquasecurity/trivy-action from 0.12.0 to 0.13.0 (#8385eef) 
- **deps:** bump ossf/scorecard-action from 2.3.0 to 2.3.1 (#c6a2830) 
- **deps:** bump github/codeql-action from 2.22.3 to 2.22.4 (#88596d9) 
- **deps:** bump actions/checkout from 4.1.0 to 4.1.1 (#43b1674) 
- **deps:** bump github/codeql-action from 2.22.2 to 2.22.3 (#35b9a21) 
- **deps:** bump github/codeql-action from 2.22.1 to 2.22.2 (#f4e8ed5) 
- **deps:** bump golang from `a76f153` to `926f7f7` (#a541ffd) 
- **deps:** bump github/codeql-action from 2.22.0 to 2.22.1 (#fde5af2) 
- **deps:** bump ossf/scorecard-action from 2.2.0 to 2.3.0 (#882f4fc) 
- **deps:** bump github/codeql-action from 2.21.9 to 2.22.0 (#f870824) 
- **deps:** bump golang from `1c9cc94` to `a76f153` (#57d48c4) 
- **deps:** bump golang from `4bc6541` to `1c9cc94` (#1e9f2ec) 
- **deps:** bump step-security/harden-runner from 2.5.1 to 2.6.0 (#e71babe) 
- **deps:** bump golang from `ec31b7f` to `4bc6541` (#929c948) 
- **deps:** bump github/codeql-action from 2.21.8 to 2.21.9 (#1be6f50) 
- **deps:** bump golang from `96634e5` to `ec31b7f` (#8a1a30b) 
- **deps:** bump actions/checkout from 4.0.0 to 4.1.0 (#e31da3b) 
- **deps:** bump docker/login-action from 2.2.0 to 3.0.0 (#3614f41) 
- **deps:** bump github/codeql-action from 2.21.7 to 2.21.8 (#041537f) 

### Ci
- fix wrong changelog upload (#8206f5e) 
- fixes releaserc formatting and adds persist-credentials: false (#f38669e) 
- release action updates version and changelog (#15ecc43) 
- fix image name for internal docker release (#ee28b1c) 
- add release tag to docker image (#db37184) 

### Feat
- add nosql support via mongo driver (#7155eec) 

### Refactor
- add db interface (#eddd5ad) 

<a name="2.1.2"></a>
## [2.1.2] - 2023-08-02
### Build
- **deps:** bump github.com/stretchr/testify from 1.8.3 to 1.8.4 (#6a3a8cb) 
- **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.3 (#bf93755) 
- **deps:** bump github.com/lib/pq from 1.10.8 to 1.10.9 (#3c1ac31) 

### Ci
- updated semantic release version (#1805ae6) 
- adds codecov, go fmt, go vet, and go lint (#0a3a00e) 

<a name="2.1.1"></a>
## [2.1.1] - 2022-10-05
### Build
- **deps:** bump github.com/lib/pq from 1.10.6 to 1.10.7 (#7e3f3a3) 
- **deps:** bump github.com/stretchr/testify from 1.7.1 to 1.8.0 (#bf8cd4f) 
- **deps:** bump github.com/lib/pq from 1.10.5 to 1.10.6 (#e4d6def) 
- **docker:** remove GOARCH build flag (#f590d82) 

### Ci
- add junit test output (#3233976) 

### Fix
- release updated dependencies (#f3e5830) 

### Test
- **proxy:** refactors parseguid to use table driven (#d59687b) 


<a name="v2.1.0"></a>
## [v2.1.0] - 2022-05-03
### Build
- **deps:** bump github.com/lib/pq from 1.10.4 to 1.10.5 (#5f91e2f) 
- **deps:** bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (#795b20d) 
- **deps:** bump github.com/lib/pq from 1.10.1 to 1.10.4 (#adebc41) 

### Ci
- **jenkinsfile:** removes protex scan (#89c4ecc) 
- **lint:** adds semantic checks to PRs (#a42845c) 
- **release:** adds semantic release to repo (#891195c) 

### Docs
- **github:** add pull request template (#36ff9bf) 

### Feat
- **env:** add option to override default mps host (#0b5fdd9) 
- **healthcheck:** adds flag for checking db status (#b3360c7) 


<a name="v2.0.0"></a>
## [v2.0.0] - 2021-09-15
### Ci
- **changelog:** add automation (#4cc5126) 

### Docs
- **security:** added SECURITY.md file (#777962e) 
- **security:** added security.md file (#3d8be20) 


<a name="v1.4.0"></a>
## v1.4.0 - 2021-06-23
### Build
- **changelog:** add config (#625fb46) 
- **docker:** use non root user (#df60f17) 
- **scan:** fixed protex and checkmarx scan (#2b0313a) 
- **scan:** fixed MPS-Router checkmarx scan (#a9633cc) 
- **scans:** enabled checkmarx (#f34d530) 
- **scans:** enabled Checkmarx (#355035f) 

### Ci
- add jenkinsfile (#5b25763) 
- **changelog:** add automation for changelog generation (#6f6a354) 
- **jenkins:** fix protex project name (#7b3fa88) 

### Docs
- **changelog:** fix version (#e47c569) 
- **changelog:** add changelog (#7d622f1) 
- **copyright:** add missing header (#8d88116) 

### Feat
- **build:** added git workflows (#7ce7e80) 

### Fix
- guid parse now supports v1-4 (#393dc3f) 
- **proxy:** Updated mps server and mps router ports as env variables (#4bb454b) 
- **test:** added unit tests for db (#bd816f6) 

### Test
- **proxy:** test forward and backward functions (#d02b571) 


[Unreleased]: https://github.com/open-amt-cloud-toolkit/mps/compare/2.0.0...HEAD
[2.0.0]: https://github.com/open-amt-cloud-toolkit/mps/compare/v2.1.0...2.0.0
[v2.1.0]: https://github.com/open-amt-cloud-toolkit/mps/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/open-amt-cloud-toolkit/mps/compare/v1.4.0...v2.0.0
