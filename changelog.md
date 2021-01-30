### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-v2-tekton-converter/releases/download/v0.0.2/jx-v2-tekton-converter-linux-amd64.tar.gz | tar xzv 
sudo mv jx-v2-tekton-converter /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-v2-tekton-converter/releases/download/v0.0.2/jx-v2-tekton-converter-darwin-amd64.tar.gz | tar xzv
sudo mv jx-v2-tekton-converter /usr/local/bin
```

## Changes

### Bug Fixes

* improve the converter for v2 (jenkins-x-bot-test)
* lets not fail if we have no git clone task (James Strachan)
* add tekton pipelines (James Strachan)
* remove old volume mount (James Strachan)
* initial import (James Strachan)

### Chores

* upgrade pipelines (jenkins-x-bot-test)
* upgrade pipelines (James Strachan)
