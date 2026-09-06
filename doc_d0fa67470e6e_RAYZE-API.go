###[SEND-TEXT]
package main

import (
    "net/http"
    "os"
    "strings"
)

func main() {
    body := strings.NewReader(`{
        "number":  "5511999999999",
        "message": "Olá! Este é um teste do RyzeAPI."
    }`)
    req, _ := http.NewRequest("POST", "https://ryzeapi.cloud/api/message/text/"+os.Getenv("Instance_Name"), body)
    req.Header.Set("token", os.Getenv("Token_Instance"))
    req.Header.Set("Content-Type", "application/json")
    http.DefaultClient.Do(req)
}



###[BOTOOES]

package main

import (
    "net/http"
    "os"
    "strings"
)

func main() {
    body := strings.NewReader(`{
        "number":      "5511999999999",
        "headerText":  "Atendimento RyzeAPI",
        "contentText": "Como podemos ajudar você hoje?",
        "footerText":  "Disponível 24/7",
        "buttons": [
            { "id": "menu_sales",   "displayText": "Falar com vendas", "type": "REPLY" },
            { "id": "menu_support", "displayText": "Suporte técnico",   "type": "REPLY" },
            { "id": "menu_billing", "displayText": "Financeiro",        "type": "REPLY" }
        ]
    }`)
    req, _ := http.NewRequest("POST", "https://ryzeapi.cloud/api/message/button/"+os.Getenv("Instance_Name"), body)
    req.Header.Set("token", os.Getenv("Token_Instance"))
    req.Header.Set("Content-Type", "application/json")
    http.DefaultClient.Do(req)
}





###[HISTORY]
package main

import (
    "net/http"
    "os"
    "strings"
)

func main() {
    body := strings.NewReader(`{"number":"5511999999999"}`)
    req, _ := http.NewRequest("POST", "https://ryzeapi.cloud/api/chat/history/"+os.Getenv("Instance_Name"), body)
    req.Header.Set("token", os.Getenv("Token_Instance"))
    req.Header.Set("Content-Type", "application/json")
    http.DefaultClient.Do(req)
}