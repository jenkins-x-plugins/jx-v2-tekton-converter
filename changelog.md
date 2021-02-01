### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-v2-tekton-converter/releases/download/v0.0.4/jx-v2-tekton-converter-linux-amd64.tar.gz | tar xzv 
sudo mv jx-v2-tekton-converter /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-v2-tekton-converter/releases/download/v0.0.4/jx-v2-tekton-converter-darwin-amd64.tar.gz | tar xzv
sudo mv jx-v2-tekton-converter /usr/local/bin
```

## Changes

### Bug Fixes

* lets use the new pipelines if the jenkins-x.yml is empty (James Strachan)

### Chores

* add batch mode (James Strachan)
