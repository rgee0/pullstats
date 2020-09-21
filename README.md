# PullStats

[![OpenFaaS](https://img.shields.io/badge/openfaas-cloud-blue.svg)](https://www.openfaas.com)

A simple OpenFaaS function to consolidate pull counts where Docker images are dispersed across different organisations.  The default is a good example use case, where the OpenFaaS gateway, for historical reasons, resided in the `openfaas` & `functions` organisations.

## Configuration

The orgs to be interrogated can be provided as a comma separated list through the `orgs` env var in `stack.yml`

```
    environment:
      orgs: openfaas,functions
      read_timeout: 30s
      write_timeout: 30s
```

If a larger number of orgs are being queried it is recommended that the function timeouts are tailored accordingly.

## Usage

Once deployed to an OpenFaaS instance the endpoint can be called with:

### No argument / Default output

This will consolidate all the images within the provided orgs.

```sh
$ curl https://rgee0.o6s.io/pullstats -d ''

{
  "total": 32262735,
  "repos": {
    "alertmanager": 4329010,
    "alertmanager-legacy": 121,
    "alexa-leds": 2612,
    "alpine": 3558970,
    "api-key-protected": 95837,
    ...
    "queue-worker": 4620823,
    "resizer": 713307,
    "sentimentanalysis": 35042,
    "webhookstash": 245125,
    "wordcount": 21619
  }
}
```

Try piping through jq to prettify

```
$ curl -s https://rgee0.o6s.io/pullstats -d '' | jq .
```


### Specific image name

This will consolidate only the provided image name.
> Note: if the image name isnt found within the set then the output will revert to default

```sh
$ curl https://rgee0.o6s.io/pullstats -d 'gateway'

7218749
 
```
