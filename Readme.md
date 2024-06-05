### Companion AI

An AI companion app. User can create a  story, in that story, user can chat with AI. Currently support image in chat.

1. How to run:
- run `make compose_api` to run the app, require Docker.
- go to `http://localhost:8000/api/v1/docs/index.html#/` to view the app.
- UI is Work in progress.

A user can login, register using email and password

The can create a story, in each story, they can create a qa, each qa is a pair of question and answer. They can contains an attachment, currently support png and jpg image only.

2. Techstack:
- Backend: Go, Echo(Framework) + Postgres SQL Database
- LLM: Ollama (with LLaVA-phi3 (vision+text) and phi3 (text))
- Frontend: HTML + TailwindCSS + HTMX
- API Doc: Swagger Doc
