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

More information about the tool can be found in the [operations](docs/operations/delete-loop_en.md) documentation, or by calling `tempdel --help` provides more information.

Documentation on developing `tempdel` can be found in [English](docs/developing_en.md) and [German](docs/developing_de.md).

---
## What is the Cloudogu EcoSystem?
The Cloudogu EcoSystem is an open platform, which lets you choose how and where your team creates great software. Each service or tool is delivered as a Dogu, a Docker container. Each Dogu can easily be integrated in your environment just by pulling it from our registry.

We have a growing number of ready-to-use Dogus, e.g. SCM-Manager, Jenkins, Nexus Repository, SonarQube, Redmine and many more. Every Dogu can be tailored to your specific needs. Take advantage of a central authentication service, a dynamic navigation, that lets you easily switch between the web UIs and a smart configuration magic, which automatically detects and responds to dependencies between Dogus.

The Cloudogu EcoSystem is open source and it runs either on-premises or in the cloud. The Cloudogu EcoSystem is developed by Cloudogu GmbH under [AGPL-3.0-only](https://spdx.org/licenses/AGPL-3.0-only.html).

## License
Copyright © 2020 - present Cloudogu GmbH
This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3.
This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
You should have received a copy of the GNU Affero General Public License along with this program. If not, see https://www.gnu.org/licenses/.
See [LICENSE](LICENSE) for details.


---
MADE WITH :heart:&nbsp;FOR DEV ADDICTS. [Legal notice / Imprint](https://cloudogu.com/en/imprint/?mtm_campaign=ecosystem&mtm_kwd=imprint&mtm_source=github&mtm_medium=link)
