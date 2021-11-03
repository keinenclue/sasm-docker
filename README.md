# sasm-docker

A docker image in order to run [SASM](https://dman95.github.io/SASM/english.html)
This uses x11 (X Window System) in order to display the SASM GUI.
## Run

Just start the docker container:

```bash
docker compose up
```

## Install x11 on MacOS

[Follow these steps](https://gist.github.com/paul-krohn/e45f96181b1cf5e536325d1bdee6c949) but use `xhost +localhost` instead of `xhost +$(hostname).local`.

## Releases

Releases are pushed to the GitHub container registry as `ghcr.io/keinenclue/sasm-docker`.

Images are built and pushed automatically on every git tag starting with `v`.
