# Best Route

Best Route é um programa desenvolvido para consultar o menor preço possível de uma viagem informando os preços de viagem de locais relacionados e o mais importante, o local de ínicio e de fim da viagem, para assim o sistema responder a rota mais barata.

## Exemplo
Para salvar as rotas e seus respectivos preços deve ser salvo em csv e com o seguinte padrão:
Considerando que para viajar de GRU para BRC o custo é $10
```
GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20
```
Para viajar de GRU para CDG existem as seguintes rotas:

1. GRU - BRC - SCL - ORL - CDG ao custo de $40
2. GRU - ORL - CGD ao custo de $64
3. GRU - CDG ao custo de $75
4. GRU - SCL - ORL - CDG ao custo de $45

O melhor preço é da rota 1 logo, o output da consulta será GRU - BRC - SCL - ORL - CDG.

## Como executar o programa
Para usar este projeto é necessário que a linguagem ***[GO](https://golang.org/)*** esteja instalada em seu computador

Após baixar este repositório do Github, entre na raiz do projeto, e se caso preferir pode compilar o programa novamente usando o comando: "go build -o main".

Com o programa já compilado, devemos inicia-lo com o seguinte comando:
> ./main input-routes.csv

(Sendo "input-routes.csv" o arquivo de rotas inicial, o nome e o arquivo pode ser alterado, mas deve seguir o mesmo padrão mencionado anteriormente.)

Após isto você escolherá a forma de consulta que você deseja utilizar, sendo 1 para CLI e 2 para API

### CLI
A famosa interface de linha de comando, para consultar a melhor rota você deve informar o local inical e o local final desta forma: "INICIAL-FINAL", por exemplo: "GRU-CDG"

### API
Depois de escolhida a interface de aplicação, por padrão o servidor será iniciado na rota "localhost:3000", no qual escutará as seguintes rotas:

##### Melhor Rota
Para consultar a melhor rota, envie uma requisição com método "GET" para o endereço **"/best"** (localhost:3000/best), com os seguinte body:
```
{
	"start": "GRU",
	"target": "CDG"
}
```
Em caso de sucesso, a resposta será:
```
{
  "route": [
    "GRU",
    "BRC",
    "SCL",
    "ORL",
    "CDG"
  ],
  "cost": 40
}
```
#### Adicionar Rota
Para adicionar uma nova rota, envie uma requisição com método "POST" para o endereço **"/add"** (localhost:3000/add), com o seguinte body:

```
{
	"start": "GRU",
	"target": "CDG",
	"cost": 40
}
```
Em caso de sucesso, a API retornará a rota que foi adicionada.

## Detalhes sobre o projeto e seu desenvolvimento
Foi muito divertido desenvolver este desafio, principalmente quando me deparei com o problema e relembrei das aulas que já tive sobre um algoritmo em específico que soluciona o problema, o poderosíssimo Dijkstra, que calcula o menor caminho entre dois nós de um Grafo.
Inicie pesquisando para estudar melhor sobre o algoritmo e encontrei uma biblioteca abandonada em específica que me chamou atenção, fiz um fork e concertei um problema que havia na validação dos nós que já pertenciam ao grafo, e utilizei esse fork para continuar o desafio.
Para implementar o desafio decidi utilizar o mínimo de bibliotecas externas da linguagem, o resultado foi que o programa possui apenas duas dependências, o meu fork do algoritmo de Dijkstra e a biblioteca testify, que é uma preferência minha para implementar os testes.
Decidi também manter o menos acoplado possível, principalmente a implementação para ler e inserir dados, que foi pensada para que possa ser feita uma nova implementação para um banco de dados por exemplo, sem afetar o restante do código.Da mesma maneira foi feito para a uma entidade chamada "route_calculator" na qual internamente está utilizando a implementação do algoritmo de Dijkstra, alterar a implementação seria algo "relativamente simples".
Todo o código foi implementado pensando em testes unitários e mantendo uma responsabilidade mínima para cada método e entidade, para que possíveis refactors futuros fossem feitos com mais simplicidade.

### Estrutura de pacotes

```
- database
    - csv (implementação para ler e inserir em um arquivo csv)
        client.go
    database.go (definição da interface)
- models
    route.go (entidade de rota)
    response-error.go (entidade de error interna)
    validations.go (métodos para validar as requisições)
- route_calculator (entidade que calcula a melhor rota)
    - djk (implementação utilizando o algoritmo de Dijkstra)
        client.go
    route_calculator (definição da interface)
main.go
api.go
cli.go
```