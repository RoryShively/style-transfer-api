# Setup

## Dependencies

The minimum dependencies for this project are minikube, kubeclt and skaffold. Instructions will be provided below of how to download these tools for macOS using homebrew and links will be provided for download instructions on other platforms.

### Kubectl

https://kubernetes.io/docs/tasks/tools/install-kubectl/

**Homebrew**
```
# Install the kubectl cli to control kubernetes cluster
$ brew install kubectl
```


### Minikube

https://kubernetes.io/docs/tasks/tools/install-minikube/

**Homebrew**
```
# Tap cask repository
$ brew tap caskroom/cask
 
# Install virtualbox so minikube can create a VM for the one node kubernetes cluster to run in
$ brew cask install virtualbox
 
# Install minikube to be able to set up kubernetes cluster locally
$ brew cask install minikube
```

### Skaffold

https://skaffold.dev/docs/getting-started/#installing-skaffold

**Homebrew**
```
# Install skaffold
$ brew install skaffold
```


