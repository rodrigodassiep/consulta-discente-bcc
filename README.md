# Sistema de Consulta Discente - IME-USP

Este repositório contém o código-fonte do Sistema de Consulta Discente, desenvolvido para os alunos do Bacharelado em Ciência da Computação do Instituto de Matemática e Estatística da Universidade de São Paulo (IME-USP).

## Descrição

O Sistema de Consulta Discente é uma aplicação web que permite aos alunos do curso de Ciência da Computação do IME-USP avaliar professores e disciplinas. Entre as funcionalidades oferecidas pelo sistema, destacam-se:

- Avaliacao de professores
- Avaliacao de disciplinas
- Coleta de feedback anonimo dos alunos
- Geracao de relatorios de avaliacao
- Visualizacao de resultados de avaliacoes anteriores

## Aplicacao em Producao

O sistema esta hospedado em: **https://consulta-discente-bcc.up.railway.app/**

## Tecnologias Utilizadas

- **Frontend:** SvelteKit 2.0 + Tailwind CSS
- **Backend:** Golang + Gin Framework
- **Banco de Dados:** PostgreSQL
- **Autenticacao:** JWT (JSON Web Tokens)
- **Hospedagem:** Railway

---

## Como Testar o Sistema

O sistema possui tres tipos de usuarios (roles) com diferentes permissoes e funcionalidades. Abaixo estao as credenciais de teste e o que cada tipo de usuario pode fazer.

### Credenciais de Teste

| Role | Email | Senha |
|------|-------|-------|
| **Administrador** | `admin@usp.br` | `admin123` |
| **Professor** | `maria.silva@usp.br` | `prof123` |
| **Professor** | `joao.santos@usp.br` | `prof123` |
| **Professor** | `ana.costa@usp.br` | `prof123` |
| **Estudante** | `pedro.oliveira@usp.br` | `student123` |
| **Estudante** | `julia.ferreira@usp.br` | `student123` |
| **Estudante** | `lucas.almeida@usp.br` | `student123` |
| **Estudante** | `carla.mendes@usp.br` | `student123` |
| **Estudante** | `rafael.lima@usp.br` | `student123` |

---

### Fluxo de Teste por Role

#### 1. Estudante

**Acesso:** Login com credenciais de estudante (ex: `pedro.oliveira@usp.br` / `student123`)

**Funcionalidades disponiveis:**
- Visualizar pesquisas disponiveis para responder
- Responder pesquisas de avaliacao de disciplinas
- Ver historico de respostas enviadas
- Visualizar suas matriculas

**Como testar:**
1. Faca login com uma conta de estudante
2. No dashboard, veja as pesquisas ativas
3. Clique em uma pesquisa para responder
4. Responda as questoes (NPS, avaliacao por estrelas, texto livre, multipla escolha)
5. Envie suas respostas

> **Nota:** As respostas sao anonimas - professores nao conseguem identificar qual estudante enviou cada resposta.

---

#### 2. Professor

**Acesso:** Login com credenciais de professor (ex: `maria.silva@usp.br` / `prof123`)

**Funcionalidades disponiveis:**
- Visualizar suas disciplinas
- Criar novas pesquisas de avaliacao
- Adicionar questoes as pesquisas (4 tipos: NPS, Rating, Texto Livre, Multipla Escolha)
- Visualizar respostas anonimas dos alunos
- Ver estatisticas agregadas (media de notas, distribuicao, NPS score)

**Como testar:**
1. Faca login com uma conta de professor
2. No dashboard, veja suas disciplinas e pesquisas existentes
3. Crie uma nova pesquisa:
   - Clique em "Criar Pesquisa"
   - Selecione a disciplina
   - Defina titulo e datas
4. Adicione questoes:
   - Clique em "Gerenciar Questoes"
   - Adicione diferentes tipos de perguntas
5. Visualize respostas:
   - Clique em "Ver Respostas" em uma pesquisa
   - Veja estatisticas agregadas por questao
   - Para questoes de texto, clique para ver todas as respostas anonimas

---

#### 3. Administrador

**Acesso:** Login com credenciais de administrador (`admin@usp.br` / `admin123`)

**Funcionalidades disponiveis:**
- Gerenciar semestres (criar, ativar/desativar)
- Gerenciar disciplinas (criar, editar, atribuir professores)
- Gerenciar usuarios (visualizar, aprovar roles solicitados)
- Gerenciar matriculas de estudantes
- Visualizar todas as pesquisas e respostas do sistema

**Como testar:**
1. Faca login com a conta de administrador
2. No dashboard, explore as diferentes secoes:
   - **Semestres:** Crie um novo semestre ou ative/desative existentes
   - **Disciplinas:** Adicione novas disciplinas e atribua professores
   - **Usuarios:** Veja todos os usuarios cadastrados
   - **Matriculas:** Vincule estudantes a disciplinas/semestres

---

### Tipos de Questoes

O sistema suporta 4 tipos de questoes nas pesquisas:

| Tipo | Descricao | Visualizacao |
|------|-----------|--------------|
| **NPS (0-10)** | Net Promoter Score | Media + Score NPS (Promotores - Detratores) |
| **Rating (1-5)** | Avaliacao por estrelas | Media + Distribuicao por estrela |
| **Texto Livre** | Resposta aberta | Lista de respostas anonimas |
| **Multipla Escolha** | Opcoes pre-definidas | Grafico de distribuicao |

---

### Testando o Fluxo Completo

Para testar o sistema de ponta a ponta:

1. **Admin:** Crie um semestre ativo e uma disciplina
2. **Admin:** Matricule um estudante na disciplina
3. **Professor:** Crie uma pesquisa para a disciplina
4. **Professor:** Adicione questoes de diferentes tipos
5. **Estudante:** Responda a pesquisa
6. **Professor:** Visualize as respostas anonimas e estatisticas

---

### PWA (Progressive Web App)

O sistema pode ser instalado como aplicativo no celular ou desktop:

**No iPhone/iPad:**
1. Abra o site no Safari
2. Toque no botao de compartilhar
3. Selecione "Adicionar a Tela de Inicio"

**No Android:**
1. Abra o site no Chrome
2. Toque no menu (3 pontos)
3. Selecione "Instalar aplicativo"

**No Desktop (Chrome/Edge):**
1. Clique no icone de instalacao na barra de endereco
2. Confirme a instalacao

---

## Desenvolvimento Local

### Requisitos
- Go 1.24+
- Node.js 18+
- PostgreSQL

### Backend
```bash
cd server
cp .env.example .env
# Configure as variaveis de ambiente
go run .
```

### Frontend
```bash
cd client
npm install
npm run dev
```

---

## Licenca

Este projeto esta licenciado sob a Licenca MIT. Veja o arquivo LICENSE para mais detalhes.
