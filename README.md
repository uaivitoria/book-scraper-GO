# 📚 Book Scraper — Go

Scraper desenvolvido como desafio técnico para o Programa Trainee Crawler/RPA & IA da IN8 Holding / Devnology.

Coleta dados de todos os livros do site [books.toscrape.com](https://books.toscrape.com) e salva em JSON e CSV.

Desenvolvido em **Go** usando apenas a biblioteca padrão + `golang.org/x/net/html`.

---

## Como rodar localmente

### Sem Docker

```bash
# 1. Baixe as dependências
go mod download

# 2. Execute o scraper
go run main.go
```

Os arquivos `books.json` e `books.csv` serão gerados na pasta atual.

### Com Docker

```bash
# 1. Build da imagem
docker build -t book-scraper-go .

# 2. Rodar o container
docker run --rm -v $(pwd)/output:/app/output book-scraper-go
```

---

## Estrutura dos dados extraídos

### Schema JSON

```json
[
  {
    "title": "A Light in the Attic",
    "price_gbp": 51.77,
    "rating": 3,
    "in_stock": true
  }
]
```

### Schema CSV

| Coluna | Tipo | Descrição |
|---|---|---|
| `title` | string | Título completo do livro |
| `price_gbp` | float | Preço em libras esterlinas |
| `rating` | int (1–5) | Avaliação numérica |
| `in_stock` | bool | Disponibilidade em estoque |

---

## Como o pipeline funciona

| Stage | O que faz |
|---|---|
| **lint** | Roda o `go vet` para verificar erros e boas práticas. Falha se houver problemas. |
| **test** | Executa os testes unitários com `go test`. Falha se algum teste não passar. |
| **build** | Compila o binário, constrói a imagem Docker e faz push para o GitLab Registry. |
| **deploy** | Simula o deploy no AWS ECS. Roda apenas na branch `main`. |

---

## Decisões técnicas

**Limguagem Go:**
O desafio sugeria Go como linguagem principal. Mesmo sem experiência prévia com a linguagem,
decidi encarar o desafio como uma oportunidade de aprendizado real. Com o apoio de IA para
entender a sintaxe e os conceitos, consegui desenvolver o projeto do zero e compreender o
que estava sendo feito em cada etapa.

**Arquitetura de pastas:**
Organizei o projeto seguindo o padrão de mercado do Go, separando em `cmd/` para o ponto
de entrada e `internal/` para a lógica de negócio. Além de deixar o projeto mais legível,
essa estrutura permite escalar o projeto facilmente e adicionar novos pacotes, uma API ou
novos exporters sem bagunçar o código.

**Correção do leak de memória no `findAll`:**
Analisando o código percebi que a função acumulava slices a cada chamada recursiva, o que
em projetos maiores poderia causar sobrecarga de memória e corromper os dados. Refatorei
usando `walkNodes` com ponteiro para a slice, garantindo que todas as chamadas escrevem no
mesmo lugar em vez de criar cópias desnecessárias.

**Sem frameworks externos:**
Usei apenas `golang.org/x/net/html` para parsing de HTML e a biblioteca padrão do Go para
todo o resto. Menos dependências significa menos vulnerabilidades e builds mais rápidos.

**Delay entre requisições:**
Adicionei `time.Sleep(500ms)` entre páginas para respeitar o servidor, seguindo as boas
práticas éticas de web scraping descritas no desafio.

## O que eu faria diferente com mais tempo


- **API REST:** transformar o scraper em uma API HTTP com endpoints para consultar, 
filtrar e disparar a coleta. Com docker-compose orquestrando a API + PostgreSQL, 
qualquer sistema poderia consumir os dados via HTTP em vez de depender de arquivos JSON/CSV
- **docker-compose.yml:** orquestrar API + banco de dados com um único comando `docker-compose up`
- **Logging estruturado:**
Adicionaria logs detalhados usando a biblioteca `slog` (padrão do Go desde 1.21),
registrando em tempo real a quantidade de livros coletados por página, erros de
requisição com o motivo da falha, e o tempo total de execução. Isso facilita muito
o monitoramento em produção e a identificação de problemas sem precisar rodar o
scraper novamente.

---

## Como usei IA neste desafio

Usei o Claude (Anthropic) como ferramenta de apoio durante todo o desenvolvimento:

- **Entendimento do desafio:** pedi uma análise do PDF e um passo a passo de execução
- **Aprendizado de Go:** não tinha experiência prévia com Go. O Claude explicou conceitos como structs, interfaces, goroutines e o sistema de módulos, arquitetura de pastas, limpeza de variavel
- **Geração do código:** o Claude gerou a estrutura base que revisei e ajustei para garantir que entendia cada parte
- **Testes unitários:** aprendi a criar HTML falso com `html.Parse` para testes sem internet
- **Dockerfile:** entendi o conceito de multi-stage build com binário estático Go
- **Pipeline CI/CD:** entendi a boa aplicação de um CI/CD
- **O que funcionou:** pedir explicações junto com o código, arquivo por arquivo, estudei a arquitetura de pastas da linguagem GO, para melhor organizar e estruturar o teste
- **O que não funcionou:** pedir o projeto todo de uma vez gerou código genérico demais 