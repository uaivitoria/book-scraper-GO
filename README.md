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

**Go:** Eu já havia feito o desafio em python, mas refleti e resolvi me desafiar e fazer na linguagem sugerida pelo desafio. 

**Sem frameworks externos:** usei apenas `golang.org/x/net/html` para parsing de HTML e a biblioteca padrão do Go para todo o resto. Menos dependências = menos vulnerabilidades e builds mais rápidos.

**Funções recursivas para navegação HTML:** `findFirst` e `findAll` navegam a árvore HTML de forma limpa, equivalente ao BeautifulSoup do Python.

**Multi-stage build no Docker:** o stage `builder` compila o binário e o stage `runtime` usa só o Alpine puro com o binário — imagem final extremamente leve.

**Usuário não-root no container:** boa prática de segurança.

**Delay entre requisições:** `time.Sleep(500ms)` entre páginas para respeitar o servidor.

**User-Agent identificável:** o header deixa claro que é um scraper educacional.

---

## O que eu faria diferente com mais tempo

- **Concorrência com goroutines:** coletar múltiplas páginas em paralelo usando goroutines e channels, aproveitando o ponto forte do Go
- **Retry com backoff:** retentar automaticamente em caso de falha de rede
- **Página de detalhe:** coletar descrição, UPC e número de resenhas de cada livro
- **IA na extração:** usar um LLM para extrair informações estruturadas das descrições
- **Logging estruturado:** registrar quantidade coletada, erros e tempo total
- **Banco de dados:** persistir os dados em PostgreSQL
- **docker-compose.yml:** orquestrar scraper + banco com um único comando

---

## Como usei IA neste desafio

Usei o Claude (Anthropic) como ferramenta de apoio durante todo o desenvolvimento:

- **Entendimento do desafio:** pedi uma análise do PDF e um passo a passo de execução
- **Aprendizado de Go:** não tinha experiência prévia com Go — o Claude explicou conceitos como structs, interfaces, goroutines e o sistema de módulos
- **Geração do código:** o Claude gerou a estrutura base que revisei e ajustei para garantir que entendia cada parte
- **Testes unitários:** aprendi a criar HTML falso com `html.Parse` para testes sem internet
- **Dockerfile:** entendi o conceito de multi-stage build com binário estático Go
- **Pipeline CI/CD:** entendi as diferenças entre o pipeline Go e Python
- **O que funcionou:** pedir explicações junto com o código, arquivo por arquivo
- **O que não funcionou:** pedir o projeto todo de uma vez gerou código genérico demais
