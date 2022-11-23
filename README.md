# goexpert-desafio-client-server-api
Projeto do desafio client-server-api do curso Go Expert da Full Cycle.

## server
Para rodar o server, a partir da raiz do projeto, basta dar um ```go run server/server.go```.

Para fazer consultas no banco de dados SQLite3, basta usar o docker (de dentro da pasta *server*):

    docker run --rm -it -v "$(pwd)/database:/workspace" -w /workspace keinos/sqlite3
    
Você vai acessar o container e aí é só mandar um ```.open dollar.db``` para abrir a base de dados. Por fim, é SQL normal (lembre-se do ; no final das consultas).

## client
O cliente é mais simples, é só dar um ```go run client.go```, a partir da raiz do projeto. Cada vez que ele roda, ele vai no server, busca a cotação e grava no arquivo (faz um append). O arquivo fica dentro da pasta ```/client```.


