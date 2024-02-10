# Desafio de Sistemas de Cotação de Moedas em Go
## Visão Geral

Este repositório contém duas aplicações em Go que trabalham em conjunto para fornecer cotações de câmbio de USD para BRL:

- **client.go**: Solicita cotações de câmbio ao `server.go`.
- **server.go**: Obtém cotações de câmbio de uma API externa e retorna o resultado ao cliente.

## Arquitetura do Sistema

A imagem a seguir apresenta a arquitetura do sistema, detalhando a interação entre `client.go` e `server.go`, e como eles se conectam com a API externa e o banco de dados SQLite.

## Como Funciona
**server.go**

- Consome a API de câmbio de Dólar para Real do endereço https://economia.awesomeapi.com.br/json/last/USD-BRL.
- Retorna a taxa de câmbio no formato JSON para o cliente.
- Utiliza o pacote context para implementar timeouts:
  - Um máximo de 200ms para chamar a API de cotação da moeda. 
  - Um máximo de 10ms para persistir os dados no banco de dados SQLite.
- Escuta na porta 8080 e serve o endpoint /cotacao.

**client.go**
- Faz uma requisição HTTP ao `server.go` para obter a cotação atual de USD para BRL.
- Recebe apenas o valor de "lance" (bid) da taxa de câmbio da resposta do servidor.
- Implementa um timeout máximo de 300ms para receber o resultado do `server.go`, usando o pacote `context`.
- Escreve a cotação atual em um arquivo chamado `cotacao.txt` no formato: `Dólar: {valor}`.

## Requisitos
 - Linguagem de programação Go
- SQLite
- Consumo de API externa
- Uso do pacote context para gerenciamento de timeouts
- Log de erro para tempo de execução insuficiente

## Configuração

1. Certifique-se de ter o Go instalado na sua máquina.
2. Clone este repositório.
3. Navegue até o diretório do repositório.

## Executando as Aplicações

Para iniciar o servidor:
```go
go run server.go
```

Para executar o cliente:
```go
go run client.go
```
