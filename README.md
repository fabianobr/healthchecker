# healthchecker
O objetivo deste projeto é simples, avaliar o status de vários serviços (ou microserviços). Muito útil quando há muitos serviços a serem avaliados, evitando de conectar um a um via Postman, Insomnia ou outras ferramentas. O segundo objetivo é aprender Go (ou Goland).

## Como funciona
A aplicação funciona como um web service (porta 3001) que faz praticamente 2 coisas:
  1. Ler o arquivo de configurações "services.json"
  2. Checar a saúde de cada aplicação listada no arquivo de configurações

## Passos de funcionamento
  1. Ler e materializar arquivos de configurações, conhecer todos os serviços
  2. Checar a saúde de cada aplicação

## Endpoints
  * /: Lista todas as aplicações e sua saúde, no formato JSON, vide "StatusTime" e "StatusText";
  * /check: Força outra checagem de todos os serviços
  * /reload: Força a leitura do arquivo de configurações "services.json"

## TO-DOs
  1. Generalizar o nome do arquivo JSON via argumentos na linha de comando/
  2. Generalizar o numero da porta HTTP;
