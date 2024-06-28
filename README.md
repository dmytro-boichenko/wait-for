# CLI `wait-for` tool

## Preface
The wait-for tool used for controlling docker-compose builds for correct waiting for required resources like databases.

The current implementation supports MySQL, Postgres, Redis and Elasticsearch.

## Usage
```shell
$ ./wait-for -h                                                                                                                                                                                                                                                        5s 23:12:16
NAME:
   wait-for - used for controlling in docker-compose builds for correct waiting for required resources like databases

USAGE:
   wait-for [global options] command [command options] 

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --service value, -s value  service name. Supported: mysql, elasticsearch, redis, postgres
   --timeout value, -t value  timeout in seconds between repeats (default: 1)
   --limit value, -l value    number of repeats (default: 30)
   --help, -h                 show help

```

For example:  
&nbsp;&nbsp;&nbsp;```./wait-for -service mysql postgres```

### Required environment variables:
#### MySQL
|Environment Variable|Description|Default value|
| --- | --- | --- |
|`DB_HOST`|Host|`localhost`|
|`DB_PORT`|Port|3306|
|`DB_USER`|Username|`user`|
|`DB_PASSWORD`|Password|`password`|
|`DB_NAME`|Database name|`db`|
#### Elasticsearch
|Environment Variable|Description|Default value|
| --- | --- | --- |
|`ELASTIC_HOST`|Host name with schema|`http://localhost`|
|`ELASTIC_PORT`|Port|9200|
#### Redis
|Environment Variable|Description|Default value|
| --- | --- | --- |
|`REDIS_HOST`|Host|`localhost`|
|`REDIS_PORT`|Port|6379|
#### Postgres
|Environment Variable|Description|Default value|
| --- | --- | --- |
|`DB_HOST`|Host|`localhost`|
|`DB_PORT`|Port|5432|
|`DB_USER`|Username|`postgres`|
|`DB_PASSWORD`|Password|`postgres`|
|`DB_NAME`|Database name|`db`|