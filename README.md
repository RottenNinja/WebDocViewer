# WebDocViewer

Web based markdown viewer for documentation

## Setup
Just add any volume containing mark down files to the /www/docs of the container and it will display them on port 8080


### Docker compose
```yaml
version: "3.2"
services:

  webdocs:
    container_name: webdocs
    image: hgithub/webdocviewer:latest
    ports:
      - '8080:8080'
    volumes:
      - ./example_docs:/www/docs

  #Optional
  watchtower:
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --interval 60  --include-restarting --include-stopped --revive-stopped --cleanup webdocs
```



### Required config

Mandatory `config.wdv.json` file must exist in the root of the document folder containing the following:
```json
{
  "title": "Title goes here",
  "logoUri": "docs/logo-dark.png"
}
```


## Optional folder configs

*Be warned this json file is not parsed and is directly sent to the html javascript of the browser to deal with*

Configs exist for folders such as having a `.wdv.json` file in the root allows the following options:
```json
{
  "order": {
    "first": ["Home.md","About.md"],
    "last": ["FAQ.md"],
    "rest": "ASC"
  }
}
```
Here you can state the order that the navigation will display the files in `rest` either values `ASC` or `DESC` as well as if you want any files to come "FIRST" or "LAST" in the navigation. Just list the file names in the arrays `fist` and `last` respectively.


## Optional file configs

*Be warned this json file is not parsed and is directly sent to the html javascript of the browser to deal with*

For files there is also the option to have the `filename.wdv.json` such as `Home.md.mdv.json` that can contain the following config:
```json
{
  //Empty for now
}
```