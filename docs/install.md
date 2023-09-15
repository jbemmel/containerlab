---
hide:
  - navigation
---
Containerlab is distributed as a Linux deb/rpm package and can be installed on any Debian- or RHEL-like distributive in a matter of a few seconds.

### Pre-requisites

The following requirements must be satisfied to let containerlab tool run successfully:

* A user should have `sudo` privileges to run containerlab.
* A Linux server/VM[^2] and [Docker](https://docs.docker.com/engine/install/) installed.
* Load container images (e.g. Nokia SR Linux, Arista cEOS) that are not downloadable from a container registry. Containerlab will try to pull images at runtime if they do not exist locally.

### Install script

Containerlab can be installed using the [installation script](https://github.com/srl-labs/containerlab/blob/main/get.sh) which detects the operating system type and installs the relevant package:

!!! note
    Containerlab is distributed via deb/rpm packages, thus only Debian- and RHEL-like distributives can leverage package installation.  
    Other systems can follow the [manual installation](#manual-installation) procedure.

```bash
# download and install the latest release (may require sudo)
bash -c "$(curl -sL https://get.containerlab.dev)"

# download a specific version - 0.10.3 (may require sudo)
bash -c "$(curl -sL https://get.containerlab.dev)" -- -v 0.10.3

# with wget
bash -c "$(wget -qO - https://get.containerlab.dev)"
```

### Package managers

It is possible to install official containerlab releases via public APT/YUM repository.

=== "APT"
    ```bash
    echo "deb [trusted=yes] <https://apt.fury.io/netdevops/> /" | \
    sudo tee -a /etc/apt/sources.list.d/netdevops.list

    sudo apt update && sudo apt install containerlab
    ```
=== "YUM"
    ```
    yum-config-manager --add-repo=<https://yum.fury.io/netdevops/> && \
    echo "gpgcheck=0" | sudo tee -a /etc/yum.repos.d/yum.fury.io_netdevops_.repo

    sudo yum install containerlab
    ```
=== "APK"
    Download `.apk` package from [Github releases](https://github.com/srl-labs/containerlab/releases).
=== "AUR"
    Arch Linux users can download a package from this [AUR repository](https://aur.archlinux.org/packages/containerlab-bin).

??? "Manual package installation"
    Alternatively, users can manually download the deb/rpm package from the [Github releases](https://github.com/srl-labs/containerlab/releases) page.

    example:
    ```bash
    # manually install latest release with package managers
    LATEST=$(curl -s https://github.com/srl-labs/containerlab/releases/latest | sed -e 's/.*tag\/v\(.*\)\".*/\1/')
    # with yum
    yum install "https://github.com/srl-labs/containerlab/releases/download/v${LATEST}/containerlab_${LATEST}_linux_amd64.rpm"
    # with dpkg
    curl -sL -o /tmp/clab.deb "https://github.com/srl-labs/containerlab/releases/download/v${LATEST}/containerlab_${LATEST}_linux_amd64.deb" && dpkg -i /tmp/clab.deb

    # install specific release with yum
    yum install https://github.com/srl-labs/containerlab/releases/download/v0.7.0/containerlab_0.7.0_linux_386.rpm
    ```

The package installer will put the `containerlab` binary in the `/usr/bin` directory as well as create the `/usr/bin/clab -> /usr/bin/containerlab` symlink. The symlink allows the users to save on typing when they use containerlab: `clab <command>`.

### Container

Containerlab is also available in a container packaging. The latest containerlab release can be pulled with:

```
docker pull ghcr.io/srl-labs/clab
```

To pick any of the released versions starting from release 0.19.0, use the version number as a tag, for example, `docker pull ghcr.io/srl-labs/clab:0.19.0`

Since containerlab itself deploys containers and creates veth pairs, its run instructions are a bit more complex, but still, it is a copy-paste-able command.

For example, if your lab files are contained within the current working directory - `$(pwd)` - then you can launch containerlab container as follows:

```bash
docker run --rm -it --privileged \
    --network host \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/run/netns:/var/run/netns \
    -v /etc/hosts:/etc/hosts \
    -v /var/lib/docker/containers:/var/lib/docker/containers \
    --pid="host" \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    ghcr.io/srl-labs/clab bash
```

Within the started container you can use the same `containerlab deploy/destroy/inspect` commands to manage your labs.

!!!note
    Containerlab' container command is itself `containerlab`, so you can deploy a lab without invoking a shell, for example:
    ```bash
    docker run --rm -it --privileged \
    # <run options omitted>
    -w $(pwd) \
    ghcr.io/srl-labs/clab deploy -t somelab.clab.yml
    ```

### Manual installation

If the linux distributive can't install deb/rpm packages, containerlab can be installed from the archive:

```bash
# get the latest available tag
LATEST=$(curl -s https://github.com/srl-labs/containerlab/releases/latest | \
       sed -e 's/.*tag\/v\(.*\)\".*/\1/')

# download tar.gz archive
curl -L -o /tmp/clab.tar.gz "https://github.com/srl-labs/containerlab/releases/download/v${LATEST}/containerlab_${LATEST}_Linux_amd64.tar.gz"

# create containerlab directory
mkdir -p /etc/containerlab

# extract downloaded archive into the containerlab directory
tar -zxvf /tmp/clab.tar.gz -C /etc/containerlab

# (optional) move containerlab binary somewhere in the $PATH
mv /etc/containerlab/containerlab /usr/bin && chmod a+x /usr/bin/containerlab
```

### Windows Subsystem Linux (WSL)

Containerlab [runs](https://twitter.com/ntdvps/status/1380915270328401922) on WSL, but you need to [install docker-ce](https://docs.docker.com/engine/install/) inside the WSL2 linux system instead of using Docker Desktop[^3].

If you are running Ubuntu 20.04 as your WSL2 machine, you can run [this script](https://gist.github.com/hellt/e8095c1719a3ea0051165ff282d2b62a) to install docker-ce.

```bash
curl -L https://gist.githubusercontent.com/hellt/e8095c1719a3ea0051165ff282d2b62a/raw/1dffb71d0495bb2be953c489cd06a25656d974a4/docker-install.sh | \
bash
```

Once installed, issue `sudo service docker start` to start the docker service inside WSL2 machine.

??? "Running VM-based routers inside WSL"
    At the moment of this writing, KVM support was not available out-of-the box with WSL2 VMs. There are [ways](https://www.reddit.com/r/bashonubuntuonwindows/comments/ldbyxa/what_is_the_current_state_of_kvm_acceleration_on/) to enable KVM support, but they were not tested with containerlab. This means that running traditional VM based routers via [vrnetlab integration](manual/vrnetlab.md) is not readily available.

    It appears to be that next versions of WSL2 kernels will support KVM.

### Mac OS

Running containerlab on Mac OS is possible[^4] by means of a separate docker image with containerlab inside.

!!!warning
    ARM-based Macs (M1/2) are not supported, and no binaries are generated for this platform. This is mainly due to the lack of network images built for arm64 architecture as of now. Nevertheless, it is technically possible to run Containerlab on ARM-based MacBook. Check [How to start Containerlab on ARM-based MacBook](#containerlab-on-arm-based-macs) for details.

To use this container use the following command:

```shell linenums="1"
CLAB_WORKDIR=~/clab

docker run --rm -it --privileged \
    --network host \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /run/netns:/run/netns \
    --pid="host" \
    -w $CLAB_WORKDIR \
    -v $CLAB_WORKDIR:$CLAB_WORKDIR \
    ghcr.io/srl-labs/clab bash
```

The first command in the snippet above sets the working directory which you intend to use on your Mac OS. The `~/clab` in the example above expands to `/Users/<username>/clab` and means that we intent to have our containerlab labs to be stored in this directory.

!!!note
    1. It is best to create a directory under the `~/some/path` unless you know what to do[^5]
    2. vrnetlab based nodes will not be able to start, since Docker VM does not support virtualization.
    3. Docker Desktop for Mac introduced cgroups v2 support in 4.3.0 version; to support the images that require cgroups v1 follow [these instructions](https://github.com/docker/for-mac/issues/6073).
    4. Docker Desktop relies on a LinuxKit based HyperKit VM. Unfortunately, it is shipped with a minimalist kernel, and some modules such as VRF are disabled by default. Follow [these instructions](https://medium.com/@notsinge/making-your-own-linuxkit-with-docker-for-mac-5c1234170fb1) to rebuild it with more modules.

When the container is started, you will have a bash shell opened with the directory contents mounted from the Mac OS. There you can use `containerlab` commands right away.

???tip "Step by step example"
    Let's imagine I want to run a lab with two SR Linux containers running directly on a Mac OS.

    First, I need to have Docker Desktop for Mac installed and running.

    Then I will create a directory under the `$HOME` path on my mac:

    ```
    mkdir -p ~/clab
    ```

    Then I will create a clab file defining my lab in the newly created directory:

    ```bash
    cat <<EOF > ~/clab/2srl.clab.yml
    name: 2srl

    topology:
      nodes:
        srl1:
          kind: srl
          image: ghcr.io/nokia/srlinux
        srl2:
          kind: srl
          image: ghcr.io/nokia/srlinux

      links:
        - endpoints: ["srl1:e1-1", "srl2:e1-1"]
    EOF
    ```

    Now when the clab file is there, launch the container and don't forget to use path to the directory you created:

    ```bash
    CLAB_WORKDIR=~/clab

    docker run --rm -it --privileged \
        --network host \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v /run/netns:/run/netns \
        --pid="host" \
        -w $CLAB_WORKDIR \
        -v $CLAB_WORKDIR:$CLAB_WORKDIR \
        ghcr.io/srl-labs/clab bash
    ```

    Immediately you will get into the directory inside the container with your lab file available:

    ```
    root@docker-desktop:/Users/romandodin/clab# ls
    2srl.clab.yml
    ```

    Now you can launch the lab, as containerlab is already part of the image:
    ```
    root@docker-desktop:/Users/romandodin/clab# clab dep -t 2srl.clab.yml
    INFO[0000] Parsing & checking topology file: 2srl.clab.yml 
    INFO[0000] Creating lab directory: /Users/romandodin/clab/clab-2srl 
    INFO[0000] Creating root CA                             
    INFO[0000] Creating docker network: Name='clab', IPv4Subnet='172.20.20.0/24', IPv6Subnet='2001:172:20:20::/64', MTU='1500' 
    INFO[0000] Creating container: srl1                     
    INFO[0000] Creating container: srl2                     
    INFO[0001] Creating virtual wire: srl1:e1-1 <--> srl2:e1-1 
    INFO[0001] Adding containerlab host entries to /etc/hosts file 
    +---+----------------+--------------+-----------------------+------+-------+---------+----------------+----------------------+
    | # |      Name      | Container ID |         Image         | Kind | Group |  State  |  IPv4 Address  |     IPv6 Address     |
    +---+----------------+--------------+-----------------------+------+-------+---------+----------------+----------------------+
    | 1 | clab-2srl-srl1 | 574bf836fb40 | ghcr.io/nokia/srlinux | srl  |       | running | 172.20.20.2/24 | 2001:172:20:20::2/64 |
    | 2 | clab-2srl-srl2 | f88531a74ffb | ghcr.io/nokia/srlinux | srl  |       | running | 172.20.20.3/24 | 2001:172:20:20::3/64 |
    +---+----------------+--------------+-----------------------+------+-------+---------+----------------+----------------------+
    ```

#### Containerlab on ARM-based Macs

The easiest option to run Containerlab on ARM-based M1/M2 Macs is to build a docker-in-docker container. We'll provide an example of a custom [devcontainer](https://code.visualstudio.com/docs/devcontainers/containers) that can be opened in [VSCode](https://code.visualstudio.com) with [Remote Development extension pack](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack) installed.

Create `.devcontainer` directory in the root of the Containerlab repository with the following content:

```text
.devcontainer
|- devcontainer.json
|- Dockerfile
```

=== "Dockerfile"

    ```Dockerfile
    # The devcontainer will be based on Python 3.9
    # The base container already has entrypoint, vscode user account, etc. out of the box
    FROM mcr.microsoft.com/devcontainers/python:0-3.9-bullseye

    # containelab version will be set in devcontainer.json
    ARG _CLAB_VERSION

    # install some basic tools inside the container
    # adjust this list based on your demands
    RUN apt-get update \
        && apt-get install -y --no-install-recommends \
        sshpass \
        curl \
        iputils-ping \
        htop \
        && rm -rf /var/lib/apt/lists/* \
        && rm -Rf /usr/share/doc && rm -Rf /usr/share/man \
        && apt-get clean

    # install preferred version of the containerlab
    RUN bash -c "$(curl -sL https://get.containerlab.dev)" -- -v ${_CLAB_VERSION} \
        && pip3 install --user yamllint
    ```
=== "devcontainer.json"

    ```json
    // For format details, see https://aka.ms/devcontainer.json. For config options, see the
    // README at: https://github.com/devcontainers/templates/tree/main/src/python
    {
        "name": "clab-for-arm",
        "build": {
            "dockerfile": "Dockerfile",
            "args": {
                "_CLAB_VERSION": "0.43.0"
            }
        },
        "features": {
            // Containerlab will run in a docker-in-docker container
            // it is also possible to use docker-outside-docker feature
            "ghcr.io/devcontainers/features/docker-in-docker:2.2.0": {
                "version": "latest"
            }
        },
        // add any required extensions that must be pre-installed in the devcontainer
        "customizations": {
            "vscode": {
                "extensions": [
                    // various tools
                    "tuxtina.json2yaml",
                    "vscode-icons-team.vscode-icons",
                    "mutantdino.resourcemonitor"
                ]
            }
        }
    }
    ```

Once the devcontainer is defined as described above:

* Open the devcontainer in VSCode
* Import the required images for your cLab inside the container (if you are using Docker-in-Docker option)
* Start you Containerlab

### Upgrade

To upgrade `containerlab` to the latest available version issue the following command[^1]:

```
containerlab version upgrade
```

This command will fetch the installation script and will upgrade the tool to its most recent version.

or leverage `apt`/`yum` utilities if containerlab repo was added as explained in the [Package managers](#package-managers) section.

### From source

To build containerlab from source:

=== "with go build"
    To build containerlab from source, clone the repository and issue `go build` at its root.
=== "with goreleaser"
    When we release containerlab we use [goreleaser](https://goreleaser.com/) project to build binaries for all supported platforms as well as the deb/rpm packages.  
    Users can install `goreleaser` and do the same locally by issuing the following command:
    ```
    goreleaser --snapshot --skip-publish --rm-dist
    ```

### Uninstall

To uninstall containerlab when it was installed via installation script or packages:

=== "Debian-based system"
    ```
    apt remove containerlab
    ```
=== "RPM-based systems"
    ```
    yum remove containerlab
    ```
=== "Manual removal"
    Containerlab binary is located at `/usr/bin/containerlab`. In addition to the binary, containerlab directory with static files may be found at `/etc/containerlab`.

### SELinux

When SELinux set to enforced mode containerlab binary might fail to execute with `Segmentation fault (core dumped)` error. This might be because containerlab binary is compressed with [upx](https://upx.github.io/) and selinux prevents it from being decompressed by default.

To fix this:

```
sudo semanage fcontext -a -t textrel_shlib_t $(which containerlab)
sudo restorecon $(which containerlab)
```

or more globally:

```
sudo setsebool -P selinuxuser_execmod 1
```

[^1]: only available if installed from packages
[^2]: Most containerized NOS will require >1 vCPU. RAM size depends on the lab size. Architecture: AMD64.
[^3]: No need to uninstall Docker Desktop, just make sure that it is not integrated with WSL2 machine that you intend to use with containerlab. Moreover, you can make it even work with Docker Desktop with a [few additional steps](https://twitter.com/networkop1/status/1380976461641834500/photo/1), but installing docker-ce into the WSL maybe more intuitive.
[^4]: kudos to Michael Kashin who [shared](https://github.com/srl-labs/containerlab/issues/577#issuecomment-895847387) this approach with us
[^5]: otherwise make sure to add a custom shared directory to the docker on mac.
