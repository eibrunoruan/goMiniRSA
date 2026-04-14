# miniRSA

Projeto de demonstração de RSA em Go, com:
- geração de chaves
- criptografia/decriptografia
- ataque de fatoração (bônus)

## Pré-requisitos

- Go instalado (o projeto está com `go 1.26.1` no [go.mod](go.mod))

## Como rodar

Na raiz do projeto, execute:

```bash
go run ./cmd/rsa-app
```

Isso vai:
1. Gerar um par de chaves RSA (16 bits, para fins didáticos)
2. Criptografar o valor de teste (`54321`)
3. Decriptar o valor
4. Tentar quebrar a chave por fatoração e decriptar novamente

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
- [internal/rsa/breaker.go](internal/rsa/breaker.go): tentativa de quebra por fatoração
