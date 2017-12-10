# ![Spirala](https://user-images.githubusercontent.com/3391295/33789830-b4e396c6-dc7b-11e7-97db-e47d4e22aca6.png)
[![GitHub last commit](https://img.shields.io/github/last-commit/lnsp/spirala.svg?style=flat-square)](https://github.com/lnsp/spirala)
[![Docker Stars](https://img.shields.io/docker/stars/lnsp/spirala.svg?style=flat-square)](https://hub.docker.com/r/lnsp/spirala/)
[![Docker Automated build](https://img.shields.io/docker/automated/lnsp/spirala.svg?style=flat-square)](https://hub.docker.com/r/lnsp/spirala/)
[![license](https://img.shields.io/github/license/lnsp/spirala.svg?style=flat-square)](https://github.com/lnsp/spirala)

This is the repository of the Spirala Project, a self-hosted private cloud platform based on Docker Swarm.
It was created to replace easy-to-use systems like Dokku that cannot be used for Docker Swarm. The release target is to provide
a fully integrated private cloud using only Open-Source components.

## Roadmap
- [ ] Display information from Docker Client
- [ ] Allow use of SSL certificates
- [ ] Create and edit Docker components
- [ ] Use Stack templates for easy service creation
- [ ] Add workflow interface to allow visual infrastructure composition

## Installation
Before you run the platform on your Docker host, you should be sure of the following things.
- your host is part of a Docker Swarm
- your host has the role Manager

You can then just run the following script on your host.
```bash
docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8080:8080 lnsp/spirala /go/bin/spirala --endpoint unix:///var/run/docker.sock
```
You can now open your browser and point it to `http://localhost:8080`!

## Copyright
The Spirala project is developed by the Spirala community and administrated by the Spirala maintainers.
The logo credits go to Simon Martin from the Noun Project for his beautiful Penrose Pentagon.