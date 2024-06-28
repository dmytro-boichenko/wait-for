# CLI `wait-for` tool

## Preface
The tool is used for controlling in docker-compose builds for correct waiting for required resources like databases.

The current implementation supports MySQL, Postgres, Redis and Elasticsearch.

## Usage
```shell
Usage of ./wait-for:
  -limit int
    	number of repeats (default 30)
  -service string
    	service name. Mandatory. Possible values are: [mysql elasticsearch redis]
  -timeout int
    	timeout in seconds between repeats (default 1)
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