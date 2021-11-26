<h1 align="center">
    <img src="https://github.com/keinenclue/sasm-docker/blob/main/launcher/Icon.png?raw=true" width="200px" ></img><br>
    Sasm-docker
</h1>

<p align="center">
    <!--a href="https://www.gnu.org/licenses/agpl-3.0">
        <img src="https://img.shields.io/badge/License-AGPL%20v3-blue.svg" />
    </a-->
    <a href="https://github.com/keinenclue/sasm-docker/releases/latest">
        <img src="https://img.shields.io/github/v/release/keinenclue/sasm-docker?logo=github&logoColor=white" alt="GitHub release"/>
    </a>
    <a href="https://github.com/keinenclue/sasm-docker/actions/workflows/release.yml/badge.svg">
        <img src="https://github.com/keinenclue/sasm-docker/actions/workflows/release.yml/badge.svg" alt="Badge tests">
    </a>
    <a href="https://goreportcard.com/report/github.com/keinenclue/sasm-docker">
        <img src="https://goreportcard.com/badge/github.com/keinenclue/sasm-docker" alt="Go report" />
    </a>
</p>

Sasm-docker simplifies the setup and use of [SASM](https://dman95.github.io/SASM/english.html) by running it inside a docker container and using x11 (X Window System) in order to display the SASM GUI.

## Features

- **Easy setup:** Just install docker and xserver, download the launcher, and you're ready to go
- **Easy updating:** The launcher checks for and downloads updates on every start

## Run

#### Install X server

These are just examples, you don't have to use these installation methods.  
You just need to have X server on your system in the end.

- On MacOS
    1. Install XQuartz:
        - With the dmg from the [official website](https://www.xquartz.org/releases/index.html)
        - Or with the MacPorts instructions down below on the website
        - Or use Homebrew `brew cask install xquartz`
    2. Launch XQuartz. Under the XQuartz menu, select Preferences
    3. Go to the security tab and ensure "Allow connections from network clients" is checked.
    4. Restart XQuartz.
- On Windows
  - [Install VcXsrv Windows X Server](https://sourceforge.net/projects/vcxsrv/)
  - Or install through Chocolatey `choco install vcxsrv`
- On Linux you are probaply ready to go :)

#### Install docker

- On MacOS [follow these steps](https://docs.docker.com/desktop/mac/install)
- On Windows [follow these steps](https://docs.docker.com/desktop/windows/install)
- On Linux [follow these steps](https://docs.docker.com/engine/install)

#### Install the launcher

- Download the binary for your OS over here: [https://github.com/keinenclue/sasm-docker/releases/latest](https://github.com/keinenclue/sasm-docker/releases/latest)
- On MacOS open the dmg and drag the application to you applications folder

#### Known issues

- If you get the warning `warning: creating DT_TEXTREL in a PIE` add the flag `-fPIC` to `Linking options`
- Debugging of a project which uses multiple files will not work ([this is a bug in sasm](https://github.com/Dman95/SASM/issues/102))
- Debugging will not work if you have a function which is not delcared as global
- If the debugging doesn't work, clear your Data path(it's best that you completly delete the folder and then recreate it and put your files back in)

## Screenshots

<table align="center">
    <tr>
        <td align="center">
               <a href="https://user-images.githubusercontent.com/30153207/140638832-c3f91a51-81a3-44db-8a1e-0192c9fe9ec5.png">
                   <img src="https://user-images.githubusercontent.com/30153207/140638832-c3f91a51-81a3-44db-8a1e-0192c9fe9ec5.png" width="500px" alt="Screenshot launch" />
              </a>
        </td>
        <td align="center">
            <a href="https://user-images.githubusercontent.com/30153207/140639058-fed938e4-2a66-4a42-849d-86c5a4fb61a6.png" >
                <img src="https://user-images.githubusercontent.com/30153207/140639058-fed938e4-2a66-4a42-849d-86c5a4fb61a6.png" width="500px" alt="Screenshot pulling" />
            </a>
        </td>
    </tr>
    <tr>
       <td align="center">
            <a href="https://user-images.githubusercontent.com/30153207/140663399-7cb1af3b-ce8f-4551-9aae-954619710607.png">
                <img src="https://user-images.githubusercontent.com/30153207/140663399-7cb1af3b-ce8f-4551-9aae-954619710607.png"  width="500px" alt="Screenshot logs" />
            </a>
        </td>
        <td align="center">
            <a href="https://user-images.githubusercontent.com/30153207/140639009-8f6dd323-12aa-4e8f-9d6d-afbcfed45e32.png" >
                <img src="https://user-images.githubusercontent.com/30153207/140639009-8f6dd323-12aa-4e8f-9d6d-afbcfed45e32.png" width="500px" alt="Screenshot settings" />
            </a>
        </td>
    </tr>
</table>
