# confluence-temp-delete-job

A small tool that deletes files (short: `tempdel`).

## Description

Over the course of time, Confluence may carelessly leave temporary files in the directory `/opt/atlassian/confluence/temp`. This tool  runs like a background job and deletes such temp files according a file age policy.

There are three basic input parameters:

1. the directory to be scanned for files
1. the maximum age in hours of files that should not be deleted
1. the time between the scanning intervals in minutes

```bash
tempdel delete-loop --age 12 --interval 60 /opt/atlassian/confluence/temp
```

For more information, calling `tempdel --help` provides more information about the tool.

Documentation on developing `tempdel` can be found in [English](docs/developing_en.md) and [German](docs/developing_de.md).

---
### What is the Cloudogu EcoSystem?
The Cloudogu EcoSystem is an open platform, which lets you choose how and where your team creates great software. Each service or tool is delivered as a Dogu, a Docker container. Each Dogu can easily be integrated in your environment just by pulling it from our registry. We have a growing number of ready-to-use Dogus, e.g. SCM-Manager, Jenkins, Nexus, SonarQube, Redmine and many more. Every Dogu can be tailored to your specific needs. Take advantage of a central authentication service, a dynamic navigation, that lets you easily switch between the web UIs and a smart configuration magic, which automatically detects and responds to dependencies between Dogus. The Cloudogu EcoSystem is open source and it runs either on-premises or in the cloud. The Cloudogu EcoSystem is developed by Cloudogu GmbH under [MIT License](https://cloudogu.com/license.html).

### How to get in touch?
Want to talk to the Cloudogu team? Need help or support? There are several ways to get in touch with us:

* [Website](https://cloudogu.com)
* [myCloudogu-Forum](https://forum.cloudogu.com/topic/34?ctx=1)
* [Email hello@cloudogu.com](mailto:hello@cloudogu.com)

---
&copy; 2020 Cloudogu GmbH - MADE WITH :heart:&nbsp;FOR DEV ADDICTS. [Legal notice / Impressum](https://cloudogu.com/imprint.html)
