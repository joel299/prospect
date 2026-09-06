# DRY-RUN — FILA DE PROSPECÇÃO

Data de referência: `2026-09-06T10:00:00-04:00`
Próximo dia útil: `2026-09-08` — terça-feira
Motivo do salto: `2026-09-07` informado como feriado
Timezone: `America/Campo_Grande`
Janela comercial: `08:00–19:00`
Modo: `dry_run`

## Filtro aplicado

```sql
contact_ready IS TRUE
AND do_not_contact IS FALSE
AND converted IS FALSE
AND processing_status = 'pending'
```

## Resultado

Leads encontrados: `0`

Nenhum lead foi reservado, alterado, enviado ou sincronizado com Trello.

## Pendências

- Template Zernio ainda pendente.
- `contact_ready` permanece como campo oficial da tabela; nenhuma coluna foi criada ou alterada.
- Fila somente poderá avançar quando existirem leads com `contact_ready = true`.
- Envio externo permanece desabilitado.
