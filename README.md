# miniRSA

Projeto de demonstração de RSA em Go, com:
- geração de chaves
- criptografia/decriptografia
- ataque de fatoração (bônus) com backend nativo rápido

## Pré-requisitos

- Go instalado (o projeto está com `go 1.26.1` no [go.mod](go.mod))
- `msieve` instalado para a via rápida de fatoração

No macOS com Homebrew:

```bash
brew install msieve
```

## Como rodar

Na raiz do projeto, execute:

```bash
go run ./cmd/rsa-app
```

Isso vai:
1. Gerar um par de chaves RSA (64 bits por padrão, para manter a fatoração interativa)
2. Criptografar o valor de teste (`54321`)
3. Decriptar o valor
4. Tentar quebrar a chave por fatoração e decriptar novamente

Para testar outros tamanhos de chave:

```bash
RSA_BITS=128 go run ./cmd/rsa-app
```

Observação: a aplicação agora tenta usar `msieve` com ECM primeiro, aproveitando todos os cores disponíveis. Se o binário não estiver instalado, ela cai automaticamente para o Pollard-Brent em Go puro como fallback.

## Build (opcional)

Para compilar o executável:

```bash
go build -o mini-rsa ./cmd/rsa-app
```

Para executar depois:

```bash
./mini-rsa
```

## Estrutura

- [cmd/rsa-app/main.go](cmd/rsa-app/main.go): ponto de entrada
- [internal/rsa/controller.go](internal/rsa/controller.go): orquestra o fluxo
- [internal/rsa/service.go](internal/rsa/service.go): lógica RSA
- [internal/rsa/presenter.go](internal/rsa/presenter.go): saída no console
- [internal/rsa/breaker.go](internal/rsa/breaker.go): tentativa de quebra por fatoração, priorizando `msieve`
