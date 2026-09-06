# DESIGN SYSTEM — HERMES OPERATION CONSOLE

## 1. Direção visual

O produto deve parecer:

```text
premium
operacional
tecnológico
minimalista
rápido
confiável
```

Inspirado na composição dos componentes de chat fornecidos, com linguagem visual própria.

Não reproduzir literalmente marcas externas.

---

# 2. Tema

Tema principal:

```text
Dark
```

Suporte futuro a Light, mas não é requisito inicial.

---

# 3. Paleta

## Background

```text
--bg-app:        #090A0C
--bg-surface:    #0F1115
--bg-elevated:   #151820
--bg-glass:      rgba(20, 23, 30, .72)
```

## Border

```text
--border-soft:   rgba(255,255,255,.08)
--border-strong: rgba(255,255,255,.14)
```

## Text

```text
--text-primary:   #F5F7FA
--text-secondary: #A9B0BC
--text-muted:     #6F7785
```

## Accent

```text
--accent-primary: #D7B15A
--accent-hover:   #E6C36F
```

## Status

```text
--success: #33C481
--warning: #E6A94A
--danger:  #E05263
--info:    #5B8DEF
```

## WhatsApp-like semantic

Inbound/outbound devem ser claramente diferentes, sem copiar identidade visual oficial.

---

# 4. Glassmorphism

Usar moderadamente.

```css
background: rgba(16, 18, 24, 0.72);
backdrop-filter: blur(18px);
border: 1px solid rgba(255,255,255,.08);
```

Evitar blur excessivo em listas longas.

Performance primeiro.

---

# 5. Tipografia

Sugestão:

```text
Inter
Geist
```

Escala:

```text
12 caption
13 metadata
14 body small
15 body
16 UI strong
20 section
24 page title
32 hero/internal landing
```

---

# 6. Espaçamento

Base:

```text
4px
```

Escala:

```text
4
8
12
16
20
24
32
40
48
```

---

# 7. Radius

```text
sm  = 8px
md  = 12px
lg  = 16px
xl  = 20px
pill = 999px
```

---

# 8. Shadows

Usar apenas para elevação.

```text
surface shadow
modal shadow
floating composer shadow
```

Não usar glow em excesso.

---

# 9. Layout

## Sidebar

```text
240–280px
```

## Conversation list

```text
300–360px
```

## Inspector

```text
300–360px
```

Chat ocupa o restante.

---

# 10. Componentes

## Button

Variantes:

```text
primary
secondary
ghost
danger
```

## Input

Estados:

```text
default
focus
error
disabled
```

## Badge

```text
Agent ON
Agent OFF
Awaiting Reply
Follow-up
Meeting
Error
```

## Card

Usar para:
- integração;
- métricas;
- error summary;
- tool status.

## Table

Cabeçalho sticky.

Rows compactas.

## Dialog

Para:
- confirmar ação;
- configurar secret;
- reprocessar erro;
- ativar produção.

---

# 11. Chat bubbles

## Inbound

Alinhada esquerda.

## Outbound

Alinhada direita.

Mostrar:
- texto;
- hora;
- status.

Máximo de largura:

```text
72%
```

---

# 12. Composer

Somente:

```text
textarea
send button
agent/human indicator
```

Não mostrar:
- clip;
- mic;
- camera;
- phone;
- video.

---

# 13. Agent Activity

Quando Hermes estiver processando:

```text
Hermes analisando...
```

Quando tool:

```text
Consultando agenda...
```

Quando OmniRoute:

```text
Gerando resposta...
```

Não expor detalhes internos sensíveis.

---

# 14. Error Center

Critical:
- border danger;
- icon;
- severity;
- action.

Não usar vermelho em telas inteiras.

---

# 15. Motion

Duração:

```text
120–220ms
```

Usar:
- fade;
- slide pequeno;
- skeleton;
- subtle scale.

Respeitar:

```text
prefers-reduced-motion
```

---

# 16. Accessibility

- contraste AA;
- keyboard navigation;
- focus visible;
- aria labels;
- semantic HTML;
- status não depender apenas de cor.

---

# 17. Performance

Evitar:
- sombras pesadas em centenas de rows;
- backdrop blur em listas virtuais;
- animação contínua;
- rerender global por evento WebSocket.

Usar:
- virtualization;
- memoization;
- query cache;
- event-specific updates.
