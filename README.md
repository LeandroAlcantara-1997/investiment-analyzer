# investiment-analyzer

## Ferramentas utilizadas:

* Docker;
* Golang 1.23;
* Jaeger para traces;
* Redis;
* Postgres;
* OpenTelemetry;

## Como executar?

Tenha instalado a lingaguem golang na versão 1.23 ou posterior e o docker.
Adicione um arquivo .env na pasta build
API_PORT=8080
API_VERSION=
API_NAME=
DB_NAME=investment-analyzer
DB_USER=user
DB_PASSWORD=passw0rd
DB_HOST=postgres-database
DB_PORT=5432
CACHE_HOST=redis-cache
CACHE_PORT=6379
CACHE_PASSWORD=passw0rd
CACHE_READ_TIMEOUT=2
CACHE_WRITE_TIMEOUT=2
ALLOW_ORIGINS=http://localhost:3000,https://localhost:8080
OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4318
ENVIRONMENT=local

~~~ make 
make docker-build
~~~

~~~ make 
make docker-up
~~~


Você pode consultar a documentação swagger clicando (swagger)[http://localhost:8080/swagger/index.html#/]

O Jaeger pode ser visto utilizando esse endereço http://localhost:16686