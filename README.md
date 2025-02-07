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
[.env.example](../../Documentos/investiment-analyzer/.env.example)

~~~ make 
make docker-build
~~~

~~~ make 
make docker-up
~~~


Você pode consultar a documentação swagger clicando (swagger)[http://localhost:8080/swagger/index.html#/]

O Jaeger pode ser visto utilizando esse endereço http://localhost:16686