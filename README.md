# SuperK
I'm new to [Golang](https://golang.org/), so I decided to learn the language by creating this experimental tool to run [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) commands in a console user interface.

*This is still a Work In Progress.*

# How to use
## Set up the dev environment
This project leverages VS Code DevContainers which allows us to use a Docker container as a full-featured development environment. Follow the instructions in [Remote development in Containers](https://code.visualstudio.com/remote-tutorials/containers/getting-started) to install Docker and the Remote-Containers extension for VS Code.

Then ```git clone``` this repo and open it in VS Code. VS Code should realize you have a *.devContainer* definition as part of the project and offer to reopen it in a container. You may also use Ctrl+Shift+P and type *"Remote-Containers: Reopen in Container"* instead.

Note that the first time you open the project in a container, it will take a few minutes to build that container.

## Setup K8s cluster
You may create a local [KIND](https://kind.sigs.k8s.io/) cluster for testing purposes by executing ```make kind-create```. KIND allows us to run local K8s clusters using Docker container "nodes".

If you want to target your own K8s cluster (e.g. [AKS](https://azure.microsoft.com/es-es/services/kubernetes-service/)), ```unset KUBECONFIG``` first and then ensure ```kubectx``` points to the appropriate cluster.

## Try the tool
You may build the tool by executing ```make build``` and then run it with ```./superk```, or just execute ```make run``` to do it all in one step.

## Debug the tool
- To debug the tool execute ```make debug``` to start a debug server and then launch VS Code with *"Connect to server"* configuration (or just press F5).
- To run all tests execute ```make test```.
- To run or debug a specific test, open the correspondent **_test.go* file in VS Code and click *"run test"* or *"debug test"* on top of the test function.

## Other useful links:
- GOCUI Library: [Go Console User Interface](https://github.com/jroimartin/gocui).
- GOCUI docs: [package gocui docs](https://pkg.go.dev/github.com/jroimartin/gocui?tab=doc).
- A good example that also uses GOCUI Library: [azbrowse](https://github.com/lawrencegripper/azbrowse).
- Termbox Library used by GOCUI behind the scenes: [nsf/termbox-go](https://github.com/nsf/termbox-go)


