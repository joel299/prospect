package composio

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ToolExecutor interface {
	Execute(context.Context, string, map[string]any) (map[string]any, error)
}

type Client struct {
	BaseURL string
	APIKey  string
	CLI     string
}

func (c Client) Execute(ctx context.Context, name string, args map[string]any) (map[string]any, error) {
	cli := c.CLI
	if cli == "" {
		cli = "/root/.local/bin/composio"
	}
	switch name {
	case "trello_sync_lead":
		return c.trelloSync(ctx, cli, args)
	case "calcom_get_availability":
		return c.execute(ctx, cli, "CAL_GET_AVAILABLE_SLOTS_INFO", "cal_magged-lemon", args)
	case "calcom_create_booking":
		return c.execute(ctx, cli, "CAL_POST_NEW_BOOKING_REQUEST", "cal_magged-lemon", args)
	default:
		return nil, fmt.Errorf("unsupported Composio operation: %s", name)
	}
}

func (c Client) execute(ctx context.Context, cli, slug, account string, args map[string]any) (map[string]any, error) {
	payload, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	callCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(callCtx, cli, "execute", slug, "--account", account, "-d", string(payload))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Composio %s failed: %w", slug, err)
	}
	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("Composio %s returned invalid JSON", slug)
	}
	if ok, exists := result["successful"].(bool); exists && !ok {
		return nil, fmt.Errorf("Composio %s returned unsuccessful result", slug)
	}
	return result, nil
}

func (c Client) trelloSync(ctx context.Context, cli string, args map[string]any) (map[string]any, error) {
	board, err := c.execute(ctx, cli, "TRELLO_GET_BOARDS_LISTS_BY_ID_BOARD", "trello_lisa-levite", map[string]any{
		"idBoard": "DPwl3zEU", "cards": "all", "fields": "name,pos", "card_fields": "name,desc,idList,due,shortLink,url", "filter": "open",
	})
	if err != nil {
		return nil, err
	}
	leadID := strings.TrimSpace(fmt.Sprint(args["lead_id"]))
	if leadID == "" {
		return nil, fmt.Errorf("lead_id is required")
	}
	name := strings.TrimSpace(fmt.Sprint(args["place_name"]))
	cardID := strings.TrimSpace(fmt.Sprint(args["card_id"]))
	listID := ""
	var cards []map[string]any
	if data, ok := board["data"].(map[string]any); ok {
		if lists, ok := data["lists"].([]any); ok {
			for _, raw := range lists {
				list, _ := raw.(map[string]any)
				if list == nil {
					continue
				}
				if fmt.Sprint(list["name"]) == "Cadastro" {
					listID = fmt.Sprint(list["id"])
				}
				if cs, ok := list["cards"].([]any); ok {
					for _, rc := range cs {
						if c, ok := rc.(map[string]any); ok {
							cards = append(cards, c)
						}
					}
				}
			}
		}
	}
	if cardID == "" {
		for _, card := range cards {
			if (name != "" && strings.EqualFold(strings.TrimSpace(fmt.Sprint(card["name"])), name)) || fmt.Sprint(card["shortLink"]) == leadID {
				cardID = fmt.Sprint(card["id"])
				break
			}
		}
	}
	if cardID == "" {
		if listID == "" {
			return nil, fmt.Errorf("Cadastro list not found")
		}
		create := map[string]any{"idList": listID, "name": name, "desc": description(args), "pos": "bottom"}
		if due := strings.TrimSpace(fmt.Sprint(args["due"])); due != "" && due != "<nil>" {
			create["due"] = due
		}
		created, err := c.execute(ctx, cli, "TRELLO_ADD_CARDS", "trello_lisa-levite", create)
		if err != nil {
			return nil, err
		}
		cardID = findString(created, "id")
		if cardID == "" {
			return nil, fmt.Errorf("Trello create returned no card id")
		}
	}
	comment, err := c.execute(ctx, cli, "TRELLO_ADD_CARDS_ACTIONS_COMMENTS_BY_ID_CARD", "trello_lisa-levite", map[string]any{"idCard": cardID, "text": "Enviar Tamplete"})
	if err != nil {
		return nil, err
	}
	readback, err := c.execute(ctx, cli, "TRELLO_GET_CARDS_ACTIONS_BY_ID_CARD", "trello_lisa-levite", map[string]any{"idCard": cardID, "filter": "commentCard", "limit": "50", "format": "list"})
	if err != nil {
		return nil, err
	}
	return map[string]any{"successful": true, "card_id": cardID, "comment": comment, "readback": readback}, nil
}

func description(args map[string]any) string {
	keys := []string{"place_name", "whatsapp", "email", "website", "city", "address", "source", "lead_status", "pipeline_stage"}
	var b strings.Builder
	for _, k := range keys {
		if v := strings.TrimSpace(fmt.Sprint(args[k])); v != "" && v != "<nil>" {
			fmt.Fprintf(&b, "%s: %s\n", k, v)
		}
	}
	return strings.TrimSpace(b.String())
}

func findString(v map[string]any, key string) string {
	if s, ok := v[key].(string); ok {
		return s
	}
	if d, ok := v["data"].(map[string]any); ok {
		if s, ok := d[key].(string); ok {
			return s
		}
	}
	return ""
}
