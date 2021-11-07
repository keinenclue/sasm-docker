<h1 align="center">
    Sasm-docker
</h1>

<p align="center">
    <!--a href="https://www.gnu.org/licenses/agpl-3.0">
        <img src="https://img.shields.io/badge/License-AGPL%20v3-blue.svg" />
    </a-->
    <a href="https://github.com/keinenclue/sasm-docker/actions/workflows/release.yml/badge.svg">
        <img src="https://github.com/keinenclue/sasm-docker/actions/workflows/release.yml/badge.svg" alt="Badge tests">
    </a>
    <a href="https://goreportcard.com/report/github.com/keinenclue/sasm-docker">
        <img src="https://img.shields.io/badge/go%20report-A-green.svg?style=flat" alt="Go report" />
    </a>
</p>

Sasm-docker simplifies the setup and use of [SASM](https://dman95.github.io/SASM/english.html) by running it inside a docker container and using x11 (X Window System) in order to display the SASM GUI.

# Features

- **Easy setup:** Just install docker and xserver, download the launcher, and you're ready to go
- **Easy updating:** The launcher checks for and downloads updates on every start

## Run
#### Install X server
- On MacOS [follow these steps](https://gist.github.com/paul-krohn/e45f96181b1cf5e536325d1bdee6c949) but use `xhost +localhost` instead of `xhost +$(hostname).local`.
- On Windows ...
- On Linux you are probaply ready to go :)

#### Install docker
- On MaxOS ...
- On Windows ...
- On Linux [follow these steps](https://docs.docker.com/engine/install)

#### Install the launcher
-  Download the binary for your OS over here: [https://github.com/keinenclue/sasm-docker/releases/latest](https://github.com/keinenclue/sasm-docker/releases/latest)
